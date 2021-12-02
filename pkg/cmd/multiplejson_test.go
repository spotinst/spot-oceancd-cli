package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

func TestMultipleJsonFiles(result *testing.T) {
	file := "../../samples/all.json"
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		result.Fatal(err)
	}
	reader := strings.NewReader(string(bytes))
	decoder := json.NewDecoder(reader)
	resources := []string{}
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			//	log.Fatal(err)
		}
		fmt.Printf("%T: %v", t, t)

		kind := reflect.TypeOf(t).String()

		if kind == "string" {
			val := reflect.ValueOf(t).String()
			r, err := utils.GetEntityKindByName(val)
			if err == nil {
				resources = append(resources, r)
			}
		}
		if decoder.More() {
			fmt.Printf(" (more)")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("resources %v ", resources)
}
func TestMultipleJsonUnmarshal(result *testing.T) {
	file := "../../samples/all.json"
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		result.Fatal(err)
	}
	reader := strings.NewReader(string(bytes))
	envs := []model.EnvironmentRequest{}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(envs)
	if err != nil {
		result.Fail()
	}

}
