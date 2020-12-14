package api

import (
	"fmt"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetEndpointsOptions options to get a deployment
	GetEndpointsOptions struct {
		Name      string
		Namespace string
	}
)

var services = make(map[string]*v1.Endpoints)

// GetEndpoints returns a new checkDeployment object
func GetEndpoints(client kubernetes.Interface, options GetEndpointsOptions) (*v1.Endpoints, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if deployment, found := services[key]; found {
		return deployment, nil
	}
	resource, err := client.CoreV1().Endpoints(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	services[key] = resource
	return resource, nil
}
