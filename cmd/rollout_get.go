package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	"spot-oceancd-cli/pkg/utils"
	"spot-oceancd-cli/viewcontroller"
	"strings"
	"time"
)

type GetOptions struct {
	Watch          bool
	NoColor        bool
	TimeoutSeconds int
}

//  rolloutGetCmd represents the get command
var (
	rolloutGetDescription      = "Obtain information on OceanCD rollouts including details on the verification results performed in each phases"
	rolloutGetShortDescription = "Visual representation of OceanCD rollouts"
	rolloutGetExample          = fmt.Sprintf("  # %s\n  %s %s",
		"Get statuses of your running rollouts", rootCmd.Name(), "rollout get example_rollout")

	rolloutGetWatchExample = fmt.Sprintf("  # %s\n  %s %s\n",
		"Watch statuses of your running rollouts", rootCmd.Name(), "rollout get example_rollout -w")
	rolloutGetOptions = GetOptions{}

	rolloutGetCmd = &cobra.Command{
		Use:     "get SPOTDEPLOYMENT_NAME",
		Short:   rolloutGetShortDescription,
		Long:    rolloutGetDescription,
		Example: strings.Join([]string{rolloutGetExample, rolloutGetWatchExample}, "\n\n"),
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			runRolloutGetAction(args)
		},
	}
)

func init() {
	rolloutCmd.AddCommand(rolloutGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rolloutGetCmd.Flags().BoolVarP(&rolloutGetOptions.Watch, "watch", "w", false, "Watch live updates to the rollout")
	rolloutGetCmd.Flags().BoolVar(&rolloutGetOptions.NoColor, "no-color", false, "Do not colorize output")
	rolloutGetCmd.Flags().IntVarP(&rolloutGetOptions.TimeoutSeconds, "timeout-seconds", "t", 0, "Timeout after specified seconds")
}

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/cmd/get/get_rollout.go#L42
func runRolloutGetAction(args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utils.SetupSignalHandler(cancel)

	rolloutId := args[0]
	controller := viewcontroller.NewRolloutViewController(rolloutId, rolloutGetOptions.NoColor)

	detailedRollout, err := controller.GetRollout()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if !rolloutGetOptions.Watch {
		controller.PrintRollout(&detailedRollout)
	} else {
		rolloutUpdates := make(chan *rollout.DetailedRollout)
		controller.RegisterCallback(func(detailedRollout *rollout.DetailedRollout) {
			rolloutUpdates <- detailedRollout
		})
		if rolloutGetOptions.TimeoutSeconds > 0 {
			ts := time.Duration(rolloutGetOptions.TimeoutSeconds)
			ctx, cancel = context.WithTimeout(ctx, ts*time.Second)
			defer cancel()
		}
		stopCh := ctx.Done()
		go controller.WatchRollout(stopCh, rolloutUpdates)
		controller.Run(ctx)
		close(rolloutUpdates)
	}
}
