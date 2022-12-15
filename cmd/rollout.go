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
	rolloutUse         = "rollout"
	rolloutDescription = "This command consists of multiple subcommands which can perform changes on a SpotDeployment rollout"

	rolloutIdExample             = "rol-a78dsds9s"
	rolloutActionExampleTemplate = "  # %s\n  %s %s %s %s"

	rolloutCmd = &cobra.Command{
		Use:     rolloutUse,
		Short:   rolloutDescription,
		Long:    rolloutDescription,
		Example: strings.Join([]string{rolloutGetExample, abortExample, pauseExample, promoteExample, promoteFullExample, retryExample}, "\n\n"),
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
	rolloutId := args[0]
	actionRequest := map[string]string{"action": action}

	err := oceancd.SendRolloutAction(rolloutId, actionRequest)
	if err != nil {
		fmt.Printf("Failed to %s the rollout %s: %s\n", action, rolloutId, err.Error())
	} else {
		fmt.Printf("Successfully %s resource %s\n", actionPastForm, rolloutId)
	}
}

func validateRolloutActionArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("You must specify a rollout id.")
		return errors.New("error: Rollout id not specified")
	} else if len(args) > 1 {
		fmt.Println("You can only specify one rollout id.")
		return errors.New(fmt.Sprintf("error: Too many arguments: %+v", args))
	} else {
		rolloutId := args[0]
		if strings.HasPrefix(rolloutId, "rol-") == false {
			fmt.Printf(`%s is not a valid rollout id`, rolloutId)
			return errors.New(fmt.Sprintf("error: Invalid rollout id: %s", rolloutId))
		}

		if false == regexp.MustCompile(`^[a-zA-Z\d]*$`).MatchString(strings.TrimPrefix(rolloutId, "rol-")) {
			fmt.Printf(`%s is not a valid rollout id`, rolloutId)
			return errors.New(fmt.Sprintf("error: Invalid rollout id: %s", rolloutId))
		}
	}

	return nil
}

func getRolloutActionExample(description string, action string) string {
	return fmt.Sprintf(rolloutActionExampleTemplate, description, rootCmd.Name(), rolloutUse, action, rolloutIdExample)
}
