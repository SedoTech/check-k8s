package main

import (
	"fmt"
	"io"

	icinga "github.com/SedoTech/icinga-checks-library"
	"k8s.io/client-go/kubernetes"
	"SedoTech/check-k8s/pkg/checks/statefulset"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
)

type (
	checkStatefulSetPodRestartsCmd struct {
		out       io.Writer
		Client    kubernetes.Interface
		Name      string
		Namespace string
		Duration  string
		Result    string
	}
)

func newCheckStatefulSetPodRestartsCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkStatefulSetPodRestartsCmd{out: out}

	cmd := &cobra.Command{
		Use:          "podRestarts DEPLOYMENT",
		Short:        "check if a k8s StatefulSet has no pod restarts",
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
	cmd.Flags().StringVarP(&c.Result, "result", "r", "WARNING", "the result state if the check fails")
	cmd.Flags().StringVarP(&c.Duration, "duration", "d", "15m", "which duration we want to check for restarts")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkStatefulSetPodRestartsCmd) run() {
	checkStatefulSet := statefulset.NewCheckStatefulSet(c.Client, c.Name, c.Namespace)
	result := checkStatefulSet.CheckPodRestarts(
		statefulset.CheckPodRestartsOptions{
			Duration: c.Duration,
			Result:   c.Result,
		})
	result.Exit()
}
