package cmd

import (
	"context"
	"testing"

	"github.com/verchol/applier/pkg/model"
)

func TestListServices(result *testing.T) {
	err := ListServices(context.Background())

	if err != nil {
		result.Fatal(err)
	}
}
func TestListRollouts(result *testing.T) {
	err := ListRolloutSpecs(context.Background())

	if err != nil {
		result.Fatal(err)
	}
}

func TestCreateService(result *testing.T) {
	s := model.Service{}
	s.Microservice.Name = "oinventory-service"
	s.Microservice.K8sResources.ServiceWorkload.VersionLabelKey = "serviceVersion"
	s.Microservice.K8sResources.ServiceWorkload.Type = "deployment"
	s.Microservice.K8sResources.Labels = []model.ServiceLabel{{Key: "app", Value: "inventory"}}
	err := CreateService(context.Background(), &s)

	if err != nil {
		result.Fatal(err)
	}

}
