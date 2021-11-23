package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

func SaveToFile(file string, bytes []byte) error {
	return nil
}

func YamlToJson(yamlFile string, o interface{}) error {
	bytes, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	err = os.WriteFile(yamlFile, bytes, 0644)

	if err != nil {
		return err
	}

	return nil
}
