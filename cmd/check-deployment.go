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
	checkDeploymentCmd struct {
		Settings                  environment.EnvSettings
		out                       io.Writer
		Deployment                *v1.Deployment
		Name                      string
		Namespace                 string
		AvailableReplicasWarning  string
		AvailableReplicasCritical string
		UpdateStrategyResult      string
		UpdateStrategyValue       string
	}
)

func newCheckDeploymentCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkDeploymentCmd{out: out, Settings: settings}

	cmd := &cobra.Command{
		Use:          "deployment",
		Short:        "check if a k8s deployment resource is healthy",
		SilenceUsage: false,
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
	cmd.AddCommand(
		newCheckDeploymentAvailableReplicasCmd(settings, out),
		newCheckDeploymentUpdateStrategyCmd(settings, out),
	)

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace where the deployment is")
	cmd.Flags().StringVar(&c.AvailableReplicasWarning, "availableReplicasWarning", "2:", "minimum of replicas in spec")
	cmd.Flags().StringVar(&c.AvailableReplicasCritical, "availableReplicasCritical", "2:", "minimum of available replicas")
	cmd.Flags().StringVar(&c.UpdateStrategyResult, "updateStrategyResult", "WARNING", "minimum of available replicas")
	cmd.Flags().StringVar(&c.UpdateStrategyValue, "updateStrategyString", "RollingUpdate", "minimum of available replicas")

	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Deployment)
	results := checkDeployment.CheckAll(deployment.CheckAllOptions{
		CheckAvailableReplicasOptions: deployment.CheckAvailableReplicasOptions{ThresholdWarning: c.AvailableReplicasWarning, ThresholdCritical: c.AvailableReplicasCritical},
		CheckUpdateStrategyOptions:    deployment.CheckUpdateStrategyOptions{Result: c.UpdateStrategyResult, UpdateStrategy: c.UpdateStrategyValue},
		CheckPodRestartsOptions:       deployment.CheckPodRestartsOptions{Settings: c.Settings, ThresholdWarning: c.UpdateStrategyResult, ThresholdCritical: c.UpdateStrategyValue},
	})
	results.Exit()
}
