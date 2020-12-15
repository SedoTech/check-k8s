package main

import (
	"fmt"
	"io"

	icinga "github.com/SedoTech/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"github.com/SedoTech/check-k8s/pkg/checks/statefulset"
	"github.com/SedoTech/check-k8s/pkg/environment"
	"github.com/SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
)

type (
	checkStatefulSetAvailableReplicasCmd struct {
		out               io.Writer
		Client            kubernetes.Interface
		Name              string
		Namespace         string
		ThresholdWarning  string
		ThresholdCritical string
	}
)

func newCheckStatefulSetAvailableReplicasCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkStatefulSetAvailableReplicasCmd{out: out}

	cmd := &cobra.Command{
		Use:          "availableReplicas",
		Short:        "check if a k8s StatefulSet has a minimum of available replicas",
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

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the StatefulSet")
	cmd.Flags().StringVarP(&c.ThresholdCritical, "critical", "c", "2:", "critical threshold for minimum available replicas")
	cmd.Flags().StringVarP(&c.ThresholdWarning, "warning", "w", "2:", "warning threshold for minimum available replicas")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkStatefulSetAvailableReplicasCmd) run() {
	checkStatefulSet := statefulset.NewCheckStatefulSet(c.Client, c.Name, c.Namespace)
	result := checkStatefulSet.CheckAvailableReplicas(
		statefulset.CheckAvailableReplicasOptions{
			ThresholdWarning:  c.ThresholdWarning,
			ThresholdCritical: c.ThresholdCritical,
		})
	result.Exit()
}
