package api

import (
	"fmt"

	"k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetStatefulSetOptions options to get a StatefulSet
	GetStatefulSetOptions struct {
		Name      string
		Namespace string
	}
)

var statefulSets = make(map[string]*v1.StatefulSet)

// GetStatefulSet returns a new checkStatefulSet object
func GetStatefulSet(client kubernetes.Interface, options GetStatefulSetOptions) (*v1.StatefulSet, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if statefulSet, found := statefulSets[key]; found {
		return statefulSet, nil
	}
	resource, err := client.AppsV1().StatefulSets(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	statefulSets[key] = resource
	return resource, nil
}
