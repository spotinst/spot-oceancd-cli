package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// retryCmd represents the retry command
var (
	retryDescription      = "Available for the last rolled back SpotDeployment only. With this action you will be able to retry your full rollout"
	retryShortDescription = "Retry a rollout"
	retryExample          = getRolloutActionExample(retryShortDescription, oceancd.RetryAction)

	retryCmd = &cobra.Command{
		Use:     oceancd.RetryAction + " ROLLOUT_ID ",
		Short:   retryShortDescription,
		Long:    retryDescription,
		Example: retryExample,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutAction(oceancd.RetryAction, args, "retried")
		},
	}
)

func init() {
	rolloutCmd.AddCommand(retryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
