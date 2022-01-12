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

func NewGetCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get  oceancd resource representations(microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {

			return HandleGetMultipeEntities(context.Background(), args)

		},
	}

	return cmd
}
func HandleGet(ctx context.Context, entityToGet string) error {

	params := strings.Split(entityToGet, "/")
	if len(params) != 2 {
		color.Red("get failed as entity %v  should in be in form type/name \n", entityToGet)

		return errors.New("invalid input, should be in form of type/name")
	}
	entityType := params[0]
	entityName := params[1]

	entityType, err := utils.GetEntityKindByName(entityType)
	if err != nil {
		return err
	}

	output, err := commands.GetEntity(context.Background(), entityType, entityName, "json")
	if err != nil {
		color.Red("get failed for %v %v\n", entityType, entityName)
		return err
	}
	color.Green("%v \n", output)
	return nil

}

func HandleWideOption(ctx context.Context) {}
func HandleGetMultipeEntities(ctx context.Context, args []string) error {

	for _, a := range args {
		err := HandleGet(ctx, a)
		if err != nil {
			color.Red("get failed for entity %s\n", a)
			return err
		}
	}
	return nil
}
