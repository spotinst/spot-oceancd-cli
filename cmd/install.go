package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const installScript = ""

func NewInstallCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "intall",
		Short: "install oceancd controller",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ListResources(cmd.Context(), args)
		},
	}
	cmd.PersistentFlags().StringP("url", "o", "", "manifest file with resource definition")
	pflag := cmd.PersistentFlags().Lookup("url")
	viper.BindPFlag("output", pflag)
	return cmd
}
