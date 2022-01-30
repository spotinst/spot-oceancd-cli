/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/spf13/cobra"
)

const (
	OceanCDNamespace  = "oceancd"
	OceanCDDeployment = "spot-oceancd-controller"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check controller status",
	Run: func(cmd *cobra.Command, args []string) {
		runControllerStatusCmd(context.Background())
	},
}

func init() {
	controllerCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runControllerStatusCmd(ctx context.Context) {
	config, configErr := ctrl.GetConfig()
	if configErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", configErr.Error())
	}

	clientSet, clientError := kubernetes.NewForConfig(config)
	if clientError != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", configErr.Error())
	}

	d, err := clientSet.AppsV1().Deployments(OceanCDNamespace).Get(ctx, OceanCDDeployment, v1.GetOptions{})
	if err != nil {
		fmt.Printf("oceancd controller not found - %s\n", err.Error())
		return
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		fmt.Println("oceancd controller is not ready")
		return
	}

	fmt.Println("oceancd controller is running ")

	return
}
