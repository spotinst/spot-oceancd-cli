package main

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "oceancd",
		Short: "OceanCD cli",
		Long:  `Cli for creation and manaing oceancd deployment and verification for K8ss`,
	}

	rootCmd.AddCommand(NewListCommand())
	return rootCmd
}
