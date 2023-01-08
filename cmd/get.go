package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lensesio/tableprinter"
	"github.com/spf13/cobra"
	"os"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model"
	"spot-oceancd-cli/pkg/utils"
)

var (
	output string

	getDescription = `Display one or many resources.
Prints a table of the most important information about the specified resources.`
	getExamples = `  # List all strategies in ps output format
  oceancd get stgs

  # List a single strategy with specified NAME in ps output format
  oceancd get stgs app-canary

  # List strategies in JSON output format
  oceancd get stgs -o json

  # List a single strategy in JSON output format
  oceancd get -o json stgs app-canary

  # List a single strategy in YAML output format
  oceancd get -o yml stgs app-canary`

	getCmd = &cobra.Command{
		Use:     "get [(-o|--output=)json|yaml|yml|wide] (TYPE [NAME] ...) [flags]",
		Short:   "Display one or many resources",
		Long:    getDescription,
		Example: getExamples,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateGetArgs(cmd, args)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
		},
		Run: func(cmd *cobra.Command, args []string) {
			runGetCmd(context.Background(), args)
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringVarP(&output, "output", "o", "wide", "Output format. One of: json|yaml|wide")
}

func validateGetArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("You must specify the type of resource to get. Use \"oceancd api-resources\" for a complete list of supported resources.")
		return errors.New("error: Required resource not specified")
	} else {
		entityType := args[0]
		_, err := utils.GetEntityKindByName(entityType)
		if err != nil {
			fmt.Printf("Unknown resource type '%s'. Use \"oceancd api-resources\" for a complete list of supported resources.\n", entityType)
			return err
		}
	}

	return nil
}

func runGetCmd(ctx context.Context, args []string) {
	var resources []interface{}
	resourceType := args[0]
	resourceNames := args[1:]

	entityType, err := utils.GetEntityKindByName(resourceType)
	if len(args) == 1 {
		resources, err = oceancd.ListEntities(context.Background(), entityType)
		if err != nil {
			fmt.Printf("Failed to get resource '%s' - %s\n", entityType, err.Error())
			return
		}
	} else {
		for _, resourceName := range resourceNames {
			resource, getErr := oceancd.GetEntity(context.Background(), entityType, resourceName)
			if getErr != nil {
				fmt.Printf("Failed to get resource '%s/%s' - %s\n", entityType, resourceName, getErr.Error())
				return
			}

			resources = append(resources, resource)
		}
	}

	switch output {
	case "yaml", "yml":
		resourcesStr, yamlErr := utils.ConvertEntitiesToYamlString(resources)
		if yamlErr != nil {
			fmt.Printf("Failed to convert resources to yaml - %s\n", yamlErr.Error())
		}
		fmt.Println(resourcesStr)
	case "json":
		resourcesStr, jsonErr := utils.ConvertEntitiesToJsonString(resources)
		if jsonErr != nil {
			fmt.Printf("Failed to convert resources to json - %s\n", jsonErr.Error())
		}
		fmt.Println(resourcesStr)
	case "wide":
		handlePrint(ctx, entityType, resources)
	case "":
		if len(resources) == 1 {
			resourcesStr, yamlErr := utils.ConvertEntitiesToYamlString(resources)
			if yamlErr != nil {
				fmt.Printf("Failed to convert resources to yaml - %s\n", yamlErr.Error())
			}
			fmt.Println(resourcesStr)
		} else {
			handlePrint(ctx, entityType, resources)
		}
	default:
		fmt.Printf("Unknown output '%s'. Please choose one of: json|yaml|wide\n", output)
	}

	return
}

func handlePrint(ctx context.Context, entityType string, resources []interface{}) {
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = false, false, false, false
	printer.CenterSeparator = " "
	printer.ColumnSeparator = " "
	printer.RowSeparator = " "

	switch entityType {
	case model.VerificationProviderEntity:
		entitiesDetails := utils.GetVerificationProviderEntitiesDetails(resources)
		printer.Print(entitiesDetails)
	case model.VerificationTemplateEntity:
		entitiesDetails := utils.GetVerificationTemplateEntitiesDetails(resources)
		printer.Print(entitiesDetails)
	case model.StrategyEntity:
		entitiesDetails := utils.GetStrategyEntitiesDetails(resources)
		printer.Print(entitiesDetails)
	case model.RolloutSpecEntity:
		entitiesDetails := utils.GetRolloutSpecEntitiesDetails(resources)
		printer.Print(entitiesDetails)
	case model.ClusterEntity:
		entitiesDetails := utils.GetClusterEntitiesDetails(resources)
		printer.Print(entitiesDetails)
	}
}
