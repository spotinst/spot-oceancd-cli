/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

// versionCmd represents the version command
var (
	shortened     = false
	versionOutput = "yaml"
	versionCmd    = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, versionOutput)
			fmt.Print(resp)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", true, "Print just the version number.")
	versionCmd.Flags().StringVarP(&versionOutput, "output", "o", "yaml", "Output format. One of 'yaml' or 'json'.")
}
