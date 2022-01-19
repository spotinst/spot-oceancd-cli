package main

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/model"
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
	cmd.PersistentFlags().StringP("output", "o", "", "manifest file with resource definition")
	pflag := cmd.PersistentFlags().Lookup("output")
	viper.BindPFlag("output", pflag)
	return cmd
}

func ListWithWideFlag(ctx context.Context, entityType string) error {

	if entityType != model.RolloutSpecEntity {
		return errors.New("flag output=wide currently supported only for rolloutSpec entity")
	}
	items, err := cmd.ListEntities(ctx, model.RolloutSpecEntity)
	if err != nil {
		return err
	}
	services, err := cmd.ListServices(ctx)

	if err != nil {
		return err
	}
	err = cmd.OutputEntitiesWide(entityType, items, services)

	return nil

}
func ListResources(ctx context.Context, args []string) error {

	go utils.WaitSpinner()

	entityType, err := utils.GetEntityKindByName(args[0])

	if err != nil {
		return err
	}

	if viper.Get("output") == "wide" {
		return ListWithWideFlag(ctx, entityType)
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
