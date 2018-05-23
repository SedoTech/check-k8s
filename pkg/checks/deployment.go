package checks

import (
	"fmt"

	"github.com/benkeil/check-k8s/pkg/print"
	"github.com/benkeil/icinga-checks-library"
	"k8s.io/api/apps/v1"
)

type (
	// CheckDeployment interface to check a deployment
	CheckDeployment interface {
		CheckUpdateStrategy(string, string) icinga.Result
		CheckAvailableReplicas(string, string) icinga.Result
		//CheckAll() ServiceCheckResults
		PrintYaml()
	}

	checkDeploymentImpl struct {
		deployment *v1.Deployment
	}
)

// NewCheckDeployment creates a new instance of CheckDeployment
func NewCheckDeployment(deployment *v1.Deployment) CheckDeployment {
	return &checkDeploymentImpl{deployment}
}

// CheckUpdateStrategy checks if the deployment has the RollingUpdate strategy
func (c *checkDeploymentImpl) CheckUpdateStrategy(strategy string, result string) icinga.Result {
	name := "Deployment.Spec.Strategy.Type"
	var updateStretegy v1.DeploymentStrategyType
	switch strategy {
	case "RollingUpdate":
		updateStretegy = v1.RollingUpdateDeploymentStrategyType
	case "Recreate":
		updateStretegy = v1.RecreateDeploymentStrategyType
	default:
		panic("invalid v1.DeploymentStrategyType")
	}

	statusCheck, err := icinga.NewStatusCheckCompare(result)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't compare status: %v", err))
	}
	comparator := func() bool {
		return updateStretegy != c.deployment.Spec.Strategy.Type
	}
	status := statusCheck.Compare(comparator)
	return icinga.NewResult(name, status, fmt.Sprintf("deployment has update strategy %s", updateStretegy))
}

// CheckAvailableReplicas checks if the deployment has a minimum of available replicas
func (c *checkDeploymentImpl) CheckAvailableReplicas(tw string, tc string) icinga.Result {
	name := "Deployment.Status.AvailableReplicas"
	statusCheck, err := icinga.NewStatusCheck(tw, tc)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}
	replicas := c.deployment.Status.AvailableReplicas
	status := statusCheck.CheckInt32(replicas)
	message := fmt.Sprintf("deployment has %v available replica(s) [thresholds warn: %v crit: %v]", replicas, tw, tc)

	return icinga.NewResult(name, status, message)
}

// CheckAll runs all tests and returns an instance of ServiceCheckResults
//func (c *checkDeploymentImpl) CheckAll() ServiceCheckResults {
//	results := NewServiceCheckResults()
//	results.Add(c.CheckUpdateStrategy())
//	results.Add(c.CheckAvailableReplicas())
//	return results
//}

func (c *checkDeploymentImpl) PrintYaml() {
	print.Yaml(c.deployment)
}
