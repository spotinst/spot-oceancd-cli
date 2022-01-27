package oceancd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
	"spot-oceancd-cli/pkg/oceancd/model"
)

func CreateResource(ctx context.Context, entityType string, resourceToCreate interface{}) error {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, entityType)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(resourceToCreate).
		Post(apiUrl)

	// todo caduri - dont forget to remove comment
	//b, _ := response.Request.RawRequest.GetBody()
	//reqBody, _ := ioutil.ReadAll(b)
	//fmt.Printf("[%v\n]", string(reqBody))
	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		// todo caduri - parse error from server
		errorMsg := fmt.Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}

func UpdateResource(ctx context.Context, entityType string, entityName string, resourceToUpdate interface{}) error {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, entityType, url.QueryEscape(entityName))

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(resourceToUpdate).
		//	SetResult(model.OperationResponse{}).
		Put(apiUrl)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		// todo caduri - parse error from server
		errorMsg := fmt.Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}


func DeleteEntity(ctx context.Context, entityType string, entityName string) error {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, entityType, url.QueryEscape(entityName))

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Delete(apiUrl)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		errorMsg := fmt.Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}

func GetEntity(ctx context.Context, entityType string, entityName string) (interface{}, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, entityType, url.QueryEscape(entityName))

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return "", err
	}

	if status := response.StatusCode(); status != 200 {
		// todo caduri - parse error from server
		errorMsg := fmt.Sprintf("response status is invalid : %v\n", string(response.Body()))
		return "", errors.New(errorMsg)
	}

	items, err := unmarshalEntityResponse(entityType, response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New("entity does not exist")
	}

	return items[0], nil
}
func ListEntities(ctx context.Context, entityType string) ([]interface{}, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, entityType)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").Get(apiUrl)
	if err != nil {
		return nil, err
	}

	items, err := unmarshalEntityResponse(entityType, response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func unmarshalEntityResponse(entityType string, response []byte) ([]interface{}, error) {

	switch entityType {
	case model.EnvEntity:

		type MarshalHelper struct {
			Request  map[string]interface{} `json:"request"`
			Response struct {
				Items []*model.EnvironmentSpec `json:"items"`
			} `json:"response"`
		}
		helper := MarshalHelper{} //getListMarshallHelper(entityType)

		err := json.Unmarshal(response, &helper)
		if err != nil {
			return nil, err
		}
		items := helper.Response.Items
		b := make([]interface{}, len(items))
		for i := range items {
			b[i] = items[i]
		}
		return b, nil
	case model.ServiceEntity:

		type MarshalHelper struct {
			Request  map[string]interface{} `json:"request"`
			Response struct {
				Items []*model.Microservice `json:"items"`
			} `json:"response"`
		}
		helper := MarshalHelper{}
		err := json.Unmarshal(response, &helper)

		if err != nil {
			return nil, err
		}
		items := helper.Response.Items
		b := make([]interface{}, len(items))
		for i := range items {
			b[i] = items[i]
		}
		return b, nil
	case model.RolloutSpecEntity:

		type MarshalHelper struct {
			Request  map[string]interface{} `json:"request"`
			Response struct {
				Items []*model.RolloutSpec `json:"items"`
			} `json:"response"`
		}
		helper := MarshalHelper{}
		err := json.Unmarshal(response, &helper)
		if err != nil {
			return nil, err
		}
		items := helper.Response.Items
		b := make([]interface{}, len(items))
		for i := range items {
			b[i] = items[i]
		}
		return b, nil
	case model.ClusterEntity:

		type MarshalHelper struct {
			Request  map[string]interface{} `json:"request"`
			Response struct {
				Items []*model.ClusterSpec `json:"items"`
			} `json:"response"`
		}
		helper := MarshalHelper{}
		err := json.Unmarshal(response, &helper)
		if err != nil {
			return nil, err
		}
		items := helper.Response.Items
		b := make([]interface{}, len(items))
		for i := range items {
			b[i] = items[i]
		}
		return b, nil
	case model.NotificationProviderEntity:
		type MarshalHelper struct {
			Request  map[string]interface{} `json:"request"`
			Response struct {
				Items []*model.NotificationProviderSpec `json:"items"`
			} `json:"response"`
		}
		helper := MarshalHelper{}
		err := json.Unmarshal(response, &helper)
		if err != nil {
			return nil, err
		}
		items := helper.Response.Items
		b := make([]interface{}, len(items))
		for i := range items {
			b[i] = items[i]
		}
		return b, nil
	}

	errorMsg := fmt.Sprintf("unsupported entity %v", entityType)
	return nil, errors.New(errorMsg)

}
