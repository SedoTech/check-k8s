package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks"
	icinga "github.com/benkeil/icinga-checks-library"

	"github.com/benkeil/check-k8s/cmd/api"
	"github.com/spf13/cobra"
	"k8s.io/api/apps/v1"
)

type (
	checkDeploymentUpdateStrategyCmd struct {
		out            io.Writer
		Deployment     *v1.Deployment
		Name           string
		Namespace      string
		ReturnStatus   string
		UpdateStrategy string
	}
)

func newCheckDeploymentUpdateStrategyCmd(out io.Writer) *cobra.Command {
	c := &checkDeploymentUpdateStrategyCmd{out: out}

	cmd := &cobra.Command{
		Use:          "updateStrategy",
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
	cmd.Flags().StringVarP(&c.ReturnStatus, "result", "r", "WARNING", "the result state if the check fails")
	cmd.Flags().StringVarP(&c.UpdateStrategy, "string", "s", "RollingUpdate", "the expected update strategy")

	return cmd
}

func (c *checkDeploymentUpdateStrategyCmd) run() {
	checkDeployment := checks.NewCheckDeployment(c.Deployment)
	result := checkDeployment.CheckUpdateStrategy(c.UpdateStrategy, c.ReturnStatus)
	result.Exit()
}
