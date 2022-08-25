package oceancd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
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

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return err
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
		err = parseErrorFromResponse(response.Body())
		return err
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
		err = parseErrorFromResponse(response.Body())
		return err
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
		err = parseErrorFromResponse(response.Body())
		return "", err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New("resource does not exist")
	}

	return items[0], nil
}

func parseErrorFromResponse(body []byte) error {
	var spotResponse map[string]interface{}
	err := json.Unmarshal(body, &spotResponse)
	if err != nil {
		return err
	}

	innerResponse, isResponseKeyExist := spotResponse["response"]
	if isResponseKeyExist == false {
		return errors.New("error: Unknown server error")
	}

	innerErrors, isErrorsKeyExist := innerResponse.(map[string]interface{})["errors"]
	if isErrorsKeyExist == false || len(innerErrors.([]interface{})) == 0 {
		return errors.New("error: Unknown server error")
	}

	errorMessage, isMessageExists := innerErrors.([]interface{})[0].(map[string]interface{})["message"]
	if isMessageExists == false {
		return errors.New("error: Unknown server error")
	}

	return errors.New(errorMessage.(string))
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

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return nil, err
	}

	return items, nil
}

func unmarshalEntityResponse(response []byte) ([]interface{}, error) {

	type MarshalHelper struct {
		Request  map[string]interface{} `json:"request"`
		Response struct {
			Items []map[string]interface{} `json:"items"`
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

}

func SendRolloutAction(rolloutId string, body map[string]string) error {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/rollout/%s"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, rolloutId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(body).
		Put(apiUrl)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return err
	}

	return nil
}
