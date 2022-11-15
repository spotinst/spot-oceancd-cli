package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
)

// workloadRestartCmd represents the restart command
var (
	workloadRestartDescription      = "Restarts your currently running pods. Such action is applicable only for your LIVE revision"
	workloadRestartShortDescription = "Restarts your currently running pods"
	workloadRestartExample          = fmt.Sprintf("  # %s\n  %s %s %s %s",
		workloadRestartShortDescription, rootCmd.Name(), workloadUse, oceancd.RestartAction, spotdeploymentNameExample)

	workloadRestartCmd = &cobra.Command{
		Use:     oceancd.RestartAction + " SPOTDEPLOYMENT_NAME",
		Short:   workloadRestartShortDescription,
		Long:    workloadRestartDescription,
		Example: workloadRestartExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				fmt.Println("You must specify SpotDeployment name.")
				return errors.New(fmt.Sprintf("Wrong number of arguments: %s\n", err.Error()))
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			runWorkloadAction(oceancd.RestartAction, args, "restarted")
		},
	}
)

func init() {
	workloadCmd.AddCommand(workloadRestartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
