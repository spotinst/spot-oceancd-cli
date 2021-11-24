package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			file, err := flags.GetString("file")
			if err != nil {
				return err
			}

			color.Green("service %v was created\n", s.Name)

			return nil

		},
	}

	cmd.PersistentFlags().StringP("file", "f", "", "manifest file with resource definition")
	return cmd
}
