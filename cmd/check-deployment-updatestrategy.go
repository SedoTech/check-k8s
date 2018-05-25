package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks/deployment"
	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	icinga "github.com/benkeil/icinga-checks-library"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
)

type (
	checkDeploymentUpdateStrategyCmd struct {
		out            io.Writer
		Client         kubernetes.Interface
		Name           string
		Namespace      string
		Result         string
		UpdateStrategy string
	}
)

func newCheckDeploymentUpdateStrategyCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkDeploymentUpdateStrategyCmd{out: out}

	cmd := &cobra.Command{
		Use:          "updateStrategy",
		Short:        "check if a k8s deployment has a specific update strategy defined",
		SilenceUsage: true,
		Args:         NameArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			c.Name = args[0]
			client, err := kube.GetKubeClient(settings.KubeContext)
			if err != nil {
				icinga.NewResult("GetKubeClient", icinga.ServiceStatusUnknown, fmt.Sprintf("can't get client: %v", err)).Exit()
			}
			c.Client = client
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the deployment")
	cmd.Flags().StringVarP(&c.Result, "result", "r", "WARNING", "the result state if the check fails")
	cmd.Flags().StringVarP(&c.UpdateStrategy, "string", "s", "RollingUpdate", "the expected update strategy")

	return cmd
}

func (c *checkDeploymentUpdateStrategyCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Client, c.Name, c.Namespace)
	result := checkDeployment.CheckUpdateStrategy(
		deployment.CheckUpdateStrategyOptions{
			Result:         c.Result,
			UpdateStrategy: c.UpdateStrategy,
		})
	result.Exit()
}
