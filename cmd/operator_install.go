/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"spot-oceancd-cli/pkg/utils"
)

var operatorManagerConfig string

// operatorInstallCmd represents the operator install command
var (
	operatorInstallDescription      = `Installs Ocean CD operator on current cluster with dependencies based on provided config.`
	operatorInstallShortDescription = "Installs Ocean CD operator on current cluster"
	operatorInstallUse              = "install"
	operatorInstallExample          = fmt.Sprintf("  # %s\n  %s %s %s %s",
		operatorInstallShortDescription, rootCmd.Name(), operatorUse, operatorInstallUse, "--config /path/to/config")

	operatorInstallCmd = &cobra.Command{
		Use:     operatorInstallUse,
		Short:   operatorInstallShortDescription,
		Long:    operatorInstallDescription,
		Example: operatorInstallExample,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorInstallFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runOperatorInstallCmd(context.Background())
		},
	}
)

func init() {
	operatorCmd.AddCommand(operatorInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	operatorInstallCmd.Flags().StringVarP(&operatorManagerConfig, "config", "c", "",
		"The configuration applied to OceanCD resources and their dependencies.")

}

func runOperatorInstallCmd(ctx context.Context) {

}

func validateOperatorInstallFlags(cmd *cobra.Command) error {
	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse --config flag: %w", err)
	}

	if cmd.Flags().Lookup("config").Changed == false {
		return nil
	}

	if pathToConfig == "" {
		return fmt.Errorf("path to config file must be specified")
	}

	fileExtensionWithDot := filepath.Ext(pathToConfig)
	if err := utils.IsFileTypeSupported(fileExtensionWithDot); err != nil {
		return err
	}

	return nil
}
