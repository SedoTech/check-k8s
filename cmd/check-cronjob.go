package main

import (
	"fmt"
	"github.com/SedoTech/icinga-checks-library"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/kubernetes"
	cronjob "SedoTech/check-k8s/pkg/checks/cronjob"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"
)

type (
	checkCronjobCmd struct {
		out       io.Writer
		Client    kubernetes.Interface
		Name      string
		Namespace string
	}
)

func newCheckCronjobCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkCronjobCmd{out: out}

	cmd := &cobra.Command{
		Use:          "cronjob",
		Short:        "check if a k8s cronjob works as expected",
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
		newCheckCronjobStatusCmd(settings, out),
	)

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the cronjob")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkCronjobCmd) run() {
	checkCronjob := cronjob.NewCheckCronjob(c.Client, c.Name, c.Namespace)
	results := checkCronjob.CheckStatus(
		cronjob.CheckStatusOptions{},
	)
	results.Exit()
}
