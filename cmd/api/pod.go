package api

import (
	"fmt"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type (
	// GetPodOptions options to get a deployment
	GetPodOptions struct {
		Name          string
		Namespace     string
		LabelSelector *meta_v1.LabelSelector
	}
)

var pods = make(map[string]*v1.Pod)
var podLists = make(map[string]*v1.PodList)

// GetPod returns a new checkDeployment object
func GetPod(settings environment.EnvSettings, options GetPodOptions) (*v1.Pod, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if pod, found := pods[key]; found {
		return pod, nil
	}
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return nil, err
	}
	resource, err := client.CoreV1().Pods(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	pods[key] = resource
	return resource, nil
}

// GetPods returns a new checkDeployment object
func GetPods(settings environment.EnvSettings, options GetPodOptions) (*v1.PodList, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if podList, found := podLists[key]; found {
		return podList, nil
	}
	client, err := kube.GetKubeClient(settings.KubeContext)
	if err != nil {
		return nil, err
	}

	labelSelector := labels.Set(options.LabelSelector.MatchLabels).AsSelector()
	fmt.Printf("selector: %v\n", labelSelector.String())
	resources, err := client.CoreV1().Pods(options.Namespace).List(meta_v1.ListOptions{LabelSelector: labelSelector.String()})
	print.Yaml(resources)

	if err != nil {
		return nil, err
	}
	podLists[key] = resources
	return resources, nil
}
