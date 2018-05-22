package main

import (
	"errors"
	"io"

	"git.i.sedorz.net/infrastructure/icinga/check-k8s/pkg/checks"
	"git.i.sedorz.net/infrastructure/icinga/check-k8s/pkg/icinga"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	checkDeploymentCmd struct {
		version string
		out     io.Writer
		name    string
		options *checks.CheckDeploymentOptions
	}

	// CheckDeploymentOptions contains options and flags from the command line
	CheckDeploymentOptions struct {
		Namespace                string
		MinimumSpecReplicas      int32
		MinimumAvailableReplicas int32
	}
)

// AddFlags binds flags to the given flagset.
func (c *CheckDeploymentOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.Namespace, "namespace", "n", "", "the namespace where the deployment is")
	fs.Int32Var(&c.MinimumSpecReplicas, "minimumSpecReplicas", 2, "minimum of replicas in spec")
	fs.Int32Var(&c.MinimumAvailableReplicas, "minimumAvailableReplicas", 2, "minimum of available replicas")
}

func newCheckDeploymentCmd(out io.Writer) *cobra.Command {
	options := &CheckDeploymentOptions{}
	c := &checkDeploymentCmd{out: out, options: options}

	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "check if a k8s deployment resource is healthy",
		//TraverseChildren: true,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("deployment name is required")
			}
			c.name = args[0]
			c.run()
			return nil
		},
	}
	cmd.AddCommand(
		newCheckDeploymentAvailableReplicasCmd(out),
	)

	options.AddFlags(cmd.PersistentFlags())
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentCmd) run() {
	checkDeployment, err := checks.NewCheckDeployment(settings, c.name, *c.options)
	if err != nil {
		exitServiceState("NewCheckDeployment", icinga.ServiceStateUnknown, err)
	}
	//checkDeployment.PrintYaml()
	exitServiceCheckResults(checkDeployment.CheckAll())
}
