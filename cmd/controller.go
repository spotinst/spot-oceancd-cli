/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// controllerCmd represents the controller command
var (
	controllerDescription = `Provides utilities for interacting with ocean cd controller.
To learn more about Ocean CD please visit https://docs.spot.io/ocean-cd/ocean-cd-overview`
	controllerCmd = &cobra.Command{
		Use:   "controller",
		Short: "Provides utilities for interacting with ocean cd controller",
		Long: controllerDescription,
	}
)

func init() {
	rootCmd.AddCommand(controllerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
