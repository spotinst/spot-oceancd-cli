package cmd

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"runtime"
)

var (
	configureDescription = `Modify oceancd config file
You can create different profiles for different tokens.`
	configureExamples = fmt.Sprintf(`  # Create a profile using interactive mode
  oceancd configure

  # Create a profile programmatically using flags
  oceancd configure --profile=PROFILE --token=TOKEN --%s=CLUSTER_ID --%s=NAMESPACE

  # Create a profile with custom api url
  oceancd configure --url=URL`, ClusterIdFlagLabel, NamespaceFlagLabel)
	configFile    = filepath.Join(userHomeDir(), "spotinst", ".oceancd.ini")
	tokenQuestion = []*survey.Question{
		{
			Name:     "Token",
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
	clusterIdQuestion = []*survey.Question{
		{
			Name: "clusterId",
			Prompt: &survey.Input{
				Message: "Enter cluster id",
			},
		},
	}
	namespaceQuestion = []*survey.Question{
		{
			Name: "namespace",
			Prompt: &survey.Input{
				Message: "Enter namespace",
			},
		},
	}

	configureCmd = &cobra.Command{
		Use:     "configure",
		Short:   "Modify oceancd config file",
		Long:    configureDescription,
		Example: configureExamples,
		Run: func(cmd *cobra.Command, args []string) {
			runConfigureCmd(context.Background())
		},
	}
)

type ConfigFileFields struct {
	Token     string
	Url       string
	Profile   string
	ClusterId string
	Namespace string
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

	if isClusterIdOverridden {
		err := survey.Ask(clusterIdQuestion, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.ClusterId = clusterId
	}

	if isNamespaceOverridden {
		err := survey.Ask(namespaceQuestion, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.Namespace = namespace
	}

	dir := filepath.Dir(configFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create file '%s' - %s\n", configFile, err.Error())
			return
		}
	} else if err != nil {
		fmt.Printf("Failed to verify if dir '%s' exists - %s\n", dir, err.Error())
	}

	// Create or update configuration.
	cfg, loadErr := ini.LooseLoad(configFile)
	if loadErr != nil {
		fmt.Printf("Failed to load file '%s' - %s\n", configFile, loadErr.Error())
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

	// Create a new `clusterId` key.
	if _, err := sec.NewKey("clusterId", answers.ClusterId); err != nil {
		fmt.Printf("Failed to create key '%s' - %s\n", answers.ClusterId, err.Error())
		return
	}

	// Create a new `namespace` key.
	if _, err := sec.NewKey("namespace", answers.Namespace); err != nil {
		fmt.Printf("Failed to create key '%s' - %s\n", answers.Namespace, err.Error())
		return
	}

	// Write out configuration to a file.
	if err := cfg.SaveTo(configFile); err != nil {
		fmt.Printf("Failed to save file '%s' - %s\n", configFile, err.Error())
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
	configureCmd.PersistentFlags().StringVar(&clusterId, ClusterIdFlagLabel, "", ClusterIdFlagDescription)
	configureCmd.PersistentFlags().StringVar(&namespace, NamespaceFlagLabel, "", NamespaceFlagDescription)
	_ = viper.BindPFlag("clusterId", workloadCmd.PersistentFlags().Lookup(ClusterIdFlagLabel))
	_ = viper.BindPFlag("namespace", workloadCmd.PersistentFlags().Lookup(NamespaceFlagLabel))
}

func userHomeDir() string {
	// Windows
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
