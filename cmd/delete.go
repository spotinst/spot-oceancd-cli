package main

import (
	"context"
	"errors"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	commands "github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/utils"
)

func NewDeleteCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {

			return HandleDeleteMultipeEntities(context.Background(), args)

		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "manifest file with resource definition")
	return cmd
}
func HandleDelete(ctx context.Context, entityToDelete string) error {

	params := strings.Split(entityToDelete, "/")
	if len(params) != 2 {
		color.Red("delete failed as entity %v  should in be in form type/name \n", entityToDelete)

		return errors.New("invalid input, should be in form of type/name")
	}
	entityType := params[0]
	entityName := params[1]

	entityType, err := utils.GetEntityKindByName(entityType)
	if err != nil {
		return err
	}

	err = commands.DeleteEntity(context.Background(), entityType, entityName)
	if err != nil {
		color.Red("delete failed for %v %v\n", entityType, entityName)
		return err
	}
	color.Green("delete succeedeed for entity %v \n", entityToDelete)
	return nil

}
func HandleDeleteMultipeEntities(ctx context.Context, args []string) error {

	for _, a := range args {
		err := HandleDelete(ctx, a)
		if err != nil {
			color.Red("delete failed for entity %s\n", a)
			return err
		}
	}
	return nil
}
