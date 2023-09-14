/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model/operator"
	"spot-oceancd-cli/pkg/utils"
	"spot-oceancd-operator-commons/component_configs"
	"spot-oceancd-operator-commons/handlers/cluster"
	"spot-oceancd-operator-commons/helpers"
	"strings"
)

var operatorManagerConfig string
var shouldCreateNamespace bool

// operatorInstallCmd represents the operator install command
var (
	isOperatorInstallCommand        = true
	operatorInstallDescription      = `Installs Ocean CD operator on current cluster with dependencies based on provided config.`
	operatorInstallShortDescription = "Installs Ocean CD operator on current cluster"
	operatorInstallUse              = "install"
	operatorInstallExample          = fmt.Sprintf("  # %s\n  %s %s %s %s",
		operatorInstallShortDescription, rootCmd.Name(), operatorUse, operatorInstallUse, "--config /path/to/config")

	operatorInstallCmd = &cobra.Command{
		Use:     operatorInstallUse,
		Short:   operatorInstallShortDescription,
		Long:    operatorInstallDescription,
		Example: operatorInstallExample,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			validateToken(context.Background())
			validateClusterId(context.Background())
			validateClusterIdNotExists(context.Background())
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorInstallFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Installing OceanCD operator manager in cluster %s\n", clusterId)

			if err := runOperatorInstallCmd(context.Background(), cmd); err != nil {
				fmt.Printf("Failed to install OceanCD operator manager\n%s\n", err)
			}

			fmt.Printf("OceanCD operator manager installation finished succesfully.\n")
		},
	}
)

func init() {
	operatorCmd.AddCommand(operatorInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	operatorInstallCmd.Flags().StringVarP(&operatorManagerConfig, "config", "c", "",
		"The configuration applied to OceanCD resources and their dependencies.")
	operatorInstallCmd.Flags().BoolVar(&shouldCreateNamespace, "create-namespace", true, "Should it create OceanCD namespace. Default true")
}

func runOperatorInstallCmd(ctx context.Context, cmd *cobra.Command) error {
	var err error

	pathToConfig := ""
	if cmd.Flags().Lookup("config").Changed {
		pathToConfig, err = cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("error: Failed to parse --config flag - %w", err)
		}
	}

	installOptions := utils.Options{
		SingleResource: true,
		PathToConfig:   pathToConfig,
	}

	configHandler, err := utils.NewConfigHandler(installOptions)
	if err != nil {
		return fmt.Errorf("error: Failed to load config file - %w", err)
	}

	err = configHandler.Handle(ctx, installOperator)
	if err != nil {
		return err
	}

	return nil
}

func validateOperatorInstallFlags(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("config").Changed == false {
		return nil
	}

	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("error: Failed to parse --config flag - %w", err)
	}

	if pathToConfig == "" {
		return fmt.Errorf("error: Path to config file must be specified")
	}

	fileExtensionWithDot := filepath.Ext(pathToConfig)
	if err = utils.IsFileTypeSupported(fileExtensionWithDot); err != nil {
		return err
	}

	return nil
}

func installOperator(ctx context.Context, data map[string]interface{}) error {
	config, err := operator.NewOMConfig(data)
	if err != nil {
		return err
	}

	payload := operator.NewOMManifestsRequest(config)
	output, err := oceancd.GetOMInstallationManifests(ctx, payload)
	if err != nil {
		return fmt.Errorf("error: Failed to fetch installation resources\n%w", err)
	}

	resources, err := helpers.ConvertToUnstructuredSlice(output.OM.Apply)
	if err != nil {
		return fmt.Errorf("error: Failed to convert manifests to unstructured\n%w", err)
	}

	operatorManagerConfigMap, err := buildOperatorManagerConfigMap(config)
	if err != nil {
		return fmt.Errorf("error: Failed to build operator manager ConfigMap\n%w", err)
	}

	configMapResource, err := convertOperatorManagerConfigMap(operatorManagerConfigMap)
	if err != nil {
		return fmt.Errorf("error: Failed to convert operator manager ConfigMap\n%w", err)
	}

	resources = append(resources, configMapResource)

	if err = createOceancdNamespace(config.OceanCDConfig.Namespace); err != nil {
		return fmt.Errorf("error: Failed to create OceanCD Namespace %s\n%w", config.OceanCDConfig.Namespace, err)
	}

	applyHandler := cluster.BaseApplyHandler{}

	if isOperatorInstallCommand {
		clusterToken, err := oceancd.CreateClusterToken(ctx)
		if err != nil {
			return fmt.Errorf("error: Failed to create cluster token\n%w", err)
		}

		operatorManagerSecret := buildOperatorManagerSecret(clusterToken, config.OceanCDConfig.Namespace)
		secretResource, err := convertOperatorManagerSecret(operatorManagerSecret)
		if err != nil {
			return fmt.Errorf("error: Failed to convert operator manager Secret\n%w", err)
		}

		if err = applyHandler.Apply(secretResource); err != nil {
			return fmt.Errorf("error: failed to apply secret resource\n%w", err)
		}

		fmt.Printf("Successfuly created Secret '%s/%s'\n", operatorManagerSecret.GetNamespace(), operatorManagerSecret.GetName())
	}

	kindByPriority := helpers.ReverseMap(output.OM.Priority)
	manifestsToApply := helpers.ConvertUnstructuredListToMapByKind(resources)

	for priority := 1; priority <= len(kindByPriority); priority++ {
		kindToHandle := kindByPriority[priority]

		if kindManifestsToApply, exists := manifestsToApply[kindToHandle]; exists {

			for _, resource := range kindManifestsToApply {

				if err = applyHandler.Apply(resource); err != nil {
					return fmt.Errorf("error: Failed to apply operator manifests resources\n%w", err)
				}

				fmt.Printf("Successfuly created %s '%s/%s'\n", kindToHandle, resource.GetNamespace(), resource.GetName())
			}
		}
	}

	return nil
}

func createOceancdNamespace(oceancdNamespace string) error {
	k8sClient, err := ctrlClient.New(ctrl.GetConfigOrDie(), ctrlClient.Options{})
	if err != nil {
		return fmt.Errorf("error: Failed to create k8s client\n%w", err)
	}

	if ok, err := isNamespaceExists(k8sClient, oceancdNamespace); err != nil {
		return fmt.Errorf("error: Failed to check if OceanCD Namespace %s exists\n%w", oceancdNamespace, err)
	} else if ok {
		return nil
	}

	if shouldCreateNamespace {

		if err = createNamespace(k8sClient, oceancdNamespace); err != nil {
			return fmt.Errorf("error: Failed to create OceanCD Namespace %s\n%w", oceancdNamespace, err)
		}
	} else {
		return fmt.Errorf("error: OceanCD namespace '%s' does not exist. Please create it or set the flag 'create-namespace' to true", oceancdNamespace)
	}

	fmt.Printf("Successfully created OceanCD Namespace '%s'\n", oceancdNamespace)

	return nil
}

func createNamespace(k8sClient k8sclient.Client, oceancdNamespace string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   oceancdNamespace,
			Labels: map[string]string{"app": "spot-oceancd-operator-manager"},
		},
	}

	if err := k8sClient.Create(context.TODO(), namespace); err != nil {
		return err
	}

	return nil
}

func isNamespaceExists(k8sClient k8sclient.Client, oceancdNamespace string) (bool, error) {
	namespace := &corev1.Namespace{}

	err := k8sClient.Get(context.TODO(), k8sclient.ObjectKey{
		Name: oceancdNamespace,
	}, namespace)

	if err != nil {

		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func buildOperatorManagerSecret(clusterToken *operator.ClusterTokenResponse, namespace string) *corev1.Secret {
	clusterUrl := viper.GetString("clusterUrl")

	retVal := &corev1.Secret{
		TypeMeta:   v1.TypeMeta{Kind: string(v1beta1.Secret)},
		ObjectMeta: v1.ObjectMeta{Name: "spot-oceancd-controller-token", Namespace: namespace},
		StringData: map[string]string{
			"token":   clusterToken.Token,
			"saasUrl": clusterUrl,
		},
	}

	retVal.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   corev1.GroupName,
		Version: "v1",
		Kind:    "Secret",
	})

	return retVal
}

func buildOperatorManagerConfigMap(config *operator.OMConfig) (*corev1.ConfigMap, error) {
	omConfig := config.GetOperatorManagerConfig()

	oceanCDBytes, err := yaml.Marshal(omConfig.OceanCDConfig)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load OceanCD config\n%w", err)
	}

	argoRolloutsBytes, err := yaml.Marshal(omConfig.ArgoRolloutsConfig)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load argo-rollouts config\n%w", err)
	}

	omConfigMap := &corev1.ConfigMap{
		TypeMeta:   v1.TypeMeta{Kind: string(v1beta1.ConfigMap)},
		ObjectMeta: v1.ObjectMeta{Name: "spot-oceancd-operator-manager-config", Namespace: config.OceanCDConfig.Namespace},
		Data: map[string]string{
			strings.TrimPrefix(component_configs.OceanCDConfigPath, "/"):      string(oceanCDBytes),
			strings.TrimPrefix(component_configs.ArgoRolloutsConfigPath, "/"): string(argoRolloutsBytes),
		},
	}

	omConfigMap.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   corev1.GroupName,
		Version: "v1",
		Kind:    "ConfigMap",
	})

	return omConfigMap, nil
}

func convertOperatorManagerConfigMap(configMap *corev1.ConfigMap) (*unstructured.Unstructured, error) {
	omConfigBytes, err := json.Marshal(configMap.DeepCopyObject())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to create operator manager configmap\n%w", err)
	}

	resource, err := helpers.ConvertToUnstructured(string(omConfigBytes))
	if err != nil {
		return nil, fmt.Errorf("error: Failed to create operator manager configmap\n%w", err)
	}

	return resource, nil
}

func convertOperatorManagerSecret(secret *corev1.Secret) (*unstructured.Unstructured, error) {
	secretBytes, err := json.Marshal(secret.DeepCopyObject())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to create operator manager secret\n%w", err)
	}

	resource, err := helpers.ConvertToUnstructured(string(secretBytes))
	if err != nil {
		return nil, fmt.Errorf("error: Failed to create operator manager secret\n%w", err)
	}

	return resource, nil
}
