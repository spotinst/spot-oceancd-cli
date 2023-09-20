/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path/filepath"
	"spot-oceancd-cli/commons/configs"
	"spot-oceancd-cli/commons/handlers/cluster"
	"spot-oceancd-cli/commons/helpers"
	"spot-oceancd-cli/pkg/oceancd"
	"spot-oceancd-cli/pkg/oceancd/model/operator"
	"spot-oceancd-cli/pkg/utils"
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
				fmt.Printf("failed to install operator: %s", err)
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

	if pathToConfig != "" {
		fmt.Printf("the %s config is being used.\n", pathToConfig)
		configHandler, err := utils.NewConfigHandler(pathToConfig)
		if err != nil {
			return fmt.Errorf("failed to initiate config handler: %w", err)
		}

		err = configHandler.Handle(ctx, installOperator)
		if err != nil {
			return fmt.Errorf("failed to install operator: %w", err)
		}

		fmt.Printf("operator installation finished succesfully.\n")
		return nil

	}

	fmt.Printf("config wasn't provided. The default config is being used.\n")

	defaultConfigBytes, err := json.Marshal(operator.DefaultInstallationConfig())
	if err != nil {
		return fmt.Errorf("failed to marshal default operator manager config: %w", err)
	}

	defaultConfigData := map[string]interface{}{}
	if err := json.Unmarshal(defaultConfigBytes, &defaultConfigData); err != nil {
		return fmt.Errorf("failed to unmarshal default operator manager config: %w", err)
	}

	err = installOperator(ctx, defaultConfigData)
	if err != nil {
		return fmt.Errorf("failed to install operator: %w", err)
	}

	fmt.Printf("operator installation finished succesfully.\n")
	return nil
}

func validateOperatorInstallFlags(cmd *cobra.Command) error {
	pathToConfig, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse --config flag: %w", err)
	}

	if cmd.Flags().Lookup("config").Changed == false {
		return nil
	}

	if pathToConfig == "" {
		return fmt.Errorf("path to config file must be specified\n")
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

	if err := fetchAndApplyManifests(ctx, config); err != nil {
		return fmt.Errorf("failed to fetch and apply installation manifests: %w", err)
	}

	operatorManagerConfigMap, err := buildOperatorManagerConfigMap(config)
	if err != nil {
		return fmt.Errorf("failed to build operator manager ConfigMap: %w", err)
	}

	resource, err := convertOperatorManagerConfigMap(operatorManagerConfigMap)
	if err != nil {
		return fmt.Errorf("failed to convert operator manager ConfigMap: %w", err)
	}

	applyHandler := cluster.BaseApplyHandler{}
	if err := applyHandler.Apply(resource); err != nil {
		return fmt.Errorf("failed to apply operator manager ConfigMap: %w", err)
	}

	return nil
}

func fetchAndApplyManifests(ctx context.Context, config *operator.InstallationConfig) error {
	manifestSets, err := oceancd.InstallOperator(ctx, operator.NewInstallationPayload(config))
	if err != nil {
		return fmt.Errorf("failed to fetch installation resources: %w", err)
	}

	if err = applyAndPatch(manifestSets.Argo); err != nil {
		return fmt.Errorf("failed to apply and patch argo-rollouts manifest set: %w", err)
	}

	if err = applyAndPatch(manifestSets.OceanCD); err != nil {
		return fmt.Errorf("failed to apply and patch OceanCD manifest set: %w", err)
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
			strings.TrimPrefix(configs.OceanCDConfigPath, "/"):      string(oceanCDBytes),
			strings.TrimPrefix(configs.ArgoRolloutsConfigPath, "/"): string(argoRolloutsBytes),
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

func applyAndPatch(set operator.ManifestSet) error {
	for _, manifest := range set.Appliable {
		resource, err := helpers.ConvertToUnstructured(manifest)
		if err != nil {
			return fmt.Errorf("failed to convert manifest to unstructured: %w; manifest: %s", err, manifest)
		}

		applyHandler := cluster.BaseApplyHandler{}
		if err = applyHandler.Apply(resource); err != nil {
			return fmt.Errorf("failed to apply manifest: %w; manifest: %s", err, manifest)
		}
	}

	for _, manifest := range set.Patchable {
		resource, err := helpers.ConvertToUnstructured(manifest)
		if err != nil {
			return fmt.Errorf("failed to convert manifest to unstructured: %w; manifest: %s", err, manifest)
		}

		payload := &cluster.PatchPayload{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
			Kind:      resource.GetKind(),
			PatchBody: manifest,
		}

		patchHandler := cluster.BasePatchHandler{}
		if err = patchHandler.Patch(payload); err != nil {
			return fmt.Errorf("failed to patch manifest: %w; manifest: %s", err, manifest)
		}
	}

	return nil
}
