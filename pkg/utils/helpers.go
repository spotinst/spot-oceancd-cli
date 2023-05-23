package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"spot-oceancd-cli/pkg/oceancd/model"
	//"fmt"
	//"os"
	//"strings"
	//
	//"github.com/fatih/color"
	//"gopkg.in/yaml.v2"
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
		return "", errors.New(fmt.Sprintf("error: Unrecognize resource type '%s'", entityType))
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
		return "", errors.New(fmt.Sprintf("error: Unrecognize resource type %s", entityType))
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

func ConvertJsonFileToMap(fileName string) (map[string]interface{}, error) {
	var retVal map[string]interface{}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}

func ConvertJsonFileToArrayOfMaps(fileName string) ([]map[string]interface{}, error) {
	var retVal []map[string]interface{}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}

func ConvertYamlFileToMap(fileName string) ([]map[string]interface{}, error) {
	retVal := make([]map[string]interface{}, 0)

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(fileBytes))

	for {
		var resource map[string]interface{}
		if err = dec.Decode(&resource); err != nil {
			break
		}

		retVal = append(retVal, resource)
	}

	if err != io.EOF {
		return nil, err
	}

	return retVal, nil
}

func ConvertYamlFileToArrayOfMaps(fileName string) ([]map[string]interface{}, error) {
	var retVal []map[string]interface{}

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}

func IsFileTypeSupported(fileType string) error {
	if fileType == "" {
		fmt.Println("File must have an extension of type json or yaml")
		return errors.New("error: Unsupported file type")
	}

	if supportedFileTypes[fileType] == false {
		fmt.Println("File must be of type json or yaml")
		return errors.New(fmt.Sprintf("error: Unsupported file typedd %v", fileType))
	}

	return nil
}

func GetNounForm(noun string, length int) string {
	if length == 0 || length > 1 {
		noun += "s"
	}

	return noun
}
