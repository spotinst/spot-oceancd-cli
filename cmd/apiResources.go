package cmd

import (
	"context"
	"github.com/lensesio/tableprinter"
	"github.com/spf13/cobra"
	"os"
)

type ApiResource struct {
	Name       string `header:"Name"`
	Shortnames string `header:"Shortnames"`
	Namespaced bool   `header:"Namespaced"`
	Kind       string `header:"Kind"`
}

var (
	apiResourcesDescription = `Print the supported API resources.
For full getting started tutorial please visit https://docs.spot.io/ocean-cd/getting-started/`
	apiResourcesCmd = &cobra.Command{
		Use:   "api-resources",
		Short: "Print the supported API resources",
		Long:  apiResourcesDescription,
		Run: func(cmd *cobra.Command, args []string) {
			printApiResources(context.Background())
		},
	}
)

func printApiResources(ctx context.Context) {
	apiResources := buildApiResourcesList(ctx)
	printer := tableprinter.New(os.Stdout)
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = false, false, false, false
	printer.CenterSeparator = " "
	printer.ColumnSeparator = " "
	printer.RowSeparator = " "
	printer.Print(apiResources)
}

func buildApiResourcesList(ctx context.Context) []ApiResource {
	return []ApiResource{
		{
			Name:       "environments",
			Shortnames: "env,envs",
			Namespaced: true,
			Kind:       "Environment",
		},
		{
			Name:       "microservices",
			Shortnames: "ms",
			Namespaced: false,
			Kind:       "Microservice",
		},
		{
			Name:       "notificationProviders",
			Shortnames: "np",
			Namespaced: false,
			Kind:       "notificationProvider",
		},
		{
			Name:       "rolloutSpecs",
			Shortnames: "rs",
			Namespaced: true,
			Kind:       "RolloutSpec",
		},
		{
			Name: "clusters",
			Kind: "Cluster",
		},
	}
}

func init() {
	rootCmd.AddCommand(apiResourcesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiResourcesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiResourcesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
