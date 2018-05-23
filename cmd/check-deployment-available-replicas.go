package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks/deployment"
	"github.com/benkeil/check-k8s/pkg/environment"
	icinga "github.com/benkeil/icinga-checks-library"

	"github.com/benkeil/check-k8s/cmd/api"
	"github.com/spf13/cobra"
	"k8s.io/api/apps/v1"
)

type (
	checkDeploymentAvailableReplicasCmd struct {
		out               io.Writer
		Deployment        *v1.Deployment
		Name              string
		Namespace         string
		ThresholdWarning  string
		ThresholdCritical string
	}
)

func newCheckDeploymentAvailableReplicasCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkDeploymentAvailableReplicasCmd{out: out}

	cmd := &cobra.Command{
		Use:          "availableReplicas",
		Short:        "check if a k8s deployment has a minimum of available replicas",
		SilenceUsage: true,
		Args:         NameArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			c.Name = args[0]
			deployment, err := api.GetDeployment(settings, api.GetDeploymentOptions{Name: c.Name, Namespace: c.Namespace})
			if err != nil {
				icinga.NewResult("GetDeployment", icinga.ServiceStatusUnknown, fmt.Sprintf("can't get deployment: %v", err)).Exit()
			}
			c.Deployment = deployment
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the deployment")
	cmd.Flags().StringVarP(&c.ThresholdCritical, "critical", "c", "2:", "critical threshold for minimum available replicas")
	cmd.Flags().StringVarP(&c.ThresholdWarning, "warning", "w", "2:", "warning threshold for minimum available replicas")

	return cmd
}

func (c *checkDeploymentAvailableReplicasCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Deployment)
	result := checkDeployment.CheckAvailableReplicas(deployment.CheckAvailableReplicasOptions{ThresholdWarning: c.ThresholdWarning, ThresholdCritical: c.ThresholdCritical})
	result.Exit()
}
