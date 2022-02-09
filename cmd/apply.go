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
	applyDescription = `Apply a configuration to a resource by file name. The resource name and kind must be specified. This resource will be
created if it doesn't exist yet.
JSON and YAML formats are accepted.

Ocean CD api reference please visit https://docs.spot.io/api/#tag/Ocean-CD`
	applyExamples = `For example files in json and yaml format
please visit our repo https://github.com/spotinst/spot-oceancd-cli
and see the samples dir`

	supportedFileTypes = map[string]bool{
		"json": true,
		"yml":  true,
		"yaml": true,
	}

	fileToApply string

	applyCmd = &cobra.Command{
		Use:     "apply (-f FILENAME)",
		Short:   "Apply a configuration to a resource by file name",
		Long:    applyDescription,
		Example: applyExamples,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFlags(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runApplyCmd(context.Background())
		},
	}
)

func runApplyCmd(ctx context.Context) {
	fileExtension := filepath.Ext(fileToApply)[1:]

	switch fileExtension {
	case "json":
		err := handleJsonFile(ctx)
		if err != nil {
			fmt.Printf("Failed to apply resource - %s\n", err.Error())
		}
	case "yaml", "yml":
		err := handleYamlFile(ctx)
		if err != nil {
			fmt.Printf("Failed to apply resource - %s\n", err.Error())
		}
	}
}

func handleYamlFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertYamlFileToArrayOfMaps(fileToApply)
	if err != nil {
		resources, err = utils.ConvertYamlFileToMap(fileToApply)
		if err != nil {
			return err
		}
	}

	for _, resource = range resources {
		err = applyResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleJsonFile(ctx context.Context) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = utils.ConvertJsonFileToArrayOfMaps(fileToApply)
	if err != nil {
		resource, err = utils.ConvertJsonFileToMap(fileToApply)
		if err != nil {
			return err
		}

		return applyResource(ctx, resource)
	}

	for _, resource = range resources {
		err = applyResource(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyResource(ctx context.Context, resource map[string]interface{}) error {
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
		return errors.New("error: Unknown resource type")
	}

	_, resourceErr := oceancd.GetEntity(ctx, entityType, resourceName)
	if resourceErr != nil {
		if resourceErr.Error() == "resource does not exist" {
			err = oceancd.CreateResource(ctx, entityType, resourceToApply)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully created resource '%s/%s'\n", entityType, resourceName)
			return nil
		}

		return resourceErr
	}

	err = oceancd.UpdateResource(ctx, entityType, resourceName, resourceToApply)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated resource '%s/%s'\n", entityType, resourceName)
	return nil
}

func validateFlags(cmd *cobra.Command, args []string) error {
	if fileToApply == "" {
		fmt.Println("You must specify a file using -f")
		return errors.New("error: Required file not specified")
	}

	fileExtensionWithDot := filepath.Ext(fileToApply)
	if fileExtensionWithDot == "" {
		fmt.Println("File must have an extension of type json or yaml")
		return errors.New("error: Unsupported file type")
	}

	fileExtension := fileExtensionWithDot[1:]
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
	applyCmd.Flags().StringVarP(&fileToApply, "file", "f", "", "manifest file with resource definition")
}
