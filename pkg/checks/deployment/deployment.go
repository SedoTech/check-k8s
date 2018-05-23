package deployment

import (
	"fmt"

	"github.com/benkeil/check-k8s/cmd/api"
	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/print"
	"github.com/benkeil/icinga-checks-library"
	"k8s.io/api/apps/v1"
)

type (
	// CheckDeployment interface to check a deployment
	CheckDeployment interface {
		CheckUpdateStrategy(CheckUpdateStrategyOptions) icinga.Result
		CheckAvailableReplicas(CheckAvailableReplicasOptions) icinga.Result
		CheckAll(CheckAllOptions) icinga.Results
		PrintYaml()
	}

	checkDeploymentImpl struct {
		deployment *v1.Deployment
	}
)

func (c *checkDeploymentImpl) PrintYaml() {
	print.Yaml(c.deployment)
}

// NewCheckDeployment creates a new instance of CheckDeployment
func NewCheckDeployment(deployment *v1.Deployment) CheckDeployment {
	return &checkDeploymentImpl{deployment}
}

// CheckUpdateStrategyOptions contains options needed to run CheckUpdateStrategy check
type CheckUpdateStrategyOptions struct {
	Result         string
	UpdateStrategy string
}

// CheckUpdateStrategy checks if the deployment has the RollingUpdate strategy
func (c *checkDeploymentImpl) CheckUpdateStrategy(options CheckUpdateStrategyOptions) icinga.Result {
	name := "Deployment.UpdateStretegy"
	var updateStretegy v1.DeploymentStrategyType
	switch options.UpdateStrategy {
	case "RollingUpdate":
		updateStretegy = v1.RollingUpdateDeploymentStrategyType
	case "Recreate":
		updateStretegy = v1.RecreateDeploymentStrategyType
	default:
		panic("invalid v1.DeploymentStrategyType")
	}

	statusCheck, err := icinga.NewStatusCheckCompare(options.Result)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't compare status: %v", err))
	}
	comparator := func() bool {
		return updateStretegy != c.deployment.Spec.Strategy.Type
	}
	status := statusCheck.Compare(comparator)
	return icinga.NewResult(name, status, fmt.Sprintf("deployment has update strategy %s", updateStretegy))
}

// CheckAvailableReplicasOptions contains options needed to run CheckAvailableReplicas check
type CheckAvailableReplicasOptions struct {
	ThresholdWarning  string
	ThresholdCritical string
}

// CheckAvailableReplicas checks if the deployment has a minimum of available replicas
func (c *checkDeploymentImpl) CheckAvailableReplicas(options CheckAvailableReplicasOptions) icinga.Result {
	name := "Deployment.AvailableReplicas"
	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}
	replicas := c.deployment.Status.AvailableReplicas
	status := statusCheck.CheckInt32(replicas)
	message := fmt.Sprintf("deployment has %v available replica(s) [thresholds warn: %v crit: %v]", replicas, options.ThresholdWarning, options.ThresholdCritical)

	return icinga.NewResult(name, status, message)
}

// CheckPodRestartsOptions contains options needed to run CheckPodRestarts check
type CheckPodRestartsOptions struct {
	Settings          environment.EnvSettings
	ThresholdWarning  string
	ThresholdCritical string
}

// CheckPodRestarts checks if the deployment has a minimum of available replicas
func (c *checkDeploymentImpl) CheckPodRestarts(options CheckPodRestartsOptions) icinga.Result {
	name := "Deployment.PodRestarts"
	api.GetPods(options.Settings, api.GetPodOptions{LabelSelector: c.deployment.Spec.Selector})
	//client, err := kube.GetKubeClient(options.settings.KubeContext)
	//labelSelector := labels.Set(map[string]string{"mylabel": "ourdaomain1"}).AsSelector()
	//resource, err := client.CoreV1().(c.namespace).List(meta_v1.ListOptions{})
	return icinga.NewResult(name, icinga.ServiceStatusOk, icinga.DefaultSuccessMessage)
}

// CheckAllOptions contains options needed to run all deployment checks
type CheckAllOptions struct {
	Settings                      environment.EnvSettings
	CheckUpdateStrategyOptions    CheckUpdateStrategyOptions
	CheckAvailableReplicasOptions CheckAvailableReplicasOptions
	CheckPodRestartsOptions       CheckPodRestartsOptions
}

// CheckAll runs all tests and returns an instance of ServiceCheckResults
func (c *checkDeploymentImpl) CheckAll(options CheckAllOptions) icinga.Results {
	results := icinga.NewResults()
	results.Add(c.CheckUpdateStrategy(options.CheckUpdateStrategyOptions))
	results.Add(c.CheckAvailableReplicas(options.CheckAvailableReplicasOptions))
	results.Add(c.CheckPodRestarts(options.CheckPodRestartsOptions))
	return results
}
