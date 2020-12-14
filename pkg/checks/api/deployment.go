package api

import (
	"fmt"

	"k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetDeploymentOptions options to get a deployment
	GetDeploymentOptions struct {
		Name      string
		Namespace string
	}
)

var deployments = make(map[string]*v1.Deployment)

// GetDeployment returns a new checkDeployment object
func GetDeployment(client kubernetes.Interface, options GetDeploymentOptions) (*v1.Deployment, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if deployment, found := deployments[key]; found {
		return deployment, nil
	}
	resource, err := client.AppsV1().Deployments(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	deployments[key] = resource
	return resource, nil
}
