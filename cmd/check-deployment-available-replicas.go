package main

import (
	"fmt"
	"io"

	icinga "github.com/benkeil/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/deployment"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
)

type (
	checkDeploymentAvailableReplicasCmd struct {
		out               io.Writer
		Client            kubernetes.Interface
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
	cmd.Flags().StringVarP(&c.ThresholdCritical, "critical", "c", "2:", "critical threshold for minimum available replicas")
	cmd.Flags().StringVarP(&c.ThresholdWarning, "warning", "w", "2:", "warning threshold for minimum available replicas")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentAvailableReplicasCmd) run() {
	checkDeployment := deployment.NewCheckDeployment(c.Client, c.Name, c.Namespace)
	result := checkDeployment.CheckAvailableReplicas(
		deployment.CheckAvailableReplicasOptions{
			ThresholdWarning:  c.ThresholdWarning,
			ThresholdCritical: c.ThresholdCritical,
		})
	result.Exit()
}
