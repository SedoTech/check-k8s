package deployment

import (
	"fmt"
	"github.com/benkeil/icinga-checks-library"
	"k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/api"
	"strings"
	"time"
)

type (
	// CheckDeployment interface to check a deployment
	CheckCronjob interface {
		CheckStatus(CheckStatusOptions) icinga.Result
	}

	checkCronjobImpl struct {
		Client    kubernetes.Interface
		Name      string
		Namespace string
	}
)

// NewCheckCronjob creates a new instance of CheckCronjob
func NewCheckCronjob(client kubernetes.Interface, name string, namespace string) CheckCronjob {
	return &checkCronjobImpl{Client: client, Name: name, Namespace: namespace}
}

// CheckStatusOptions contains options needed to run CheckStatus check
type CheckStatusOptions struct{}

// CheckStatus checks if the cronjob is running properly
func (c *checkCronjobImpl) CheckStatus(options CheckStatusOptions) icinga.Result {
	name := "Cronjob.RunInterval"

	jobLists, err := api.GetJobsByCronjob(c.Client, api.GetJobOptions{Name: c.Name, Namespace: c.Namespace})
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't get jobLists from cronjob: %v", err))
	}

	if len(jobLists.Items) == 0 {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("No jobs found for [%s] in [%s]", c.Name, c.Namespace))
	}

	job := jobLists.Items[len(jobLists.Items)-1]
	jobNameParts := strings.Split(job.Name, "-")
	if jobNameParts[0] == "" || strings.Compare(jobNameParts[0], c.Name) != 0 {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("No job found for Cronjob [%s] (job: [%s])", c.Name, job.Name))
	}

	if job.Status.StartTime == nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("No job status (and start time) found for job [%s]", job.Name))
	}

	startTime, err := time.Parse(api.JOB_TIME_FORMAT, job.Status.StartTime.String())
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't get start time from job [%s]: %v", job.Name, err))
	}

	if job.Status.Active > 0 {
		podRestarts, err := getPodRestarts(c.Client, job)
		if err != nil {
			return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't get pods for job [%s]: %v", job.Name, err))
		}
		if podRestarts > 0 {
			return icinga.NewResult(name, icinga.ServiceStatusCritical, fmt.Sprintf("Job [%s] has too many restarts [%d] since start time [%v]", job.Name, podRestarts, startTime.String()))
		}

		elapsedTime := time.Now().Second() - startTime.Second()

		return icinga.NewResult(name, icinga.ServiceStatusOk, fmt.Sprintf("Job [%s] running since %d seconds for Cronjob [%s]", job.Name, elapsedTime, c.Name))
	}

	if job.Status.Failed > 0 {
		return icinga.NewResult(name, icinga.ServiceStatusCritical, fmt.Sprintf("Job [%s] failed for Cronjob [%s] at start time [%v]", job.Name, c.Name, startTime.String()))
	}

	if job.Status.Succeeded > 0 {
		return icinga.NewResult(name, icinga.ServiceStatusOk, fmt.Sprintf("Job [%s] succeeded for Cronjob [%s] at completion time [%v]", job.Name, c.Name, job.Status.CompletionTime))
	}

	return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("Unknown running specs for job [%s] in Cronjob [%s] - [%v]", job.Name, c.Name, job))
}

func getPodRestarts(client kubernetes.Interface, job v1.Job) (int32, error) {
	podList, err := api.GetPods(client, api.GetPodOptions{LabelSelector: job.Spec.Selector})
	if err != nil {
		return 0, err
	}

	for _, pod := range podList.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready && containerStatus.RestartCount > 0 {
				return containerStatus.RestartCount, nil
			}
		}
	}

	return 0, nil
}
