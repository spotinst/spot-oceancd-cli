package main

//kube = cmdutil.NewFactory(cmdutil.NewMatchVersionFlags(cf))

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/verchol/applier/pkg/cmd"
	"github.com/verchol/applier/pkg/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

const OceanCDNamespace = "oceancd"
const OceanCDDeployment = "spot-oceancd-controller"

func NewWhoAmICommand() *cobra.Command {
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "update oceancd resources (microservices,  environmetns , etc",
		RunE: func(cmd *cobra.Command, args []string) error {
			RunWhoamiCommand(context.Background(), clientset)
			return nil

		},
	}

	cmd.PersistentFlags().StringP("context", "c", "", "kubernets context to use")
	return cmd
}

func ListClusters() ([]interface{}, error) {
	return cmd.ListEntities(context.Background(), model.ClusterEntity)

}

func RunWhoamiCommand(ctx context.Context, clientSet *kubernetes.Clientset) error {

	d, err := clientSet.AppsV1().Deployments(OceanCDNamespace).Get(ctx, OceanCDDeployment, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to find oceancd deployed component %v  in namespace %v :", OceanCDDeployment, OceanCDNamespace)
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		return fmt.Errorf("oceancd deployed component %v  in namespace %v is not ready", OceanCDDeployment, OceanCDNamespace)
	}

	color.Green("controller deployment is %v", d.Status.AvailableReplicas)
	return nil
}
