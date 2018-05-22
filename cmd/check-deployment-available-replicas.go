package main

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	checkDeploymentAvailableReplicasCmd struct {
		version string
		out     io.Writer
		name    string
		options *CheckDeploymentAvailableReplicasOptions
	}

	// CheckDeploymentAvailableReplicasOptions contains options and flags from the command line
	CheckDeploymentAvailableReplicasOptions struct {
		Namespace         string
		ThresholdWarning  string
		ThresholdCritical string
	}
)

// AddFlags binds flags to the given flagset.
func (c *CheckDeploymentAvailableReplicasOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&c.ThresholdCritical, "critical", "c", "2:", "minimum of replicas in spec")
	fs.StringVarP(&c.ThresholdWarning, "warning", "w", "2:", "minimum of available replicas")
}

func newCheckDeploymentAvailableReplicasCmd(out io.Writer) *cobra.Command {
	options := &CheckDeploymentAvailableReplicasOptions{}
	c := &checkDeploymentAvailableReplicasCmd{out: out, options: options}

	cmd := &cobra.Command{
		Use:          "availableReplicas",
		Short:        "check if a k8s deployment has enough available replicas",
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

	options.AddFlags(cmd.PersistentFlags())
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func (c *checkDeploymentAvailableReplicasCmd) run() {
	//checkDeployment, err := checks.NewCheckDeployment(settings, c.name, *c.options)
	//if err != nil {
	//	exitServiceState("NewCheckDeployment", icinga.ServiceStateUnknown, err)
	//}
	//exitServiceCheckResult(checkDeployment.CheckAvailableReplicas())
}
