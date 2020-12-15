package main

import (
	"fmt"
	"io"

	icinga "github.com/SedoTech/icinga-checks-library"
	"github.com/SedoTech/check-k8s/pkg/checks/statefulset"
	"github.com/SedoTech/check-k8s/pkg/environment"
	"github.com/SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

type (
	checkStatefulSetCmd struct {
		out                       io.Writer
		Client                    kubernetes.Interface
		Name                      string
		Namespace                 string
		AvailableReplicasWarning  string
		AvailableReplicasCritical string
		UpdateStrategyResult      string
		UpdateStrategyValue       string
		PodRestartsResult         string
		PodRestartsDuration       string
		ContainerDefinedResult    string
		ContainerDefinedValue     []string
		ProbesDefinedResult       string
		ProbesDefinedValue        []string
		Test                      []string
	}
)

func newCheckStatefulSetCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkStatefulSetCmd{out: out}

	cmd := &cobra.Command{
		Use:          "statefulset",
		Short:        "check if a k8s StatefulSet resource is healthy",
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
		newCheckStatefulSetAvailableReplicasCmd(settings, out),
		newCheckStatefulSetUpdateStrategyCmd(settings, out),
		newCheckStatefulSetPodRestartsCmd(settings, out),
		newCheckStatefulSetProbesDefinedCmd(settings, out),
		newCheckStatefulSetContainerDefinedCmd(settings, out),
	)

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the StatefulSet")
	cmd.Flags().StringVar(&c.AvailableReplicasWarning, "availableReplicasWarning", "2:", "warning threshold for minimum available replicas")
	cmd.Flags().StringVar(&c.AvailableReplicasCritical, "availableReplicasCritical", "2:", "critical threshold for minimum available replicas")
	cmd.Flags().StringVar(&c.UpdateStrategyResult, "updateStrategyResult", "WARNING", "the result state if the updateStrategy check fails")
	cmd.Flags().StringVar(&c.UpdateStrategyValue, "updateStrategyValue", "RollingUpdate", "the expected update strategy")
	cmd.Flags().StringVar(&c.PodRestartsResult, "podRestartsResult", "WARNING", "the result state if the podRestart check fails")
	cmd.Flags().StringVar(&c.PodRestartsDuration, "podRestartsDuration", "15m", "the duration during the check looks for restarts")
	cmd.Flags().StringVar(&c.ContainerDefinedResult, "containerDefinedResult", "CRITICAL", "the result state if the updateStrategy check fails")
	cmd.Flags().StringSliceVar(&c.ContainerDefinedValue, "containerDefinedValue", []string{}, "check only the defined containers, not all")
	cmd.Flags().StringVar(&c.ProbesDefinedResult, "probesDefinedResult", "WARNING", "the result state if the updateStrategy check fails")
	cmd.Flags().StringSliceVar(&c.ProbesDefinedValue, "probesDefinedContainer", []string{}, "check only the defined containers, not all")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkStatefulSetCmd) run() {
	checkStatefulSet := statefulset.NewCheckStatefulSet(c.Client, c.Name, c.Namespace)
	results := checkStatefulSet.CheckAll(statefulset.CheckAllOptions{
		CheckAvailableReplicasOptions: statefulset.CheckAvailableReplicasOptions{
			ThresholdWarning:  c.AvailableReplicasWarning,
			ThresholdCritical: c.AvailableReplicasCritical,
		},
		CheckUpdateStrategyOptions: statefulset.CheckUpdateStrategyOptions{
			Result:         c.UpdateStrategyResult,
			UpdateStrategy: c.UpdateStrategyValue,
		},
		CheckPodRestartsOptions: statefulset.CheckPodRestartsOptions{
			Result:   c.PodRestartsResult,
			Duration: c.PodRestartsDuration,
		},
		CheckContainerDefinedOptions: statefulset.CheckContainerDefinedOptions{
			Result:           c.ContainerDefinedResult,
			ContainerDefined: c.ContainerDefinedValue,
		},
		CheckProbesDefinedOptions: statefulset.CheckProbesDefinedOptions{
			Result:        c.ProbesDefinedResult,
			ProbesDefined: c.ProbesDefinedValue,
		},
	})
	results.Exit()
}
