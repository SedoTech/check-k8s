package main

import (
	"errors"
	"os"

	"github.com/benkeil/check-k8s/pkg/environment"
	"github.com/spf13/cobra"
)

var (
	globalUsage = `This program tests kubernetes resources`
	version     string
	settings    environment.EnvSettings
)

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "check-k8s",
		Short:        "check-k8s checks if a k8s resource is healthy",
		Long:         globalUsage,
		Version:      version,
		SilenceUsage: true,
	}

	settings.AddFlags(cmd.PersistentFlags())

	cmd.PersistentFlags().Parse(args)

	// set defaults from environment
	settings.Init(cmd.PersistentFlags())

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		// check commands
		newCheckDeploymentCmd(settings, out),
	)

	return cmd
}

// NameArgs returns an error if there are not exactly 1 arg containing the resource name.
func NameArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("resource name is required")
		}
		return nil
	}
}
