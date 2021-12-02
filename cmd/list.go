package main

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/utils"
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

	go WaitSpinner()

	entityType, err := utils.GetEntityKindByName(args[0])

	if err != nil {
		return err
	}

	items, err := cmd.ListEntities(ctx, entityType)

	if err != nil {
		return err
	}

	err = cmd.OutputEntities(entityType, items)
	if err != nil {
		return err
	}

	return nil

}
