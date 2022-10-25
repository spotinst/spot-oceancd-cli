package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd"
	"strconv"
)

// workloadRollbackCmd represents the rollback command
var (
	workloadRollbackDescription = "Rolls back to one of the last 20 revisions of your choice. Such action is applicable " +
		"only to your non-live versions and will trigger a new rollout on your behalf"
	workloadRollbackShortDescription = "Rolls back to one of the last 20 revisions of your choice"
	workloadRollbackExample          = fmt.Sprintf("  # %s\n  %s %s %s %s %s", workloadRollbackShortDescription,
		rootCmd.Name(), workloadUse, oceancd.RollbackAction, spotdeploymentNameExample, revisionIdExample)

	workloadRollbackCmd = &cobra.Command{
		Use:     oceancd.RollbackAction + " SPOTDEPLOYMENT_NAME REVISION_ID",
		Short:   workloadRollbackShortDescription,
		Long:    workloadRollbackDescription,
		Example: workloadRollbackExample,
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
			runWorkloadAction(oceancd.RollbackAction, args, "rolled back")
		},
	}
)

func init() {
	workloadCmd.AddCommand(workloadRollbackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
