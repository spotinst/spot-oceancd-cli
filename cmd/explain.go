/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"spot-oceancd-cli/pkg/oceancd/model"
	"spot-oceancd-cli/pkg/utils"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var (
	explainDescription = `List the fields for supported resources.
Currently refer to Ocean CD api documentation.
https://docs.spot.io/api/#tag/Ocean-CD`
	explainCmd = &cobra.Command{
		Use:   "explain RESOURCE",
		Short: "Get documentation for a resource",
		Long:  explainDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateExplainArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runExplainCmd(context.Background(), args)
		},
	}
)

func runExplainCmd(ctx context.Context, args []string) {
	resourceType := args[0]
	entityType, _ := utils.GetOceanCdEntityKindByName(resourceType)

	apiUrl := ""
	switch entityType {
	case model.VerificationProviderEntity:
		apiUrl = "https://docs.spot.io/api/#operation/OceanCDVerificationProviderCreate"
	case model.VerificationTemplateEntity:
		apiUrl = "https://docs.spot.io/api/#operation/OceanCDVerificationTemplateCreate"
	case model.StrategyEntity:
		apiUrl = "https://docs.spot.io/api/#operation/OceanCDStrategyCreate"
	case model.RolloutSpecEntity:
		apiUrl = "https://docs.spot.io/api/#operation/OceanCDRolloutSpecCreate"
	case model.ClusterEntity:
		apiUrl = "https://docs.spot.io/api/#operation/OceanCDClusterList"
	}

	fmt.Printf("To review %s fields plese visit %s", entityType, apiUrl)
}

func validateExplainArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("You must specify the type of resource to explain. Use \"oceancd api-resources\" for a complete list of supported resources.")
		return errors.New("error: Required resource not specified")
	} else if len(args) > 1 {
		fmt.Println("You can only specify one type of resource to explain. Use \"oceancd api-resources\" for a complete list of supported resources.")
		return errors.New("error: Too many resources")
	} else {
		entityType := args[0]
		_, err := utils.GetOceanCdEntityKindByName(entityType)
		if err != nil {
			fmt.Printf("Unknown resource type '%s'. Use \"oceancd api-resources\" for a complete list of supported resources.\n", entityType)
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(explainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
