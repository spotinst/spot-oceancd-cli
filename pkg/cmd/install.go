package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
)

func init() {
	//logger = log.New(os.Stderr)
}
func InstallFromDir(dir string) error {

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Printf("installing resource from %s\n", color.GreenString(f.Name()))
	}

	return nil
}
func Install(resources []string) error {

	return nil
}
func InstallJob(jobManifest string) error {
	return nil
}
func CheckInstall(job string) bool {
	return true
}

func Logs(job string) {

}
func TestInstall() {

}
