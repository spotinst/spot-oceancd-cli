package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"spot-oceancd-cli/pkg/oceancd"
	"strings"
)

// retryCmd represents the retry command
var (
	rolloutDescription = "Perform changes on a rollout level"
	rolloutCmd         = &cobra.Command{
		Use:   "rollout",
		Short: rolloutDescription,
		Long:  rolloutDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
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
	rootCmd.AddCommand(rolloutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runRolloutAction(action string, args []string, actionPastForm string) {
	rolloutID := args[0]
	actionRequest := map[string]string{"action": action}

	err := oceancd.SendRolloutAction(rolloutID, actionRequest)
	if err != nil {
		fmt.Printf("Failed to %s the rollout - %s\n", action, err.Error())
	} else {
		fmt.Printf("Successfully %s resource %s\n", actionPastForm, rolloutID)
	}
}

func validateRolloutActionArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("You must specify a rollout ID.")
		return errors.New("error: Rollout ID not specified")
	} else if len(args) > 1 {
		fmt.Println("You can only specify one rollout ID.")
		return errors.New(fmt.Sprintf("error: Too many arguments: %+v", args))
	} else {
		rolloutID := args[0]
		if strings.HasPrefix(rolloutID, "rol-") == false {
			fmt.Printf(`%s is not a valid rollout id`, rolloutID)
			return errors.New(fmt.Sprintf("error: Invalid Rollout ID: %s", rolloutID))
		}

		if false == regexp.MustCompile(`^[a-zA-Z\d]*$`).MatchString(strings.TrimPrefix(rolloutID, "rol-")) {
			fmt.Printf(`%s is not a valid rollout id`, rolloutID)
			return errors.New(fmt.Sprintf("error: Invalid Rollout ID: %s", rolloutID))
		}
	}

	return nil
}
