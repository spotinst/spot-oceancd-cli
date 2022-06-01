/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/spf13/cobra"
)

const (
	OceanCDNamespace  = "oceancd"
	OceanCDOperator   = "spot-oceancd-operator"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check operator status",
	Run: func(cmd *cobra.Command, args []string) {
		runOperatorStatusCmd(context.Background())
	},
}

func init() {
	operatorCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runOperatorStatusCmd(ctx context.Context) {
	config, configErr := ctrl.GetConfig()
	if configErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", configErr.Error())
		return
	}

	clientset, clientError := kubernetes.NewForConfig(config)
	if clientError != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", configErr.Error())
		return
	}

	discoveryClient := cacheddiscovery.NewMemCacheClient(clientset.Discovery())
	discoveryRESTMapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	gvk, gvkErr := discoveryRESTMapper.KindFor(schema.GroupVersionResource{Resource: "ClusterServiceVersion"})

	if gvkErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", gvkErr.Error())
		return
	}

	restMapping, mappingErr := discoveryRESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)

	if mappingErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", mappingErr.Error())
		return
	}

	dynamicClient, dynamicErr := dynamic.NewForConfig(config)

	if dynamicErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", dynamicErr.Error())
		return
	}

	items, listErr := dynamicClient.
		Resource(restMapping.Resource).
		Namespace(OceanCDNamespace).
		List(context.TODO(), metav1.ListOptions{})

	if listErr != nil {
		fmt.Printf("The connection to the kubernetes server was refused - %s\n", listErr.Error())
		return
	}

	for _, operator := range items.Items {
		if operator.Object["spec"].(map[string]interface{})["displayName"] == OceanCDOperator {
			version := operator.Object["spec"].(map[string]interface{})["version"]

			if operator.Object["status"].(map[string]interface{})["phase"] == "Succeeded" {
				fmt.Printf("%s %s is running\n", OceanCDOperator, version)
			} else {
				fmt.Printf("%s %s is not ready\n", OceanCDOperator, version)
			}

			return
		}
	}

	fmt.Printf("%s is not installed\n", OceanCDOperator)
}
