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
	fileTolDelete string

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete oceancd resources (microservices, environments or rolloutspecs)",
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
	fileExtension := filepath.Ext(fileTolDelete)[1:]

	switch fileExtension {
	case "json":
		err := handleDeleteByJsonFile(ctx)
		if err != nil {
			fmt.Printf("Failed to delete resource - %s\n", err.Error())
		}
	case "yaml", "yml":
		err := handleDeleteByYamlFile(ctx)
		if err != nil {
			fmt.Printf("Failed to delete resource - %s\n", err.Error())
		}
	}
}

func handleDeleteByYamlFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertYamlFileToArrayOfMaps(fileTolDelete)
	if err != nil {
		resources, err = utils.ConvertYamlFileToMap(fileTolDelete)
		if err != nil {
			return err
		}
	}

	for _, resource = range resources {
		err = deleteResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleDeleteByJsonFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertJsonFileToArrayOfMaps(fileTolDelete)
	if err != nil {
		resource, err = utils.ConvertJsonFileToMap(fileTolDelete)
		if err != nil {
			return err
		}

		return deleteResource(ctx, resource)
	}

	for _, resource = range resources {
		err = deleteResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
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
		fileExtension := filepath.Ext(fileTolDelete)[1:]

		if err := utils.IsFileTypeSupported(fileExtension); err != nil {
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
