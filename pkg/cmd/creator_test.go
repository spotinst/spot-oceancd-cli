package cmd

import "testing"

func TestCreateServiceEntity(result *testing.T) {
	file := "../../samples/services_test1.json"
	s, err := CreateServiceFromManifest(file)

	if err != nil {
		result.Fatal(err)
	}

	if s.Name != "testservice" {
		result.Fatal(s)
	}
}
