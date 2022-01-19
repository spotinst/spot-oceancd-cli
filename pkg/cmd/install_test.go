package cmd

import (
	"testing"

	"github.com/verchol/applier/pkg/utils"
)

func TestInstallDir(result *testing.T) {
	dir := "../../samples/DemoApp"
	err := InstallFromDir(dir)
	if err != nil {
		result.Fatal(err)
	}
}

func TestInstallNamespace(result *testing.T) {
	i, err := NewInstaller()
	i.Namespace = "olegtest7"
	if err != nil {
		result.Fatal(err)
	}
	token := "79b8b542e613a96ae282c2e10cc328ef98afd89bd5a778078605e7808b8892ec"
	url := "https://api.spotinst.io/ocean/cd/clusterInstaller"
	clusterId := "temp_oleg_2"
	manifest, err := utils.BringInstallScript(url, clusterId, token)
	if err != nil {
		result.Fatal(err)
	}
	resources, err := utils.RunLocalBuilder(manifest)
	if err != nil {
		result.Fatal(err)
	}
	err = i.Install("olegtest8", resources)

	if err != nil {
		result.Fatal(err)
	}
}
