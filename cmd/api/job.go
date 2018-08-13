package api

import (
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/batch/v1"
)

type (
	// GetJobOptions options to get a job
	GetJobOptions struct {
		Name               string
		Namespace          string
	}
)

var jobs = make(map[string]*v1.Job)

// GetJob returns a new checkJob object
func GetJob(client kubernetes.Interface, options GetJobOptions) (*v1.Job, error) {
	key := fmt.Sprintf("%s-%s", options.Name, options.Namespace)
	if job, found := jobs[key]; found {
		return job, nil
	}
	resource, err := client.BatchV1().Jobs(options.Namespace).Get(options.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	jobs[key] = resource

	return resource, nil
}

// GetJobs returns a new checkJob object
func GetJobsByCronjob(client kubernetes.Interface, options GetJobOptions) (*v1.JobList, error) {

	//sel := strings.Join([]string{"items.metadata.ownerReferences.name", options.Name}, "=")
	//log.Printf("sel: %v", sel)
	//resource, err := client.BatchV1().Jobs(options.Namespace).List(meta_v1.ListOptions{LabelSelector: sel})

	resource, err := client.BatchV1().Jobs(options.Namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return resource, nil
}
