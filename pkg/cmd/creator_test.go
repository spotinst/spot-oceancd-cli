package cmd

import (
	"fmt"
	"testing"
)

func TestCreateServiceEntity(result *testing.T) {
	file := "../../samples/services_test1.json"
	s, err := CreateServiceFromManifest(file, true)

	if err != nil {
		result.Fatal(err)
	}

	if s.Name != "testservice" {
		result.Fatal(s)
	}
}
func TestUpdateServiceEntity(result *testing.T) {
	file := "../../samples/services_test1.json"
	s, err := CreateServiceFromManifest(file, false) //update

	if err != nil {
		result.Fatal(err)
	}

	fmt.Printf("created services is %v", s.Name)

}
