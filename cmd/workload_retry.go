package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
	"strconv"
)

// workloadRetryCmd represents the retry command
var (
	workloadRetryDescription = "Retries your latest rolled-back deployment. " +
		"This action is restricted to one revision only and will trigger a new rollout on your behalf"
	workloadRetryShortDescription = "Retries your latest rolled-back deployment"
	workloadRetryExample          = fmt.Sprintf("  # %s\n  %s %s %s %s %s", workloadRetryShortDescription,
		rootCmd.Name(), workloadUse, oceancd.RetryAction, spotdeploymentNameExample, revisionIdExample)

	workloadRetryCmd = &cobra.Command{
		Use:     oceancd.RetryAction + " SPOTDEPLOYMENT_NAME REVISION_ID",
		Short:   workloadRetryShortDescription,
		Long:    workloadRetryDescription,
		Example: workloadRetryExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(2)(cmd, args); err != nil {
				fmt.Println("You must specify SpotDeployment name and revision ID.")
				return errors.New(fmt.Sprintf("Wrong number of arguments: %s\n", err.Error()))
			}

			if _, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("Revision ID must be a digit.")
				return errors.New(fmt.Sprintf("Wrong type of revision ID: %s\n", args[1]))
			}

			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			runWorkloadAction(oceancd.RetryAction, args, "retried")
		},
	}
)

func init() {
	workloadCmd.AddCommand(workloadRetryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
