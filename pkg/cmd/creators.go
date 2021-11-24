package cmd

import (
	"context"

	"github.com/verchol/applier/pkg/model"
)

func CreateServiceFromManifest(file string) (*model.Service, error) {

	s, err := ServiceSpecFromFile(file)
	if err != nil {
		return nil, err
	}

	err = CreateEntity(context.Background(), s, "service")
	if err != nil {
		return nil, err
	}

	return &s, nil
}
func CreateRolloutFromManifest(file string) (*model.Service, error) {

	s, err := ServiceSpecFromFile(file)
	if err != nil {
		return nil, err
	}

	err = CreateEntity(context.Background(), s, "rolloutspec")
	if err != nil {
		return nil, err
	}

	return &s, nil
}
