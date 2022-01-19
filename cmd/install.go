package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/verchol/applier/pkg/utils"
)

const installScript = ""

func NewInstallCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install oceancd controller",
		//Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInstall(cmd.Context(), args)
		},
	}
	cmd.PersistentFlags().StringP("installDir", "f", "", "installation directory path")
	pflag := cmd.PersistentFlags().Lookup("installDir")
	viper.BindPFlag("installDir", pflag)
	return cmd
}

func RunInstall(ctx context.Context, args []string) error {
	//installDir := viper.GetString("installDir")
	//err := cmd.InstallFromDir(installDir)

	fmt.Println(utils.RenderTerminalString("ocean cd"))

	utils.MessageWithProgress("create  namespace", 2)
	utils.MessageWithProgress("create  service account", 2)
	utils.MessageWithProgress("deploy  oceancd controller", 2)
	return nil
}
