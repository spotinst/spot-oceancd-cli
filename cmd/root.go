package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootDescription = `Ocean CD controls oceancd resources
For more information visit our github repo https://github.com/spotinst/spot-oceancd-cli`
	profile            string
	token              string
	url                string
	isTokenFromConfig  = false
	isProfileOverriden = false

	rootCmd = &cobra.Command{
		Use:   "oceancd",
		Short: "Ocean CD controls oceancd resources",
		Long:  rootDescription,
	}

	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "name of credentials profile to use")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "unqiue spot token for api authentication")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "Base ocean cd api url")
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	_ = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
}

func initConfig() {
	// Find home directory.
	home, _ := os.UserHomeDir()
	cfgFile := filepath.Join(home, "spotinst", ".oceancd.ini")
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	profile = viper.GetString("profile")
	if profile == "" {
		profile = "default"
		viper.Set("profile", profile)
	} else {
		isProfileOverriden = true
	}

	token = viper.GetString("token")
	if token == "" {
		isTokenFromConfig = true
		tokenKey := fmt.Sprintf("%s.%s", profile, "token")
		token = viper.GetString(tokenKey)
		viper.Set("token", token)
	}

	url = viper.GetString("url")
	if url == "" {
		urlKey := fmt.Sprintf("%s.%s", profile, "url")
		url = viper.GetString(urlKey)

		if url == "" {
			url = "https://api.spotinst.io"
		}

		viper.Set("url", url)
	}

	return
}

func validateToken(_ context.Context) {
	if token == "" {
		fmt.Println("You haven't specify your access token. You can use \"oceancd configure\" to create a config file")
		os.Exit(1)
	}
}
