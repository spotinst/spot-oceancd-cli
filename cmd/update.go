package main

import (
	"context"
	"errors"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	commands "github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/model"
)

func NewUpdateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			file, err := flags.GetString("file")

			if err != nil {
				return err
			}

			obj, err := commands.EntitySpecFromFile(file)
			if err != nil {
				return err
			}

			spec, err := commands.GetEntitySpec(obj)
			meta, ok := spec.(model.EntityMeta)
			if !ok {
				return errors.New("can't retrieve metadata from object")
			}
			err = commands.UpdateEntity(context.Background(), obj, meta.GetEntityKind(), meta.GetEntityName())
			if err != nil {
				color.Red("update fail for %v %v\n", meta.GetEntityKind(), meta.GetEntityName())
				return err
			}
			color.Green(" %v %v was updated \n", meta.GetEntityKind(), meta.GetEntityName())
			return nil

		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "manifest file with resource definition")
	return cmd
}
