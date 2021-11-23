package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/verchol/applier/pkg/model"
)

//like -f ./service_newservices.json
func ServiceManifestFromFile(file string) (model.Service, error) {
	if !strings.Contains(file, "service") {
		return model.Service{}, errors.New(fmt.Sprintf("file name should have service posfix %v", file))
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return model.Service{}, err
	}
	s := model.Service{}
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return model.Service{}, err
	}

	return s, nil

}
func ReadEntitiesDir(dir string) (model.EntityList, error) {
	files, err := ioutil.ReadDir(dir)
	list := model.EntityList{}
	if err != nil {
		return list, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "service") {
			bytes, err := ioutil.ReadFile(dir + "//" + f.Name())
			if err != nil {
				continue
			}
			s := model.ServiceRequest{}
			err = json.Unmarshal(bytes, &s)
			if err != nil {
				continue
			}
			list.Services = append(list.Services, s)

		}
		if strings.Contains(f.Name(), "environment") {
			bytes, err := ioutil.ReadFile(dir + "//" + f.Name())
			if err != nil {
				continue
			}
			e := model.EnvironmentSpec{}
			err = json.Unmarshal(bytes, &e)
			if err != nil {
				continue
			}
			list.Environments = append(list.Environments, e)

		}
		if strings.Contains(f.Name(), "rolloutspec") {
			bytes, err := ioutil.ReadFile(dir + "//" + f.Name())
			if err != nil {
				continue
			}
			r := model.RolloutSpec{}
			err = json.Unmarshal(bytes, &r)
			if err != nil {
				continue
			}
			list.Specs = append(list.Specs, r)

		}
	}

	return list, nil

}
