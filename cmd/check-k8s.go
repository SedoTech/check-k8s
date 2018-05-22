package main

import (
	"os"

	"github.com/benkeil/check-k8s/pkg/print"

	"github.com/benkeil/check-k8s/pkg/icinga"

	"github.com/benkeil/check-k8s/pkg/checks"
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
		newCheckServiceCmd(out),
		newCheckDeploymentCmd(out),
		newCheckPodCmd(out),
		newCheckEndpointCmd(out),
	)

	return cmd
}

func init() {
}

func exitServiceState(name string, state icinga.ServiceState, err error) {
	print.Printfln("%s:%s: %s", name, state, err)
	os.Exit(state.Ordinal())
}

func exitServiceCheckResult(result checks.ServiceCheckResult) {
	print.Printfln("%s: %s", result.Name(), result.State(), result.Message())
	os.Exit(result.State().Ordinal())
}

func exitServiceCheckResults(results checks.ServiceCheckResults) {
	for _, result := range results.All() {
		print.Printfln("%s:%s: %s", result.Name(), result.State(), result.Message())
	}
	os.Exit(results.CalculateServiceState().Ordinal())
}
