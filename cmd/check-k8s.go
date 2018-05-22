package main

import (
	"os"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/icinga-checks-library"

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
		newCheckDeploymentCmd(out),
		newCheckPodCmd(out),
		newCheckEndpointCmd(out),
	)

	return cmd
}

func init() {
}

func exitServiceState(name string, state icinga.Status, err error) {
	print.Printfln("%s:%s: %s", name, state, err)
	os.Exit(state.Ordinal())
}

func exitServiceCheckResult(result icinga.Result) {
	print.Printfln("%s: %s", result.Name(), result.Status(), result.Message())
	os.Exit(result.Status().Ordinal())
}

func exitServiceCheckResults(results icinga.Results) {
	for _, result := range results.All() {
		print.Printfln("%s:%s: %s", result.Name(), result.Status(), result.Message())
	}
	//os.Exit(results.CalculateServiceStatus().Ordinal())
}
