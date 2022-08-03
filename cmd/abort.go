package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// abortCmd represents the abort command
var (
	abortDescription = "The rollout will be terminated and the previous version (i.e., Stable) will be restored"
	abortCmd         = &cobra.Command{
		Use:   oceancd.AbortAction + " ROLLOUT_ID",
		Short: abortDescription,
		Long:  abortDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutActionCmd(oceancd.AbortAction, args, "rolled back")
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
