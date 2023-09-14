/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// operatorCmd represents the operator command
var (
	operatorDescription = `Provides utilities for interacting with ocean cd operator.
To learn more about Ocean CD please visit https://docs.spot.io/ocean-cd/ocean-cd-overview`
	operatorUse = "operator"
	operatorCmd = &cobra.Command{
		Use:   operatorUse,
		Short: "Provides utilities for interacting with ocean cd operator",
		Long:  operatorDescription,
	}
)

func init() {
	rootCmd.AddCommand(operatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
