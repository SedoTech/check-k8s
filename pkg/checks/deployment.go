package checks

import (
	"fmt"

	"github.com/benkeil/check-k8s/pkg/icinga"
	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	"k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	// CheckDeployment interface to check a deployment
	CheckDeployment interface {
		//CheckUpdateStrategy() ServiceCheckResult
		CheckAvailableReplicas(string, string) ServiceCheckResult
		//CheckAll() ServiceCheckResults
		PrintYaml()
	}

	checkDeploymentImpl struct {
		settings   environment.EnvSettings
		name       string
		deployment *v1.Deployment
	}
)

// NewCheckDeployment returns a new checkDeployment object
func NewCheckDeployment(settings environment.EnvSettings, name string) (CheckDeployment, error) {
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return nil, err
	}
	resource, err := client.AppsV1().Deployments("production-mls").Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &checkDeploymentImpl{settings, name, resource}, nil
}

// CheckUpdateStrategy checks if the deployment has the RollingUpdate strategy
func (c *checkDeploymentImpl) CheckUpdateStrategy() ServiceCheckResult {
	name := "Deployment.Spec.Strategy.Type"
	updateStretegy := c.deployment.Spec.Strategy.Type
	if updateStretegy != v1.RollingUpdateDeploymentStrategyType {
		return NewServiceCheckResult(name, icinga.ServiceStateWarning, fmt.Sprintf("deployment has update strategy %s", updateStretegy))
	}
	return NewServiceCheckResultOk(name)
}

// CheckSpecReplicas checks if the deployment has a minimum of replicas specified
//func (c *checkDeploymentImpl) CheckSpecReplicas(w string, c string) ServiceCheckResult {
//	name := "Deployment.Spec.Replicas"
//	minimum := c.options.MinimumSpecReplicas
//	replicas := c.deployment.Spec.Replicas
//	if *replicas < minimum {
//		return NewServiceCheckResult(name, icinga.ServiceStateWarning, fmt.Sprintf("deployment has %v of minimum %v replicas specified", *replicas, minimum))
//	}
//	return NewServiceCheckResultOk(name)
//}

// CheckAvailableReplicas checks if the deployment has a minimum of available replicas
func (c *checkDeploymentImpl) CheckAvailableReplicas(tw string, tc string) ServiceCheckResult {
	name := "Deployment.Status.AvailableReplicas"
	minimum := c.options.MinimumAvailableReplicas
	replicas := c.deployment.Status.AvailableReplicas
	if replicas < minimum {
		return NewServiceCheckResult(name, icinga.ServiceStateWarning, fmt.Sprintf("deployment has %v of minimum %v available replicas", replicas, minimum))
	}
	return NewServiceCheckResultOk(name)
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
