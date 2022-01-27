package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	OceanCDNamespace  = "oceancd"
	OceanCDDeployment = "spot-oceancd-controller"
)

// whoAmICmd represents the whoami command
var whoAmICmd = &cobra.Command{
	Use:   "whoami",
	Short: "Check controller status",
	Run: func(cmd *cobra.Command, args []string) {
		runWhoAmICmd(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(whoAmICmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoAmICmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoAmICmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runWhoAmICmd(ctx context.Context) {
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
