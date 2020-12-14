package main

import (
	"fmt"
	"io"

	"github.com/benkeil/check-k8s/pkg/checks/statefulset"
	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/benkeil/check-k8s/pkg/kube"
	icinga "github.com/benkeil/icinga-checks-library"
	"k8s.io/client-go/kubernetes"

	"github.com/spf13/cobra"
)

type (
	checkStatefulSetContainerDefinedCmd struct {
		out              io.Writer
		Client           kubernetes.Interface
		Name             string
		Namespace        string
		Result           string
		ContainerDefined []string
	}
)

func newCheckStatefulSetContainerDefinedCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkStatefulSetContainerDefinedCmd{out: out}

	cmd := &cobra.Command{
		Use:          "containerDefined",
		Short:        "check if a k8s StatefulSet has a list of containers defined",
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
	cmd.Flags().StringVarP(&c.Result, "result", "r", "CRITICAL", "the result state if the check fails")
	cmd.Flags().StringSliceVarP(&c.ContainerDefined, "string", "s", []string{}, "the containers to check if they are present (comma separated list)")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkStatefulSetContainerDefinedCmd) run() {
	checkStatefulSet := statefulset.NewCheckStatefulSet(c.Client, c.Name, c.Namespace)
	result := checkStatefulSet.CheckContainerDefined(
		statefulset.CheckContainerDefinedOptions{
			Result:           c.Result,
			ContainerDefined: c.ContainerDefined,
		})
	result.Exit()
}
