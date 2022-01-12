package cmd

import "testing"

func TestInstallDir(result *testing.T) {
	dir := "../../samples/DemoApp"
	err := InstallFromDir(dir)
	if err != nil {
		result.Fatal(err)
	}
}
