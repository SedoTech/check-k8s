package api

import (
	"fmt"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetConfigMapOptions options to get a deployment
	GetConfigMapOptions struct {
		Name      string
		Namespace string
	}
)

var configmaps = make(map[string]*v1.ConfigMap)

// GetConfigMap returns a single k8s ConfigMap resource
func GetConfigMap(client kubernetes.Interface, options GetConfigMapOptions) (*v1.ConfigMap, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if secret, found := configmaps[key]; found {
		return secret, nil
	}
	resource, err := client.CoreV1().ConfigMaps(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	configmaps[key] = resource
	return resource, nil
}

type (
	// GetConfigMapsOptions options to get a deployment
	GetConfigMapsOptions struct {
		Namespace string
	}
)

// GetConfigMaps returns a ConfigMapList
func GetConfigMaps(client kubernetes.Interface, options GetConfigMapsOptions) (*v1.ConfigMapList, error) {
	resources, err := client.CoreV1().ConfigMaps(options.Namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resources, nil
}
