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
	checkDeploymentProbesDefinedCmd struct {
		out           io.Writer
		Client        kubernetes.Interface
		Name          string
		Namespace     string
		Result        string
		ProbesDefined []string
	}
)

func newCheckDeploymentProbesDefinedCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkDeploymentProbesDefinedCmd{out: out}

	cmd := &cobra.Command{
		Use:          "probesDefined",
		Short:        "check if a k8s deployment has liveness and readiness probes defined",
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
	cmd.Flags().StringSliceVarP(&c.ProbesDefined, "string", "s", []string{}, "check only the defined containers, not all")

	return cmd
}

func (c *checkDeploymentProbesDefinedCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Client, c.Name, c.Namespace)
	result := checkDeployment.CheckProbesDefined(
		deployment.CheckProbesDefinedOptions{
			Result:        c.Result,
			ProbesDefined: c.ProbesDefined,
		})
	result.Exit()
}
