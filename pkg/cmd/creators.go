package cmd

import (
	"context"

	"github.com/verchol/applier/pkg/model"
)

func CreateServiceFromManifest(file string, create bool) (*model.Service, error) {

	s, err := EntitySpecFromFile(file)
	if err != nil {
		return nil, err
	}

	serviceRequest := s.(*model.ServiceRequest)
	service := serviceRequest.Microservice

	if create {
		err = CreateEntity(context.Background(), s, "microservice")
	} else {
		err = UpdateEntity(context.Background(), s, "microservice", service.Name)
	}
	if err != nil {
		return nil, err
	}

	return service, nil
}
func CreateRolloutFromManifest(file string, create bool) (*model.RolloutSpec, error) {

	s, err := EntitySpecFromFile(file)
	if err != nil {
		return nil, err
	}

	request := s.(*model.RolloutSpecRequest)

	if create {
		err = CreateEntity(context.Background(), request, "rolloutSpec")
	} else {
		err = UpdateEntity(context.Background(), request, "rolloutSpec", request.Spec.Name)
	}
	if err != nil {
		return nil, err
	}
	obj := request.Spec
	return obj, nil

}

func CreateEnvironmentFromManifest(file string, create bool) (*model.EnvironmentSpec, error) {

	s, err := EntitySpecFromFile(file)
	if err != nil {
		return nil, err
	}

	request := s.(*model.EnvironmentRequest)

	if create {
		err = CreateEntity(context.Background(), request, "environment")
	} else {
		err = UpdateEntity(context.Background(), request, "environment", request.Envrionment.Name)
	}
	if err != nil {
		return nil, err
	}

	obj := request.Envrionment

	return obj, nil

}
