package main

import (
	"errors"
	"io"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/cmd/api"
	"github.com/spf13/cobra"
	"k8s.io/api/apps/v1"
)

type (
	checkDeploymentCmd struct {
		out                       io.Writer
		Deployment                *v1.Deployment
		Name                      string
		Namespace                 string
		AvailableReplicasWarning  string
		AvailableReplicasCritical string
	}
)

func newCheckDeploymentCmd(out io.Writer) *cobra.Command {
	c := &checkDeploymentCmd{out: out}

	cmd := &cobra.Command{
		Use:          "deployment",
		Short:        "check if a k8s deployment resource is healthy",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("deployment name is required")
			}
			c.Name = args[0]
			deployment, err := api.GetDeployment(settings, api.GetDeploymentOptions{Name: c.Name, Namespace: c.Namespace})
			if err != nil {
				return err
			}
			c.Deployment = deployment
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.run()
			return nil
		},
	}
	cmd.AddCommand(
		newCheckDeploymentAvailableReplicasCmd(out),
	)

	cmd.PersistentFlags().StringVarP(&c.Namespace, "namespace", "n", "", "the namespace where the deployment is")
	cmd.Flags().StringVar(&c.AvailableReplicasWarning, "availableReplicasWarning", "2:", "minimum of replicas in spec")
	cmd.Flags().StringVar(&c.AvailableReplicasCritical, "availableReplicasCritical", "2:", "minimum of available replicas")

	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentCmd) run() {
	//checkDeployment, err := checks.NewCheckDeployment(settings, c.name, *c.options)
	//if err != nil {
	//	exitServiceState("NewCheckDeployment", icinga.ServiceStateUnknown, err)
	//}
	////checkDeployment.PrintYaml()
	//exitServiceCheckResults(checkDeployment.CheckAll())
	print.Yaml(c.Deployment)
}
