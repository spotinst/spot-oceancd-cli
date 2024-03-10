/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
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
			validateClusterId(context.Background())
			validateClusterIdExists(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorInstallFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			isOperatorInstallCommand = false
			shouldCreateNamespace = false

			fmt.Printf("Upgrading OceanCD operator manager in cluster %s\n", clusterId)

			if err := runOperatorInstallCmd(context.Background(), cmd); err != nil {
				fmt.Printf("Failed to upgrade operator\n%s\n", err)
			}

			fmt.Printf("Upgrade of OceanCD operator manager finished succesfully.\n")
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
