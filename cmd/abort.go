package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// abortCmd represents the abort command
var (
	abortDescription = "This command stops progressing the current SpotDeployment rollout and reverts all steps. " +
		"The previous ReplicaSet will be active. Note the 'spec.template' still represents the new SpotDeployment " +
		"rollout version. Updating the 'spec.template' back to the previous version will fully revert the SpotDeployment rollout"
	abortCmd = &cobra.Command{
		Use:   oceancd.AbortAction + " ROLLOUT_ID",
		Short: "Abort a rollout",
		Long:  abortDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutAction(oceancd.AbortAction, args, "rolled back")
		},
	}
)

func init() {
	rolloutCmd.AddCommand(abortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
