/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model/operator"
	"spot-oceancd-cli/pkg/utils"
	"spot-oceancd-operator-commons/handlers/cluster"
	"spot-oceancd-operator-commons/models"
	"strings"
)

type OperatorDeleteOptions struct {
	Namespace         string
	ArgoNamespace     string
	KeepOceanCdCrds   bool
	KeepArgo          bool
	KeepArgoCrds      bool
	KeepArgoNamespace bool
}

// operatorDeleteCmd represents the operator upgrade command
var (
	operatorDeleteDescription = `Delete Ocean CD operator from cluster.`
	operatorDeleteUse         = "delete"
	operatorDeleteExample     = fmt.Sprintf("  # %s\n  %s %s %s %s",
		operatorDeleteDescription, rootCmd.Name(), operatorUse, operatorDeleteUse, "--config /path/to/config")
	operatorDeleteOptions = OperatorDeleteOptions{}

	operatorDeleteCmd = &cobra.Command{
		Use:     operatorDeleteUse,
		Short:   operatorDeleteDescription,
		Long:    operatorDeleteDescription,
		Example: operatorDeleteExample,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
			validateClusterId(context.Background())
			validateClusterIdExists(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorDeleteFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Deleting OceanCD operator manager in cluster %s\n", clusterId)

			if err := runOperatorDeleteCmd(context.Background(), cmd); err != nil {
				fmt.Printf("Failed to delete OceanCD operator manager\n%s\n", err)
			}

			fmt.Printf("OceanCD operator manager was deleted succesfully.\n")
		},
	}
)

func init() {
	operatorCmd.AddCommand(operatorDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	operatorDeleteCmd.Flags().StringVarP(&operatorManagerConfig, "config", "c", "",
		"The configuration applied to OceanCD resources and their dependencies.")
	operatorDeleteCmd.Flags().StringVar(&operatorDeleteOptions.Namespace, "namespace", "", "OceanCD namespace")
	operatorDeleteCmd.Flags().StringVar(&operatorDeleteOptions.ArgoNamespace, "argo-namespace", "", "Argo Rollouts namespace")
	operatorDeleteCmd.Flags().BoolVar(&operatorDeleteOptions.KeepOceanCdCrds, "keep-oceancd-crds", true, "Should we keep OceanCD crds")
	operatorDeleteCmd.Flags().BoolVar(&operatorDeleteOptions.KeepArgoCrds, "keep-argo-crds", true, "Should we keep Argo Rollouts crds")
	operatorDeleteCmd.Flags().BoolVar(&operatorDeleteOptions.KeepArgo, "keep-argo", false, "Should we keep Argo Rollouts")
	operatorDeleteCmd.Flags().BoolVar(&operatorDeleteOptions.KeepArgoNamespace, "keep-argo-namespace", false, "Should we keep Argo Rollouts namespace")
}

func validateOperatorDeleteFlags(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("config").Changed == false {

		if operatorDeleteOptions.Namespace == "" {
			operatorDeleteOptions.Namespace = OceanCDNamespace
		}

		if operatorDeleteOptions.ArgoNamespace == "" {
			operatorDeleteOptions.ArgoNamespace = ArgoRolloutsNamespace
		}
	} else {

		if err := validateOperatorInstallFlags(cmd); err != nil {
			return fmt.Errorf("error: Please provide OceanCD and Argo namespace names, either by using the config file or the command flags\n%w", err)
		}
	}

	return nil
}

func runOperatorDeleteCmd(ctx context.Context, cmd *cobra.Command) error {
	pathToConfig, _ := cmd.Flags().GetString("config")
	if pathToConfig != "" {
		deleteOptions := utils.Options{
			SingleResource: true,
			PathToConfig:   pathToConfig,
		}

		configHandler, err := utils.NewConfigHandler(deleteOptions)
		if err != nil {
			return fmt.Errorf("error: Failed to load config file - %w", err)
		}

		err = configHandler.Handle(ctx, deleteOperatorByConfig)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Using namespaces: '%s' for OceanCD operator and '%s' for Argo Rollouts\n",
			operatorDeleteOptions.Namespace, operatorDeleteOptions.ArgoNamespace)

		err := deleteOperator(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteOperatorByConfig(ctx context.Context, data map[string]interface{}) error {
	config, err := operator.NewOMConfig(data)
	if err != nil {
		return err
	}

	operatorDeleteOptions.Namespace = config.OceanCDConfig.Namespace
	operatorDeleteOptions.ArgoNamespace = config.ArgoRolloutsConfig.General.Namespace

	fmt.Printf("Using namespaces: '%s' for OceanCD operator and '%s' for Argo Rollouts\n",
		operatorDeleteOptions.Namespace, operatorDeleteOptions.ArgoNamespace)

	return deleteOperator(ctx)
}

func deleteOperator(ctx context.Context) error {
	manifestsToDelete, err := oceancd.GetClusterManifests(ctx)
	if err != nil {
		return fmt.Errorf("error: Failed to fetch cluster manifests to delete\n%w", err)
	}

	if err = deleteClusterManifest(ctx, string(v1beta1.Secret), operatorDeleteOptions.Namespace, "spot-oceancd-controller-token"); err != nil {
		return fmt.Errorf("error: Failed to delete secret '%s/%s'\n%w", operatorDeleteOptions.Namespace, "spot-oceancd-controller-token", err)
	}

	if operatorDeleteOptions.KeepArgo == false {

		if err = deleteClusterManifests(ctx, operatorDeleteOptions.ArgoNamespace, manifestsToDelete.Argo, operatorDeleteOptions.KeepArgoCrds, operatorDeleteOptions.KeepArgoNamespace); err != nil {
			return fmt.Errorf("error: Failed to delete Argo Rollouts manifests\n%w", err)
		}
	}

	if err = deleteClusterManifests(ctx, operatorDeleteOptions.Namespace, manifestsToDelete.Operator, operatorDeleteOptions.KeepOceanCdCrds, false); err != nil {
		return fmt.Errorf("error: Failed to delete OceanCD operator manifests\n%w", err)
	}

	if err = deleteClusterSecrets(ctx, operatorDeleteOptions.Namespace, operatorDeleteOptions.ArgoNamespace, manifestsToDelete.Secrets); err != nil {
		return fmt.Errorf("error: Failed to delete OceanCD secrets\n%w", err)
	}

	if err = deleteClusterManifests(ctx, operatorDeleteOptions.Namespace, manifestsToDelete.OM, false, false); err != nil {
		return fmt.Errorf("error: Failed to delete OceanCD operator manager manifests\n%w", err)
	}

	if _, err = oceancd.DeleteCluster(ctx, clusterId); err != nil {
		return fmt.Errorf("error: Failed to delete OceanCD cluster from saas\n%w", err)
	}

	return nil
}

func deleteClusterSecrets(ctx context.Context, oceanCdNamespace string, argoNamespace string, secretsToDelete []operator.SecretMetadata) error {
	var namespaces *v1.NamespaceList = nil
	var err error

	for _, secretToDelete := range secretsToDelete {
		switch secretToDelete.Namespace {
		case "all":
			if namespaces == nil {
				namespaces, err = getNamespaces()
				if err != nil {
					return fmt.Errorf("error: Failed to fetch cluster namespaces\n%w", err)
				}

				for _, namespace := range namespaces.Items {
					if err = deleteClusterManifest(ctx, string(v1beta1.Secret), namespace.Name, secretToDelete.Name); err != nil {
						return fmt.Errorf("error: Failed to delete secret '%s/%s'\n%w", namespace.Name, secretToDelete.Name, err)
					}
				}
			}
		case "oceancd":
			if err = deleteClusterManifest(ctx, string(v1beta1.Secret), oceanCdNamespace, secretToDelete.Name); err != nil {
				return fmt.Errorf("error: Failed to delete secret '%s/%s'\n%w", oceanCdNamespace, secretToDelete.Name, err)
			}
		case "argo-rollouts":
			if err = deleteClusterManifest(ctx, string(v1beta1.Secret), argoNamespace, secretToDelete.Name); err != nil {
				return fmt.Errorf("error: Failed to delete secret '%s/%s'\n%w", argoNamespace, secretToDelete.Name, err)
			}
		}
	}

	return nil
}

func getNamespaces() (*v1.NamespaceList, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	k8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	namespaceClient := k8sClientSet.CoreV1().Namespaces()
	return namespaceClient.List(context.TODO(), metav1.ListOptions{})
}

func deleteClusterManifests(ctx context.Context, namespace string, manifestsToDelete []operator.ManifestMetadata, keepCrds bool, keepNamespace bool) error {
	for _, manifestToDelete := range manifestsToDelete {

		if keepCrds == false || manifestToDelete.Kind != "CustomResourceDefinition" {

			if err := deleteClusterManifest(ctx, manifestToDelete.Kind, namespace, manifestToDelete.Name); err != nil {
				return fmt.Errorf("error: Failed to delete %s '%s/%s'\n%w", manifestToDelete.Kind, namespace, manifestToDelete.Name, err)
			}
		}
	}

	if keepNamespace == false {

		if err := deleteClusterManifest(ctx, "Namespace", namespace, namespace); err != nil {
			return fmt.Errorf("error: Failed to delete Namespace '%s/%s'\n%w", namespace, namespace, err)
		}
	}

	return nil
}

func deleteClusterManifest(_ context.Context, kind string, namespace string, name string) error {
	deleteHandler := cluster.BaseDeleteHandler{}

	deleteRequest := &models.ClusterHandlerPayload{
		Kind:      kind,
		Namespace: namespace,
		Name:      name,
	}

	if err := deleteHandler.Delete(deleteRequest); err != nil {

		if strings.Contains(err.Error(), "not found") == false && strings.Contains(err.Error(), "no matches for") == false {
			fmt.Printf("Failed to delete %s '%s/%s' configuration - %s \n", kind, namespace, name, err.Error())
		}
	} else {
		fmt.Printf("Successfully deleted %s '%s/%s' configuration\n", kind, namespace, name)
	}

	return nil
}
