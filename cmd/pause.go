package cmd

import (
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// pauseCmd represents the pause command
var (
	pauseDescription = "Pause a full rollout. Once the rollout is resumed, it will restart the phase where it left off"
	pauseCmd         = &cobra.Command{
		Use:   oceancd.PauseAction + " ROLLOUT_ID",
		Short: pauseDescription,
		Long:  pauseDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return validateRolloutActionArgs(cmd, args)
		},
		Run: func(_ *cobra.Command, args []string) {
			runRolloutAction(oceancd.PauseAction, args, "paused")
		},
	}
)

func init() {
	rolloutCmd.AddCommand(pauseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
