package main

import (
	"fmt"
	"io"

	icinga "github.com/SedoTech/icinga-checks-library"
	"SedoTech/check-k8s/pkg/checks/configmaps"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

type (
	checkConfigMapsCmd struct {
		out               io.Writer
		Client            kubernetes.Interface
		Name              string
		Namespace         string
		ConfigMapsDefined []string
	}
)

func newCheckConfigMapsCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkConfigMapsCmd{out: out}

	cmd := &cobra.Command{
		Use:          "configmaps",
		Short:        "check if k8s configmap resources exists",
		SilenceUsage: false,
		Args:         cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
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

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the endpoint")
	cmd.Flags().StringSliceVarP(&c.ConfigMapsDefined, "string", "s", []string{}, "the configmap to check if they are exists (comma separated list)")
	cmd.MarkPersistentFlagRequired("namespace")
	cmd.MarkFlagRequired("string")

	return cmd
}

func (c *checkConfigMapsCmd) run() {
	checkConfigmaps := configmaps.NewCheckConfigMaps(c.Client, c.Name, c.Namespace)
	result := checkConfigmaps.CheckExists(
		configmaps.CheckExistsOptions{
			ConfigMapsDefined: c.ConfigMapsDefined,
		},
	)
	result.Exit()
}
