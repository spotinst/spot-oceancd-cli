package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"spot-oceancd-cli/pkg/oceancd"
	"strings"
)

// workloadCmd represents the retry command
var (
	workloadUse               = "workload"
	workloadDescription       = "Performs changes on a workload level"
	revisionIdExample         = "19"
	spotdeploymentNameExample = "spotdeployment-example"

	workloadCmd = &cobra.Command{
		Use:     workloadUse,
		Short:   workloadDescription,
		Long:    workloadDescription,
		Example: strings.Join([]string{workloadRestartExample, workloadRetryExample, workloadRollbackExample}, "\n\n"),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
			validateClusterId(context.Background())
			validateNamespace(context.Background())
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(workloadCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	workloadCmd.PersistentFlags().StringVar(&clusterId, ClusterIdFlagLabel, "", ClusterIdFlagDescription)
	workloadCmd.PersistentFlags().StringVar(&namespace, NamespaceFlagLabel, "", NamespaceFlagDescription)
	_ = viper.BindPFlag("clusterId", workloadCmd.PersistentFlags().Lookup(ClusterIdFlagLabel))
	_ = viper.BindPFlag("namespace", workloadCmd.PersistentFlags().Lookup(NamespaceFlagLabel))
}

func runWorkloadAction(action string, args []string, actionPastForm string) {
	spotDeploymentName := args[0]

	pathParam := map[string]string{
		"action":             action,
		"spotDeploymentName": spotDeploymentName,
		"namespace":          viper.GetString("namespace"),
	}

	if len(args) == 2 {
		pathParam["revisionId"] = args[1]
	}

	queryParam := map[string]string{
		"clusterId": viper.GetString("clusterId"),
		"kind":      "SpotDeployment",
	}

	err := oceancd.SendWorkloadAction(pathParam, queryParam)
	if err != nil {
		fmt.Printf("Failed to %s the workload %s: %s\n", action, spotDeploymentName, err.Error())
	} else {
		fmt.Printf("Successfully %s workload %s\n", actionPastForm, spotDeploymentName)
	}
}
