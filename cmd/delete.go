package cmd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/utils"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var (
	deleteDescription = `Delete resources by file names or resource and names.

JSON and YAML formats are accepted. Only one type of argument may be specified: file names or resource and names`
	deleteExamples = `  # Delete a strategy using the type and name specified in strategy.json
  oceancd delete -f ./strategy.json

  # Delete strategies with names "baz" and "foo"
  oceancd delete stg baz foo`
	fileTolDelete string

	deleteCmd = &cobra.Command{
		Use:     "delete ([-f FILENAME] | TYPE [(NAME)])",
		Short:   "Delete resources by file names or resource and names",
		Long:    deleteDescription,
		Example: deleteExamples,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return validateDeleteArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runDeleteCmd(context.Background(), args)
		},
	}
)

func runDeleteCmd(ctx context.Context, args []string) {
	if fileTolDelete != "" {
		handleDeleteByFile(ctx)
		return
	}

	handleDeleteByArgs(ctx, args)
}

func handleDeleteByArgs(ctx context.Context, args []string) {
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

func handleDeleteByFile(ctx context.Context) {
	configHandler, err := utils.NewConfigHandler(utils.Options{PathToConfig: fileTolDelete})
	if err != nil {
		fmt.Printf("Failed to delete resource - %s\n", err.Error())
		return
	}

	err = configHandler.Handle(ctx, deleteResource)
	if err != nil {
		fmt.Printf("Failed to delete resource - %s\n", err.Error())
	}
}

func deleteResource(ctx context.Context, resource map[string]interface{}) error {
	var resourceName string
	var entityType string
	var err error

	kind, isKindExist := resource["kind"]
	if isKindExist {
		entityType, err = utils.GetOceanCdEntityKindByName(kind.(string))
		if err != nil {
			return err
		}

		resourceName = resource["name"].(string)
	} else if len(resource) == 1 {

		for key, value := range resource {
			entityType, err = utils.GetOceanCdEntityKindByName(key)
			if err != nil {
				return err
			}

			resourceName = value.(map[string]interface{})["name"].(string)
		}
	} else {
		return errors.New("error: Unknown resource type")
	}

	err = oceancd.DeleteEntity(ctx, entityType, resourceName)
	if err != nil {
		err = errors.New(fmt.Sprintf("'%v/%v' - %s\n", entityType, resourceName, err.Error()))
		return err
	}

	fmt.Printf("Successfully deleted resource '%v/%v'\n", entityType, resourceName)

	return nil
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

	deleteCmd.Flags().StringVarP(&fileTolDelete, "file", "f", "", "manifest file with resource definition")
}

func validateDeleteArgs(cmd *cobra.Command, args []string) error {
	if fileTolDelete != "" {
		fileExtensionWithDot := filepath.Ext(fileTolDelete)

		if err := utils.IsFileTypeSupported(fileExtensionWithDot); err != nil {
			return err
		}

		return nil
	}

	if len(args) < 1 {
		fmt.Println("You must specify either filename or resource type and name.")
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
