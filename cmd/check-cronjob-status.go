package main

import (
	"fmt"
	"github.com/SedoTech/icinga-checks-library"
	"io"
	"k8s.io/client-go/kubernetes"
	cronjob "github.com/SedoTech/check-k8s/pkg/checks/cronjob"
	"github.com/SedoTech/check-k8s/pkg/environment"
	"github.com/SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
)

type (
	checkCronjobStatusCmd struct {
		out       io.Writer
		Client    kubernetes.Interface
		Name      string
		Namespace string
	}
)

func newCheckCronjobStatusCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkCronjobStatusCmd{out: out}

	cmd := &cobra.Command{
		Use:          "status",
		Short:        "check if a k8s cronjob is running properly",
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

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace of the cronjob")
	cmd.MarkPersistentFlagRequired("namespace")

	return cmd
}

func (c *checkCronjobStatusCmd) run() {
	checkCronjob := cronjob.NewCheckCronjob(c.Client, c.Name, c.Namespace)
	result := checkCronjob.CheckStatus(
		cronjob.CheckStatusOptions{})
	result.Exit()
}
