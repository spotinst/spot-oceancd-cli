/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
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

// operatorInstallCmd represents the operator install command
var (
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
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOperatorInstallFlags(cmd)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.NoArgs(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := runOperatorInstallCmd(context.Background(), cmd); err != nil {
				fmt.Printf("failed to install operator: %s\n", err)
			}
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

}

func runOperatorInstallCmd(ctx context.Context, cmd *cobra.Command) error {
	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse --config flag: %w", err)
	}

	installOptions := utils.Options{
		SingleOnly:   true,
		PathToConfig: pathToConfig,
	}

	configHandler, err := utils.NewConfigHandler(installOptions)
	if err != nil {
		return fmt.Errorf("failed to initiate config handler: %w", err)
	}

	err = configHandler.Handle(ctx, installOperator)
	if err != nil {
		return fmt.Errorf("failed to execute config handler: %w", err)
	}

	fmt.Printf("operator installation finished succesfully.\n")
	return nil
}

func validateOperatorInstallFlags(cmd *cobra.Command) error {
	if cmd.Flags().Lookup("config").Changed == false {
		return nil
	}

	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse --config flag: %w", err)
	}

	if pathToConfig == "" {
		return fmt.Errorf("path to config file must be specified")
	}

	fileExtensionWithDot := filepath.Ext(pathToConfig)
	if err := utils.IsFileTypeSupported(fileExtensionWithDot); err != nil {
		return err
	}

	return nil
}

func installOperator(ctx context.Context, data map[string]interface{}) error {

	config, err := operator.NewInstallationConfig(data)
	if err != nil {
		return fmt.Errorf("failed to initialize installation config: %w", err)
	}

	payload := operator.NewInstallationPayload(config)
	output, err := oceancd.GetOMInstallationManifests(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to fetch installation resources: %w", err)
	}

	resources, err := helpers.ConvertToUnstructuredSlice(output.OM.Manifests)
	if err != nil {
		return fmt.Errorf("failed to convert manifests to unstructured: %w", err)
	}

	operatorManagerConfigMap, err := buildOperatorManagerConfigMap(config)
	if err != nil {
		return fmt.Errorf("failed to build operator manager ConfigMap: %w", err)
	}

	configMapResource, err := convertOperatorManagerConfigMap(operatorManagerConfigMap)
	if err != nil {
		return fmt.Errorf("failed to convert operator manager ConfigMap: %w", err)
	}

	resources = append(resources, configMapResource)

	applyHandler := cluster.BaseApplyHandler{}
	for _, resource := range resources {
		if err := applyHandler.Apply(resource); err != nil {
			return fmt.Errorf("failed to apply operator manager ConfigMap: %w", err)
		}
	}

	return nil
}

func buildOperatorManagerConfigMap(config *operator.InstallationConfig) (*corev1.ConfigMap, error) {
	omConfig := config.GetOperatorManagerConfig()

	oceanCDBytes, err := yaml.Marshal(omConfig.OceanCDConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OceanCD config: %w", err)
	}

	argoRolloutsBytes, err := yaml.Marshal(omConfig.ArgoRolloutsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal argo-rollouts config: %w", err)
	}

	omConfigMap := &corev1.ConfigMap{
		TypeMeta:   v1.TypeMeta{Kind: string(v1beta1.ConfigMap)},
		ObjectMeta: v1.ObjectMeta{Name: "oceancd-operator-manager", Namespace: config.OceanCDConfig.Namespace},
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
		return nil, fmt.Errorf("failed to marshal operator manager configmap: %w", err)
	}

	resource, err := helpers.ConvertToUnstructured(string(omConfigBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to convert operator manager configmap: %w", err)
	}

	return resource, nil
}
