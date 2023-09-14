package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootDescription = `Ocean CD controls oceancd resources
For more information visit our github repo https://github.com/spotinst/spot-oceancd-cli`
	profile               string
	token                 string
	url                   string
	clusterUrl            string
	clusterId             string
	namespace             string
	isTokenFromConfig     = false
	isProfileOverriden    = false
	isClusterIdOverridden = false
	isNamespaceOverridden = false

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
	rootCmd.SetUsageTemplate(usageTemplate())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "name of credentials profile to use")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "unqiue spot token for api authentication")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "Base ocean cd api url")
	rootCmd.PersistentFlags().StringVar(&clusterUrl, "clusterUrl", "", "Base ocean cd cluster api url")
	_ = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	_ = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	_ = viper.BindPFlag("clusterUrl", rootCmd.PersistentFlags().Lookup("clusterUrl"))
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

	clusterUrl = viper.GetString("clusterUrl")
	if clusterUrl == "" {
		clusterUrlKey := fmt.Sprintf("%s.%s", profile, "clusterUrl")
		clusterUrl = viper.GetString(clusterUrlKey)

		if clusterUrl == "" {
			clusterUrl = "https://cluster-gateway.oceancd.io"
		}

		viper.Set("clusterUrl", clusterUrl)
	}

	clusterId = viper.GetString("clusterId")
	if clusterId == "" {
		isClusterIdOverridden = true
		clusterIdKey := fmt.Sprintf("%s.%s", profile, "clusterId")
		clusterId = viper.GetString(clusterIdKey)
		viper.Set("clusterId", clusterId)
	}

	namespace = viper.GetString("namespace")
	if namespace == "" {
		isNamespaceOverridden = true
		namespaceKey := fmt.Sprintf("%s.%s", profile, "namespace")
		namespace = viper.GetString(namespaceKey)
		viper.Set("namespace", namespace)
	}

	return
}

func validateToken(_ context.Context) {
	if token == "" {
		fmt.Println("You haven't specify your access token. You can use \"oceancd configure\" to create a config file")
		os.Exit(1)
	}
}

func validateClusterId(_ context.Context) {
	if clusterId == "" {
		fmt.Printf(`You haven't specified your cluster ID. You can use "oceancd configure" to configure the 
missing parameters using the profile variables or use the appropriate flag: --%s`, ClusterIdFlagLabel)
		fmt.Println("")
		os.Exit(1)
	}
}

func validateNamespace(_ context.Context) {
	if namespace == "" {
		fmt.Printf(`You haven't specified your namespace. You can use "oceancd configure" to configure the 
missing parameters using the profile variables or use the appropriate flag: --%s`, NamespaceFlagLabel)
		fmt.Println("")
		os.Exit(1)
	}
}

func validateClusterIdExists(_ context.Context) {
	resource, err := oceancd.GetEntity(context.Background(), model.ClusterEntity, clusterId)
	if err != nil {
		if err.Error() != "resource does not exist" {
			fmt.Printf("Failed to fetch cluster %s from saas, %s\n", clusterId, err.Error())
			os.Exit(1)
		}
	}

	if resource == nil {
		fmt.Printf("Cluster %s does not exists\n", clusterId)
		os.Exit(1)
	}
}

func validateClusterIdNotExists(_ context.Context) {
	resource, err := oceancd.GetEntity(context.Background(), model.ClusterEntity, clusterId)
	if err != nil {
		if err.Error() != fmt.Sprintf("error: Resource '%s/%s' does not exist", model.ClusterEntity, clusterId) {
			fmt.Printf("Failed to fetch cluster %s from saas, %s\n", clusterId, err.Error())
			os.Exit(1)
		}
	}

	if resource != nil {
		fmt.Printf("Cluster %s allready exists\n", clusterId)
		os.Exit(1)
	}
}

func usageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{if (eq .Name "oceancd")}}

Rollout Commands:{{$cmds := .Commands}}{{range $cmds}}{{if (and .IsAvailableCommand (eq .Name "rollout"))}}{{if .HasAvailableSubCommands}}{{$rolloutCmds := .Commands}}{{range $rolloutCmds}}
  rollout {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
