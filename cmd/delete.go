package cmd

import (
	"context"
	"errors"
	"fmt"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/utils"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete oceancd resources (microservice, environment or replicaset)",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateDeleteArgs(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteCmd(context.Background(), args)
	},
}

func runDeleteCmd(ctx context.Context, args []string) {
	resourceType := args[0]
	resourceNames := args[1:]

	entityType, _ := utils.GetEntityKindByName(resourceType)

	for _, resourceName := range resourceNames {
		err := oceancd.DeleteEntity(ctx, entityType, resourceName)
		if err != nil {
			fmt.Printf("Failed to delete '%v/%v' - %s\n", entityType, resourceName, err.Error())
			return
		}

		fmt.Printf("Successfully deleted resource '%v/%v'\n", entityType, resourceName)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateDeleteArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("You must specify resource type and name.")
		return errors.New("error: Required arguments not specified")
	}

	if len(args) < 2 {
		fmt.Println("You must specify resource name.")
		return errors.New("error: Required argument not specified")
	}

	entityType := args[0]
	_, err := utils.GetOceanCdEntityKindByName(entityType)
	if err != nil {
		fmt.Printf("Unknown resource '%s'. Use \"oceancd api-resources\" for a complete list of supported resources.\n", entityType)
		return err
	}

	return nil
}
