package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"path/filepath"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/utils"
)

var (
	supportedFileTypes = map[string]bool {
		"json": true,
		"yml": true,
		"yaml": true,
	}

	file string

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply a configuration to oceancd resources by file (microservices,  environments, replicasets)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFlags(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runApplyCmd(context.Background())
		},
	}
)

func runApplyCmd(ctx context.Context) {
	fileExtension := filepath.Ext(file)[1:]

	switch fileExtension {
	case "json":
		err := HandleJsonFile(ctx)
		if err != nil {
			fmt.Printf("Failed to apply resource - %s\n", err.Error())
		}
	case "yaml", "yml":
		err := HandleYamlFile(ctx)
		if err != nil {
			fmt.Printf("Failed to apply resource - %s\n", err.Error())
		}
	}
}

func HandleYamlFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertYamlFileToArrayOfMaps(file)
	if err != nil {
		resource, err = utils.ConvertYamlFileToMap(file)
		if err != nil {
			return err
		}

		return ApplyResource(ctx, resource)
	}

	for _, resource = range resources {
		err = ApplyResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleJsonFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertJsonFileToArrayOfMaps(file)
	if err != nil {
		resource, err = utils.ConvertJsonFileToMap(file)
		if err != nil {
			return err
		}

		return ApplyResource(ctx, resource)
	}

	for _, resource = range resources {
		err = ApplyResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyResource(ctx context.Context, resource map[string]interface{}) error {
	var resourceName string
	var entityType string
	var err error

	resourceToApply := make(map[string]interface{})
	kind, isKindExist := resource["kind"]
	if isKindExist {
		entityType, err = utils.GetOceanCdEntityKindByName(kind.(string))
		if err != nil {
			return err
		}

		delete(resource, "kind")
		resourceName = resource["name"].(string)
		resourceToApply[entityType] = resource
	} else if len(resource) == 1 {

		for key, value := range resource {
			entityType, err = utils.GetOceanCdEntityKindByName(key)
			if err != nil {
				return err
			}

			resourceName = value.(map[string]interface{})["name"].(string)
			resourceToApply[entityType] = value
		}
	} else {
		return errors.New(fmt.Sprintf("error: Unknown resource type '%s'", kind))
	}

	_, resourceErr := oceancd.GetEntity(ctx, entityType, resourceName)
	if resourceErr != nil {
		if resourceErr.Error() == "entity does not exist" {
			err = oceancd.CreateResource(ctx, entityType, resourceToApply)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully applied resource '%s/%s'\n", entityType, resourceName)
			return nil
		}

		return resourceErr
	}

	err = oceancd.UpdateResource(ctx, entityType, resourceName, resourceToApply)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully applied resource '%s/%s'\n", entityType, resourceName)
	return nil
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if file == "" {
		fmt.Println("You must specify a file using -f")
		return errors.New("error: Required file not specified")
	}

	fileExtension := filepath.Ext(file)[1:]

	if supportedFileTypes[fileExtension] == false {
		fmt.Println("File must be of type json or yaml")
		return errors.New("error: Unsupported file type")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "manifest file with resource definition")
}
