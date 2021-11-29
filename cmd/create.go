package main

import (
	"errors"
	"reflect"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	commands "github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/utils"
)

func NewCreateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			file, err := flags.GetString("file")
			entityType := utils.GetEntityKind(file)

			if err != nil {
				return err
			}
			var obj interface{}
			switch entityType {
			case "environment":
				obj, err = commands.CreateEnvironmentFromManifest(file, true)

			case "service":
				obj, err = commands.CreateServiceFromManifest(file, true)

			case "rolloutspec":
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

		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "manifest file with resource definition")
	return cmd
}
