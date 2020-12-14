package main

import (
	"fmt"
	"io"

	icinga "github.com/benkeil/icinga-checks-library"
	"SedoTech/check-k8s/pkg/checks/secrets"
	"SedoTech/check-k8s/pkg/environment"
	"SedoTech/check-k8s/pkg/kube"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

type (
	checkSecretsCmd struct {
		out            io.Writer
		Client         kubernetes.Interface
		Namespace      string
		SecretsDefined []string
	}
)

func newCheckSecretsCmd(settings environment.EnvSettings, out io.Writer) *cobra.Command {
	c := &checkSecretsCmd{out: out}

	cmd := &cobra.Command{
		Use:          "secrets",
		Short:        "check if a k8s secret resource exists",
		SilenceUsage: false,
		Args:         cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
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
	cmd.Flags().StringSliceVarP(&c.SecretsDefined, "string", "s", []string{}, "the secrets to check if they are exists (comma separated list)")
	cmd.MarkPersistentFlagRequired("namespace")
	cmd.MarkFlagRequired("string")

	return cmd
}

func (c *checkSecretsCmd) run() {
	checkSecrets := secrets.NewCheckSecrets(c.Client, c.Namespace)
	result := checkSecrets.CheckExists(
		secrets.CheckExistsOptions{
			SecretsDefined: c.SecretsDefined,
		},
	)
	result.Exit()
}
