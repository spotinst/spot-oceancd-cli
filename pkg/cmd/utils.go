package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

//like -f ./service_newservices.json
func ServiceSpecFromFile(file string) (model.Service, error) {
	if !strings.Contains(file, "microservice") {
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
	err = json.Unmarshal(bytes, &s)
	return s, nil

}
func EntitySpecFromFile(file string) (interface{}, error) {

	entityType := utils.GetEntityKind(file)
	var obj interface{}
	switch entityType {
	case "environment":
		obj = &model.EnvironmentSpec{}
	case "service":
		obj = &model.Service{}
	case "rolloutspec":
		obj = &model.RolloutSpec{}

	default:
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return obj, err
	}

	return obj, nil

}

func CreateEntity(ctx context.Context, obj interface{}, entityType string) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	apiTemplate := "https://api.spotinst.io/ocean/cd/%v"
	api := fmt.Sprintf(apiTemplate, entityType)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(obj).
		//	SetResult(model.OperationResponse{}).
		Post(api)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		return errors.New(fmt.Sprintf("response status is invalide ,  %v", status))
	}

	return nil
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
