package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks/endpoints"
	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	icinga "github.com/benkeil/icinga-checks-library"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

type (
	checkEndpointsCmd struct {
		out                        io.Writer
		Client                     kubernetes.Interface
		Name                       string
		Namespace                  string
		AvailableAddressesWarning  string
		AvailableAddressesCritical string
	}
)

func newCheckEndpointsCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkEndpointsCmd{out: out}

	cmd := &cobra.Command{
		Use:          "endpoint",
		Short:        "check if a k8s endpoint resource is healthy",
		SilenceUsage: false,
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
	cmd.AddCommand(
		newCheckEndpointsAvailableAddressesCmd(settings, out),
	)

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the endpoint")
	cmd.Flags().StringVar(&c.AvailableAddressesWarning, "availableAddressesWarning", "2:", "warning threshold for minimum available replicas")
	cmd.Flags().StringVar(&c.AvailableAddressesCritical, "availableAdressesCritical", "1:", "critical threshold for minimum available replicas")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkEndpointsCmd) run() {
	checkEndpoints := endpoints.NewCheckEndpoints(c.Client, c.Name, c.Namespace)
	results := checkEndpoints.CheckAll(endpoints.CheckAllOptions{
		CheckAvailableAddressesOptions: endpoints.CheckAvailableAddressesOptions{
			ThresholdWarning:  c.AvailableAddressesWarning,
			ThresholdCritical: c.AvailableAddressesCritical,
		},
	})
	results.Exit()
}
