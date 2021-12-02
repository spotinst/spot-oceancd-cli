package main

import (
	"context"
	"errors"
	"reflect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	commands "github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

func NewCreateCommand() *cobra.Command {
	filesVar := FilesFlag{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {

			return HandleMultipeFiles(context.Background(), filesVar, args)

		},
	}

	//cmd.PersistentFlags().StringP("file", "f", "", "manifest file with resource definition")
	cmd.PersistentFlags().VarP(&filesVar, "file", "f", "manifest file with resource definition")
	return cmd
}
func HandleCreationFromfile(ctx context.Context, file string) error {
	entityType := utils.GetEntityKindByFilename(file)

	var obj interface{}
	var err error
	switch entityType {
	case model.EnvEntity:
		obj, err = commands.CreateEnvironmentFromManifest(file, true)

	case model.ServiceEntity:
		obj, err = commands.CreateServiceFromManifest(file, true)

	case model.RolloutSpecEntity:
		obj, err = commands.CreateRolloutFromManifest(file, true)
	default:
		return errors.New("make sure that file name matches one of supported entiteis(environment, service, rolloutspec")

	}
	if err != nil {
		return err
	}
	val := reflect.ValueOf(obj).Elem()
	nameField := val.FieldByName("Name")
	f := nameField.Interface()
	objName := reflect.ValueOf(f)

	color.Green(" %v %v was created \n", entityType, objName)

	return nil
}
func HandleMultipeFiles(ctx context.Context, f FilesFlag, arg []string) error {

	for _, file := range f.files {
		err := HandleCreationFromfile(ctx, file)
		if err != nil {
			color.Red("creation failed for manifest %s\n", file)
			return err
		}
	}
	return nil
}
