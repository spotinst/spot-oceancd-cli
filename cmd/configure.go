package cmd

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"runtime"
)

var (
	configFile      = filepath.Join(userHomeDir(), "spotinst", ".oceancd.ini")
	tokenQuestion = []*survey.Question{
		{
			Name: "Token",
			Prompt:   &survey.Input{Message: "Enter your spot access token"},
			Validate: survey.Required,
		},
	}
	profileQuestion = []*survey.Question{
		{
			Name: "profile",
			Prompt: &survey.Input{
				Message: "Enter profile",
				Default: "default",
			},
		},
	}

	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure config fileToApply params",
		Run: func(cmd *cobra.Command, args []string) {
			runConfigureCmd(context.Background())
		},
	}
)

type ConfigFileFields struct {
	Token   string
	Url     string
	Profile string
}

func runConfigureCmd(ctx context.Context) {
	answers := ConfigFileFields{Url: url}
	if isTokenFromConfig {
		err := survey.Ask(tokenQuestion, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.Token = token
	}

	if isProfileOverriden == false {
		err := survey.Ask(profileQuestion, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.Profile = profile
	}

	dir := filepath.Dir(configFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create fileToApply '%s' - %s\n", configFile, err.Error())
			return
		}
	} else if err != nil {
		fmt.Printf("Failed to verify if dir '%s' exists - %s\n", dir, err.Error())
	}

	// Create or update configuration.
	cfg, loadErr := ini.LooseLoad(configFile)
	if loadErr != nil {
		fmt.Printf("Failed to load fileToApply '%s' - %s\n", configFile, loadErr.Error())
		return
	}

	// Create a new `default` section.
	sec, secErr := cfg.NewSection(answers.Profile)
	if secErr != nil {
		fmt.Printf("Failed to create new section in '%s' - %s\n", answers.Profile, secErr.Error())
		return
	}

	// Create a new `token` key.
	if _, err := sec.NewKey("token", answers.Token); err != nil {
		fmt.Printf("Failed to create key '%s' - %s\n", answers.Token, err.Error())
		return
	}

	// Create a new `url` key.
	if _, err := sec.NewKey("url", answers.Url); err != nil {
		fmt.Printf("Failed to create key '%s' - %s\n", answers.Url, err.Error())
		return
	}

	// Write out configuration to a fileToApply.
	if err := cfg.SaveTo(configFile); err != nil {
		fmt.Printf("Failed to save fileToApply '%s' - %s\n", configFile, err.Error())
		return
	}

}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func userHomeDir() string {
	// Windows
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
