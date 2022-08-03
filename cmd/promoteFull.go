package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// promoteFullCmd represents the promoteFull command
var (
	promoteFullDescription = "Promote a phase to the end of the rollout, triggering a success"
	promoteFullCmd         = &cobra.Command{
		Use:   oceancd.PromoteFullAction + " ROLLOUT_ID",
		Short: promoteFullDescription,
		Long:  promoteFullDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutActionCmd(oceancd.PromoteFullAction, args, "fully promoted")
		},
	}
)

func init() {
	rolloutCmd.AddCommand(promoteFullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
