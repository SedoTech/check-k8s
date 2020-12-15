package main

import (
	"fmt"
	"io"

	icinga "github.com/SedoTech/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/deployment"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
)

type (
	checkDeploymentContainerDefinedCmd struct {
		out              io.Writer
		Client           kubernetes.Interface
		Name             string
		Namespace        string
		Result           string
		ContainerDefined []string
	}
)

func newCheckDeploymentContainerDefinedCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkDeploymentContainerDefinedCmd{out: out}

	cmd := &cobra.Command{
		Use:          "containerDefined",
		Short:        "check if a k8s deployment has a list of containers defined",
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
	cmd.Flags().StringVarP(&c.Result, "result", "r", "CRITICAL", "the result state if the check fails")
	cmd.Flags().StringSliceVarP(&c.ContainerDefined, "string", "s", []string{}, "the containers to check if they are present (comma separated list)")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentContainerDefinedCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Client, c.Name, c.Namespace)
	result := checkDeployment.CheckContainerDefined(
		deployment.CheckContainerDefinedOptions{
			Result:           c.Result,
			ContainerDefined: c.ContainerDefined,
		})
	result.Exit()
}
