package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks/endpoints"
	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	icinga "github.com/benkeil/icinga-checks-library"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
)

type (
	checkEndpointsAvailableAddressesCmd struct {
		out               io.Writer
		Client            kubernetes.Interface
		Name              string
		Namespace         string
		ThresholdWarning  string
		ThresholdCritical string
	}
)

func newCheckEndpointsAvailableAddressesCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkEndpointsAvailableAddressesCmd{out: out}

	cmd := &cobra.Command{
		Use:          "availableAddresses",
		Short:        "check if a k8s deployment has a minimum of available replicas",
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

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the endpoint")
	cmd.Flags().StringVarP(&c.ThresholdWarning, "warning", "w", "2:", "warning threshold for minimum available addresses")
	cmd.Flags().StringVarP(&c.ThresholdCritical, "critical", "c", "1:", "critical threshold for minimum available addresses")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkEndpointsAvailableAddressesCmd) run() {
	checkEndpoints := endpoints.NewCheckEndpoints(c.Client, c.Name, c.Namespace)
	result := checkEndpoints.CheckAvailableAddresses(
		endpoints.CheckAvailableAddressesOptions{
			ThresholdWarning:  c.ThresholdWarning,
			ThresholdCritical: c.ThresholdCritical,
		})
	result.Exit()
}
