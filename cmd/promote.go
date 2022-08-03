package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// promoteCmd represents the promote command
var (
	promoteDescription = "Promote one phase to the next"
	promoteCmd         = &cobra.Command{
		Use:   oceancd.PromoteAction + " ROLLOUT_ID",
		Short: promoteDescription,
		Long:  promoteDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutActionCmd(oceancd.PromoteAction, args, "promoted")
		},
	}
)

func init() {
	rolloutCmd.AddCommand(promoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
