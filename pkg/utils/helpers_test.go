package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verchol/applier/pkg/model"
)

func TestJsonToYaml(result *testing.T) {
	dir, _ := os.Getwd()
	fmt.Printf("current dir %v", dir)
	bytes, err := ioutil.ReadFile("../../yamls/rolloutSpec_0.json")
	r := model.RolloutSpec{}
	json.Unmarshal(bytes, &r)

	if err != nil {
		result.Fatalf("error - %v", err)
	}

	err = YamlToJson("yaml1.yaml", r)
	if err != nil {
		result.Fatalf("error - %v", err)
	}
}
func TestReadServiceManifest(result *testing.T) {
	file := "../../samples/services_test1.json"

	s, err := ServiceManifestFromFile(file)

	if err != nil {
		result.Fatal(err)
	}

	fmt.Printf("service %v", s)
}
func TestReadDir(result *testing.T) {
	const dir = "../../samples"
	list, err := ReadEntitiesDir(dir)
	if err != nil {
		result.Fatalf("error %v", err)
	}

	assert.NotEmpty(result, list.Services)
	assert.NotEmpty(result, list.Environments)
	assert.NotEmpty(result, list.Specs)

}
