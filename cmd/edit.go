package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/utils"
)

var (
	editDescription = `Edit a configuration of a resource by file name. The resource name and kind must be specified. 
JSON and YAML formats are accepted.

Ocean CD api reference please visit https://docs.spot.io/api/#tag/Ocean-CD`
	editExamples = `For example files in json and yaml format please visit our repo https://github.com/spotinst/spot-oceancd-cli
and see the samples dir`

	editCmd = &cobra.Command{
		Use:     "edit (-f FILENAME)",
		Short:   "Edit a resource by file name",
		Long:    editDescription,
		Example: editExamples,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFlags()
		},
		Run: func(cmd *cobra.Command, args []string) {
			runEditCmd(context.Background())
		},
	}
)

func runEditCmd(ctx context.Context) {
	configHandler, err := utils.NewConfigHandler(fileToApply)
	if err != nil {
		fmt.Printf("Failed to edit resource - %s\n", err.Error())
		return
	}

	err = configHandler.Handle(ctx, editResource)
	if err != nil {
		fmt.Printf("Failed to edit resource - %s\n", err.Error())
	}
}

func editResource(ctx context.Context, resource map[string]interface{}) error {
	var resourceName string
	var entityType string
	var err error

	resourceToEdit := make(map[string]interface{})
	kind, isKindExist := resource["kind"]
	if isKindExist {
		entityType, err = utils.GetOceanCdEntityKindByName(kind.(string))
		if err != nil {
			return err
		}

		delete(resource, "kind")
		resourceName = resource["name"].(string)
		resourceToEdit[entityType] = resource
	} else if len(resource) == 1 {

		for key, value := range resource {
			entityType, err = utils.GetOceanCdEntityKindByName(key)
			if err != nil {
				return err
			}

			resourceName = value.(map[string]interface{})["name"].(string)
			resourceToEdit[entityType] = value
		}
	} else {
		return errors.New("error: Unknown resource type")
	}

	_, resourceErr := oceancd.GetEntity(ctx, entityType, resourceName)
	if resourceErr != nil {

		if resourceErr.Error() == "resource does not exist" {
			fmt.Printf("Failed to edit resource '%s/%s'. Resource doesn't exist.\n", entityType, resourceName)
			return nil
		}

		return resourceErr
	}

	err = oceancd.UpdateResource(ctx, entityType, resourceName, resourceToEdit)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated resource '%s/%s'\n", entityType, resourceName)
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	editCmd.Flags().StringVarP(&fileToApply, "file", "f", "", "manifest file with resource definition")
}
