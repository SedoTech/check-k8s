package statefulset

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/benkeil/icinga-checks-library"
	"k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/api"
	"SedoTech/check-k8s/pkg/utils"
)

type (
	// CheckStatefulSet interface to check a StatefulSet
	CheckStatefulSet interface {
		CheckUpdateStrategy(CheckUpdateStrategyOptions) icinga.Result
		CheckAvailableReplicas(CheckAvailableReplicasOptions) icinga.Result
		CheckPodRestarts(CheckPodRestartsOptions) icinga.Result
		CheckProbesDefined(CheckProbesDefinedOptions) icinga.Result
		CheckContainerDefined(CheckContainerDefinedOptions) icinga.Result
		CheckAll(CheckAllOptions) icinga.Results
	}

	checkStatefulSetImpl struct {
		Client    kubernetes.Interface
		Name      string
		Namespace string
	}
)

// NewCheckStatefulSet creates a new instance of CheckStatefulSet
func NewCheckStatefulSet(client kubernetes.Interface, name string, namespace string) CheckStatefulSet {
	return &checkStatefulSetImpl{Client: client, Name: name, Namespace: namespace}
}

// CheckAllOptions contains options needed to run all StatefulSet checks
type CheckAllOptions struct {
	Client                        kubernetes.Interface
	CheckUpdateStrategyOptions    CheckUpdateStrategyOptions
	CheckAvailableReplicasOptions CheckAvailableReplicasOptions
	CheckPodRestartsOptions       CheckPodRestartsOptions
	CheckProbesDefinedOptions     CheckProbesDefinedOptions
	CheckContainerDefinedOptions  CheckContainerDefinedOptions
}

// CheckAll runs all tests and returns an instance of ServiceCheckResults
func (c *checkStatefulSetImpl) CheckAll(options CheckAllOptions) icinga.Results {
	results := icinga.NewResults()
	results.Add(c.CheckUpdateStrategy(options.CheckUpdateStrategyOptions))
	results.Add(c.CheckAvailableReplicas(options.CheckAvailableReplicasOptions))
	results.Add(c.CheckPodRestarts(options.CheckPodRestartsOptions))
	results.Add(c.CheckProbesDefined(options.CheckProbesDefinedOptions))
	results.Add(c.CheckContainerDefined(options.CheckContainerDefinedOptions))
	return results
}

// CheckUpdateStrategyOptions contains options needed to run CheckUpdateStrategy check
type CheckUpdateStrategyOptions struct {
	Result         string
	UpdateStrategy string
}

// CheckUpdateStrategy checks if the StatefulSet has the RollingUpdate strategy
func (c *checkStatefulSetImpl) CheckUpdateStrategy(options CheckUpdateStrategyOptions) icinga.Result {
	name := "StatefulSet.UpdateStrategy"
	var updateStrategy v1.StatefulSetUpdateStrategyType
	switch options.UpdateStrategy {
	case "RollingUpdate":
		updateStrategy = v1.RollingUpdateStatefulSetStrategyType
	case "OnDelete":
		updateStrategy = v1.OnDeleteStatefulSetStrategyType
	default:
		icinga.NewResult("CheckUpdateStrategy", icinga.ServiceStatusUnknown, fmt.Sprintf("invalid StatefulSetStrategy: %v", options.UpdateStrategy)).Exit()
	}

	statusCheck, err := icinga.NewStatusCheckCompare(options.Result)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't compare status: %v", err))
	}

	statefulSet, err := api.GetStatefulSet(c.Client, api.GetStatefulSetOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetStatefulSet", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	comparator := func() bool {
		return updateStrategy != statefulSet.Spec.UpdateStrategy.Type
	}
	status := statusCheck.Compare(comparator)

	return icinga.NewResult(name, status, fmt.Sprintf("StatefulSet has update strategy %s", updateStrategy))
}

// CheckAvailableReplicasOptions contains options needed to run CheckAvailableReplicas check
type CheckAvailableReplicasOptions struct {
	ThresholdWarning  string
	ThresholdCritical string
}

// CheckAvailableReplicas checks if the StatefulSet has a minimum of available replicas
func (c *checkStatefulSetImpl) CheckAvailableReplicas(options CheckAvailableReplicasOptions) icinga.Result {
	name := "StatefulSet.AvailableReplicas"

	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	statefulSet, err := api.GetStatefulSet(c.Client, api.GetStatefulSetOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetStatefulSet", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	replicas := statefulSet.Status.ReadyReplicas
	status := statusCheck.CheckInt32(replicas)
	message := fmt.Sprintf("StatefulSet has %v available replica(s)", replicas)

	return icinga.NewResult(name, status, message)
}

// CheckPodRestartsOptions contains options needed to run CheckPodRestarts check
type CheckPodRestartsOptions struct {
	Result   string
	Duration string
}

// CheckPodRestarts checks if the StatefulSet has a minimum of available replicas
func (c *checkStatefulSetImpl) CheckPodRestarts(options CheckPodRestartsOptions) icinga.Result {
	name := "StatefulSet.PodRestarts"

	statefulSet, err := api.GetStatefulSet(c.Client, api.GetStatefulSetOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResultUnknownMessage("GetStatefulSet", fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	podList, err := api.GetPods(c.Client, api.GetPodOptions{LabelSelector: statefulSet.Spec.Selector})
	if err != nil {
		return icinga.NewResultUnknownMessage("GetPods", fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	duration, err := time.ParseDuration(options.Duration)
	if err != nil {
		return icinga.NewResultUnknownMessage("ParseDuration", fmt.Sprintf("can't parse duration: %v", err))
	}

	statusCheck, err := icinga.NewStatusCheckCompare(options.Result)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	// maybe we need a field selector for the events to check that
	// kubectl get event -n production-mls --field-selector=involvedObject.name=mls-67799db556-4wbtw
	// kubectl get event -n production-mls --field-selector=involvedObject.name=mls-67799db556-4wbtw,reason=Killing
	// Reasons:
	//   - Unhealthy
	//   - Killing
	// Types:
	//   - Warning

	// contains faild containers grouped by pod name
	failedContainerMap := make(map[string][]string)

	for _, pod := range podList.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			terminated := containerStatus.LastTerminationState.Terminated
			if terminated != nil && time.Since(terminated.FinishedAt.Time).Minutes() < duration.Minutes() {
				failedContainerMap[pod.GetObjectMeta().GetName()] = append(failedContainerMap[pod.GetObjectMeta().GetName()], containerStatus.Name)
			}
		}
	}

	status := statusCheck.CompareBool(len(failedContainerMap) > 0)
	message := icinga.DefaultSuccessMessage
	if status != icinga.ServiceStatusOk {
		var buffer bytes.Buffer
		for podName, containers := range failedContainerMap {
			buffer.WriteString(fmt.Sprintf("%v: %v ", podName, containers))
		}
		message = buffer.String()
	}

	return icinga.NewResult(name, status, strings.Trim(message, "\n"))
}

// CheckProbesDefinedOptions contains options needed to run CheckProbesDefined check
type CheckProbesDefinedOptions struct {
	Result        string
	ProbesDefined []string
}

// CheckProbesDefined checks if the StatefulSet has the RollingUpdate strategy
func (c *checkStatefulSetImpl) CheckProbesDefined(options CheckProbesDefinedOptions) icinga.Result {
	name := "StatefulSet.ProbesDefined"

	statusCheck, err := icinga.NewStatusCheckCompare(options.Result)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't compare status: %v", err))
	}

	statefulSet, err := api.GetStatefulSet(c.Client, api.GetStatefulSetOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult("GetStatefulSet", icinga.ServiceStatusUnknown, fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	missingProbes := []string{}
	for _, container := range statefulSet.Spec.Template.Spec.Containers {
		if (len(options.ProbesDefined) > 0 && utils.Contains(container.Name, options.ProbesDefined)) || len(options.ProbesDefined) == 0 {
			if container.ReadinessProbe == nil || container.LivenessProbe == nil {
				missingProbes = append(missingProbes, container.Name)
			}
		}
	}
	status := statusCheck.CompareBool(len(missingProbes) > 0)

	message := icinga.DefaultSuccessMessage
	if status != icinga.ServiceStatusOk {
		message = fmt.Sprintf("containers without probes: %v", missingProbes)
	}

	return icinga.NewResult(name, status, message)
}

// CheckContainerDefinedOptions contains options needed to run CheckContainerDefined check
type CheckContainerDefinedOptions struct {
	Result           string
	ContainerDefined []string
}

// CheckContainerDefined checks if the StatefulSet has the RollingUpdate strategy
func (c *checkStatefulSetImpl) CheckContainerDefined(options CheckContainerDefinedOptions) icinga.Result {
	name := "StatefulSet.ContainerDefined"

	if len(options.ContainerDefined) == 0 {
		return icinga.NewResultUnknownMessage(name, fmt.Sprint("no containers defined to check"))
	}

	statusCheck, err := icinga.NewStatusCheckCompare(options.Result)
	if err != nil {
		return icinga.NewResultUnknownMessage(name, fmt.Sprintf("can't compare status: %v", err))
	}

	statefulSet, err := api.GetStatefulSet(c.Client, api.GetStatefulSetOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResultUnknownMessage("GetStatefulSet", fmt.Sprintf("cant't get StatefulSet: %v", err))
	}

	foundContainers := make(map[string]bool)
	for _, container := range options.ContainerDefined {
		foundContainers[container] = false
	}

	for _, container := range statefulSet.Spec.Template.Spec.Containers {
		if utils.Contains(container.Name, options.ContainerDefined) {
			foundContainers[container.Name] = true
		}
	}

	missingContainer := []string{}
	for container, found := range foundContainers {
		if !found {
			missingContainer = append(missingContainer, container)
		}
	}

	status := statusCheck.CompareBool(len(missingContainer) > 0)

	message := icinga.DefaultSuccessMessage
	if status != icinga.ServiceStatusOk {
		message = fmt.Sprintf("missing containers: %v", missingContainer)
	}

	return icinga.NewResult(name, status, message)
}
