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
	checkStatefulSetProbesDefinedCmd struct {
		out           io.Writer
		Client        kubernetes.Interface
		Name          string
		Namespace     string
		Result        string
		ProbesDefined []string
	}
)

func newCheckStatefulSetProbesDefinedCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkStatefulSetProbesDefinedCmd{out: out}

	cmd := &cobra.Command{
		Use:          "probesDefined",
		Short:        "check if a k8s StatefulSet has liveness and readiness probes defined",
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
	cmd.Flags().StringSliceVarP(&c.ProbesDefined, "string", "s", []string{}, "check only the defined containers, not all")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkStatefulSetProbesDefinedCmd) run() {
	checkStatefulSet := statefulset.NewCheckStatefulSet(c.Client, c.Name, c.Namespace)
	result := checkStatefulSet.CheckProbesDefined(
		statefulset.CheckProbesDefinedOptions{
			Result:        c.Result,
			ProbesDefined: c.ProbesDefined,
		})
	result.Exit()
}
