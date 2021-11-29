package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/verchol/applier/pkg/model"
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

func OutputServicesTable(items []model.Service) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Labels", "Wokload Type"})

	for _, s := range items {
		labels := ""
		for _, l := range s.K8sResources.Labels {
			label := fmt.Sprintf("%v=%v,", l.Key, l.Value)
			labels = labels + label
		}
		row := []string{s.Name, labels, s.K8sResources.Type}
		table.Append(row)
	}
	table.Render() // Send output
}

func OutputEnvsTable(items []model.EnvironmentSpec) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Cluster", "Namespace"})

	for _, e := range items {

		row := []string{e.Name, e.ClusterId, e.Namespace}
		table.Append(row)
	}
	table.Render() // Send output
}

func OutputSRolloutsTable(items []model.RolloutSpec) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Environment", "Service"})

	for _, s := range items {

		row := []string{s.Name, s.Environment, color.New(color.FgGreen).Sprint(s.Microservice)}
		table.Append(row)
	}
	table.Render() // Send output
}

func GetEntityKind(file string) string {
	if strings.Contains(file, "service") {
		return "service"
	}
	if strings.Contains(file, "env") {
		return "environment"
	}
	if strings.Contains(file, "rolloutspec") {
		return "rolloutspec"
	}

	return ""
}
