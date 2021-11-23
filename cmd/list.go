package main

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/verchol/applier/pkg/cmd"
)

const (
	installationNamespace = "oceancd"
	namixServer           = "oceancd"
)

func NewListCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list oceancd resources (microservices,  environmetns , etc",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ListResources(cmd.Context(), args)
		},
	}

	return cmd
}
func WaitSpinner() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner
	time.Sleep(4 * time.Second)                                 // Run for some time to simulate work
	s.Stop()
}
func ListResources(ctx context.Context, args []string) error {
	if args[0] == "services" {
		go WaitSpinner()
		cmd.ListServices(ctx)
	}

	if args[0] == "rolloutspec" || args[0] == "rs" {
		go WaitSpinner()
		cmd.ListRolloutSpecs(ctx)
	}

	return nil
}
