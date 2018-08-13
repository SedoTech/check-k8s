package api

import (
	"fmt"
	"k8s.io/api/batch/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	// GetCronjobOptions options to get a cronjob
	GetCronjobOptions struct {
		Name      string
		Namespace string
	}
)

const JOB_TIME_FORMAT = "2006-01-02 15:04:05 +0200 CEST"

var cronjobs = make(map[string]*v1beta1.CronJob)

// GetCronjob returns a new checkCronjob object
func GetCronjob(client kubernetes.Interface, options GetCronjobOptions) (*v1beta1.CronJob, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if cronjob, found := cronjobs[key]; found {
		return cronjob, nil
	}
	resource, err := client.BatchV1beta1().CronJobs(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cronjobs[key] = resource

	return resource, nil
}
