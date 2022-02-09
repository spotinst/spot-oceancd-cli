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
	case "Microservice", "microservice", "Microservices", "microservices", "ms":
		entityType = model.ServiceEntity
	case "RolloutSpec", "RolloutSpecs", "rolloutSpec", "rolloutSpecs", "Rolloutspec", "Rolloutspecs", "rolloutspec", "rolloutspecs", "rs":
		entityType = model.RolloutSpecEntity
	case "Environment", "Environments", "environment", "environments", "env", "envs":
		entityType = model.EnvEntity
	case "Cluster", "Clusters", "cluster", "clusters":
		entityType = model.ClusterEntity
	case "NotificationProvider", "NotificationProviders", "notificationProvider", "notificationProviders", "np":
		entityType = model.NotificationProviderEntity
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
		return errors.New("error: Unsupported file type")
	}

	return nil
}
