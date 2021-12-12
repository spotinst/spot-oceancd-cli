package cmd

import (
	"context"
	"testing"

	"github.com/verchol/applier/pkg/model"
)

func TestListServices(result *testing.T) {
	_, err := ListServices(context.Background())

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
	s.Name = "olegsvc1"
	s.K8sResources.ServiceWorkload.VersionLabelKey = "serviceVersion"
	s.K8sResources.ServiceWorkload.Type = "deployment"
	s.K8sResources.Labels = []model.ServiceLabel{{Key: "app", Value: "inventory"}}
	err := CreateService(context.Background(), &s)

	if err != nil {
		result.Fatal(err)
	}

}

func TestCreateRollout(result *testing.T) {
	s := model.RolloutSpec{}
	s.Name = "firstrollout"
	s.Environment = "dev"
	s.Microservice = "mymycroservice"
	err := CreateRollout(context.Background(), &s)

	if err != nil {
		result.Fatal(err)
	}

}
