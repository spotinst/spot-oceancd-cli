package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	fp "path/filepath"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/utils"
)

var (
	applyDescription = `Apply a configuration to a resource by file name. The resource name and kind must be specified. 
This resource will be created if it doesn't exist yet.
JSON and YAML formats are accepted.

Ocean CD api reference please visit https://docs.spot.io/api/#tag/Ocean-CD`
	applyExamples = `For example files in json and yaml format please visit our repo https://github.com/spotinst/spot-oceancd-cli
and see the samples dir`

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
			return validateFlags()
		},
		Run: func(cmd *cobra.Command, args []string) {
			runApplyCmd(context.Background())
		},
	}
)

func runApplyCmd(ctx context.Context) {
	configHandler, err := utils.NewConfigHandler(fileToApply, utils.Options{})
	if err != nil {
		fmt.Printf("Failed to apply resource - %s\n", err.Error())
		return
	}

	err = configHandler.Handle(ctx, applyResource)
	if err != nil {
		fmt.Printf("Failed to apply resource - %s\n", err.Error())
	}
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

func validateFlags() error {
	if fileToApply == "" {
		fmt.Println("You must specify a file using -f")
		return errors.New("error: Required file not specified")
	}

	fileExtensionWithDot := fp.Ext(fileToApply)
	if err := utils.IsFileTypeSupported(fileExtensionWithDot); err != nil {
		return err
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
