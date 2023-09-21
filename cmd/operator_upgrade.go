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

// operatorUpgradeCmd represents the operator upgrade command
var (
	operatorUpgradeDescription = `Upgrades Ocean CD operator based on provided config.`
	operatorUpgradeUse         = "upgrade"
	operatorUpgradeExample     = fmt.Sprintf("  # %s\n  %s %s %s %s",
		operatorUpgradeDescription, rootCmd.Name(), operatorUse, operatorUpgradeUse, "--config /path/to/config")

	operatorUpgradeCmd = &cobra.Command{
		Use:     operatorUpgradeUse,
		Short:   operatorUpgradeDescription,
		Long:    operatorUpgradeDescription,
		Example: operatorUpgradeExample,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorUpgradeFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runOperatorInstallCmd(context.Background(), cmd); err != nil {
				fmt.Printf("failed to upgrade operator: %s\n", err)
			}
		},
	}
)

func init() {
	operatorCmd.AddCommand(operatorUpgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	operatorUpgradeCmd.Flags().StringVarP(&operatorManagerConfig, "config", "c", "",
		"The configuration applied to OceanCD resources and their dependencies.")

}

func validateOperatorUpgradeFlags(cmd *cobra.Command) error {
	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse --config flag: %w", err)
	}

	if cmd.Flags().Lookup("config").Changed == false {
		return fmt.Errorf("--config flag using is required")
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
