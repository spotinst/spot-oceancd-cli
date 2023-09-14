package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"spot-oceancd-cli/pkg/oceancd/model"
)

var (
	supportedFileTypes = map[string]bool{
		".json": true,
		".yml":  true,
		".yaml": true,
	}
)

func GetOceanCdEntityKindByName(entityType string) (string, error) {
	if entityType == "cluster" || entityType == "Cluster" || entityType == "clusters" || entityType == "Clusters" {
		return "", fmt.Errorf("error: Unrecognize resource type '%s'", entityType)
	}

	return GetEntityKindByName(entityType)
}

func GetEntityKindByName(entityType string) (string, error) {

	switch entityType {
	case model.VerificationProviderKind, model.VerificationProviderEntity, model.VerificationProviderEntityPlural,
		"VerificationProviders", "Verificationprovider", "verificationprovider", "verificationproviders",
		model.VerificationProviderShorts[0], model.VerificationProviderShorts[1]:
		entityType = model.VerificationProviderEntity
	case model.VerificationTemplateKind, model.VerificationTemplateEntity, model.VerificationTemplateEntityPlural,
		"VerificationTemplates", "Verificationtemplate", "verificationtemplate", "verificationtemplates",
		model.VerificationTemplateShorts[0], model.VerificationTemplateShorts[1]:
		entityType = model.VerificationTemplateEntity
	case model.RolloutSpecKind, model.RolloutSpecEntity, model.RolloutSpecEntityPlural, model.RolloutSpecShort,
		"RolloutSpecs", "Rolloutspec", "Rolloutspecs", "rolloutspec", "rolloutspecs":
		entityType = model.RolloutSpecEntity
	case model.StrategyKind, model.StrategyEntity, model.StrategyEntityPlural, "Strategies", model.StrategyEntityShorts[0],
		model.StrategyEntityShorts[1]:
		entityType = model.StrategyEntity
	case model.ClusterKind, model.ClusterEntity, model.ClusterEntityPlural, "Clusters":
		entityType = model.ClusterEntity
	default:
		return "", fmt.Errorf("error: Unrecognize resource type %s", entityType)
	}

	return entityType, nil
}

func ConvertEntityToJsonString(resource interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func ConvertEntitiesToJsonString(resources []interface{}) (string, error) {
	if len(resources) == 1 {
		return ConvertEntityToJsonString(resources[0])
	}

	jsonBytes, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func ConvertEntityToYamlString(resource interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(resource)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func ConvertEntitiesToYamlString(resources []interface{}) (string, error) {
	if len(resources) == 1 {
		return ConvertEntityToYamlString(resources[0])
	}

	yamlBytes, err := yaml.Marshal(resources)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func IsFileTypeSupported(fileType string) error {
	if fileType == "" {
		return fmt.Errorf("error: Unsupported file type. File must have an extension of type json or yaml")
	}

	if supportedFileTypes[fileType] == false {
		return fmt.Errorf("error: Unsupported file type %v. File must be of type json or yaml", fileType)
	}

	return nil
}

func GetNounForm(noun string, length int) string {
	if length == 0 || length > 1 {
		noun += "s"
	}

	return noun
}
