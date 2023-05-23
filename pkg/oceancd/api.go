package oceancd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
	"spot-oceancd-cli/pkg/oceancd/model/phase"
	"spot-oceancd-cli/pkg/oceancd/model/rollout"
	"spot-oceancd-cli/pkg/oceancd/model/verification"
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

func GetRollout(rolloutId string) (rollout.Rollout, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")
	rolloutInfo := rollout.Rollout{}

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/rollout/%s/status"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, rolloutId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return rolloutInfo, err
	}

	if response.StatusCode() != 200 {
		if response.StatusCode() == 400 {
			return rolloutInfo, errors.New(fmt.Sprintf("error: rollout %s does not exist", rolloutId))
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutInfo, err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return rolloutInfo, err
	}

	if len(items) == 0 {
		return rolloutInfo, errors.New("resource does not exist")
	}

	if len(items) > 1 {
		return rolloutInfo, errors.New(fmt.Sprintf("found more that 1 rollout resource: %+v", items))
	}

	bytes, err := json.Marshal(items[0])
	if err != nil {
		return rolloutInfo, err
	}

	err = json.Unmarshal(bytes, &rolloutInfo)
	if err != nil {
		return rolloutInfo, err
	}

	return rolloutInfo, nil
}

func GetRolloutPhases(rolloutId string) ([]phase.Phase, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")
	rolloutPhases := make([]phase.Phase, 0)

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/rollout/%s/phase"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, rolloutId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return rolloutPhases, err
	}

	if status := response.StatusCode(); status != 200 {
		if response.StatusCode() == 400 {
			return rolloutPhases, errors.New(fmt.Sprintf("error: Rollout phases for rollout %s do not exist", rolloutId))
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutPhases, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return rolloutPhases, err
	}

	if len(items) == 0 {
		return rolloutPhases, errors.New(fmt.Sprintf("found no phases for rollout %s", rolloutId))
	}

	for _, item := range items {
		bytes, err := json.Marshal(item)
		if err != nil {
			return rolloutPhases, fmt.Errorf("failed to parse a rollout phase: %s", err)
		}

		rolloutPhase := phase.Phase{}
		err = json.Unmarshal(bytes, &rolloutPhase)
		if err != nil {
			return rolloutPhases, fmt.Errorf("failed to parse a rollout phase: %s", err)
		}

		rolloutPhases = append(rolloutPhases, rolloutPhase)
	}

	return rolloutPhases, nil
}

func GetRolloutVerifications(rolloutId string) ([]verification.Verification, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")
	rolloutVerifications := make([]verification.Verification, 0)

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/rollout/%s/verification"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, rolloutId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return rolloutVerifications, err
	}

	if status := response.StatusCode(); status != 200 {
		if response.StatusCode() == 400 {
			return rolloutVerifications, errors.New(fmt.Sprintf("error: Rollout verifications for rollout %s do not exist", rolloutId))
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutVerifications, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return rolloutVerifications, err
	}

	if len(items) == 0 {
		return rolloutVerifications, nil
	}

	for _, item := range items {
		bytes, err := json.Marshal(item)
		if err != nil {
			return rolloutVerifications, fmt.Errorf("failed to parse a rollout verification: %s", err)
		}

		rolloutVerification := verification.Verification{}
		err = json.Unmarshal(bytes, &rolloutVerification)
		if err != nil {
			return rolloutVerifications, fmt.Errorf("failed to parse a rollout verification: %s", err)
		}

		rolloutVerifications = append(rolloutVerifications, rolloutVerification)
	}

	return rolloutVerifications, nil
}

func SendWorkloadAction(pathParams PathParams, queryParams QueryParams) error {
	token := viper.GetString("token")
	client := resty.New()

	response, err := client.R().
		SetAuthToken(token).
		SetQueryParams(queryParams).
		SetPathParams(pathParams).
		Put(buildWorkloadApiUrl(pathParams))

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return err
	}

	return nil
}

func GetRolloutDefinition(rolloutId string) (map[string]interface{}, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")
	var rolloutDefinition map[string]interface{}

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/rollout/%s/definition"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, rolloutId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return rolloutDefinition, err
	}

	if response.StatusCode() != 200 {
		if response.StatusCode() == 400 {
			return rolloutDefinition, errors.New(fmt.Sprintf("error: rollout %s does not exist", rolloutId))
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutDefinition, err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return rolloutDefinition, err
	}

	if len(items) == 0 {
		return rolloutDefinition, errors.New("resource does not exist")
	}

	if len(items) > 1 {
		return rolloutDefinition, errors.New(fmt.Sprintf("found more that 1 rollout resource: %+v", items))
	}

	bytes, err := json.Marshal(items[0])
	if err != nil {
		return rolloutDefinition, err
	}

	err = json.Unmarshal(bytes, &rolloutDefinition)

	return rolloutDefinition, nil
}

func buildWorkloadApiUrl(params PathParams) string {
	urlTemplate := fmt.Sprintf("%s/ocean/cd/workload/{spotDeploymentName}/namespace/{namespace}",
		viper.GetString("url"))

	if params["action"] != RestartAction {
		urlTemplate += "/revision/{revisionId}"
	}

	urlTemplate += "/{action}"
	return urlTemplate
}
