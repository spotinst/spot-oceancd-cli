package oceancd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
	"spot-oceancd-cli/pkg/oceancd/model/operator"
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
		return nil, err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return nil, err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("error: Resource '%s/%s' does not exist", entityType, entityName)
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

	return fmt.Errorf("error: %s", errorMessage.(string))
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
			return rolloutInfo, fmt.Errorf("error: Rollout %s does not exist", rolloutId)
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutInfo, err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return rolloutInfo, err
	}

	if len(items) == 0 {
		return rolloutInfo, fmt.Errorf("error: Rollout %s does not exist", rolloutId)
	}

	if len(items) > 1 {
		return rolloutInfo, fmt.Errorf("error: Found more that 1 rollout resource for %s", rolloutId)
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
			return rolloutPhases, fmt.Errorf("error: Rollout phases for rollout %s do not exist", rolloutId)
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutPhases, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return rolloutPhases, err
	}

	if len(items) == 0 {
		return rolloutPhases, fmt.Errorf("error: Found no phases for rollout %s", rolloutId)
	}

	for _, item := range items {
		bytes, err := json.Marshal(item)
		if err != nil {
			return rolloutPhases, fmt.Errorf("error: Failed to parse a rollout phase - %w", err)
		}

		rolloutPhase := phase.Phase{}
		err = json.Unmarshal(bytes, &rolloutPhase)
		if err != nil {
			return rolloutPhases, fmt.Errorf("error: Failed to parse a rollout phase - %w", err)
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
			return rolloutVerifications, fmt.Errorf("error: Failed to parse a rollout verification - %w", err)
		}

		rolloutVerification := verification.Verification{}
		err = json.Unmarshal(bytes, &rolloutVerification)
		if err != nil {
			return rolloutVerifications, fmt.Errorf("error: Failed to parse a rollout verification - %w", err)
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
			return rolloutDefinition, errors.New(fmt.Sprintf("error: Rollout %s definition does not exist", rolloutId))
		}

		err = parseErrorFromResponse(response.Body())
		return rolloutDefinition, err
	}

	items, err := unmarshalEntityResponse(response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return rolloutDefinition, err
	}

	if len(items) == 0 {
		return rolloutDefinition, fmt.Errorf("error: Rollout %s definition does not exist", rolloutId)
	}

	if len(items) > 1 {
		return rolloutDefinition, fmt.Errorf("error: Found more that 1 rollout resource for %s", rolloutId)
	}

	bytes, err := json.Marshal(items[0])
	if err != nil {
		return rolloutDefinition, err
	}

	err = json.Unmarshal(bytes, &rolloutDefinition)

	return rolloutDefinition, nil
}

func GetOMInstallationManifests(_ context.Context, payload operator.OMManifestsRequest) (*operator.OMManifestsResponse, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/%v"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, "omInstaller")

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(payload).
		Post(apiUrl)

	if err != nil {
		return nil, err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return nil, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to parse response\n%w", err)
	}

	if len(items) != 1 {
		return nil, fmt.Errorf("error: Wrong number of installation items received, expected 1, got %d", len(items))
	}

	itemBytes, err := json.Marshal(items[0])
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load installation manifests response - %w", err)
	}

	output := &operator.OMManifestsResponse{}

	err = json.Unmarshal(itemBytes, output)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load installation manifests response - %w", err)
	}

	return output, nil
}

func DeleteCluster(_ context.Context, clusterId string) (*operator.DeleteClusterResponse, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/cluster/%s"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl, clusterId)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Delete(apiUrl)

	if err != nil {
		return nil, err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return nil, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to parse response - %w", err)
	}

	if len(items) != 1 {
		return nil, fmt.Errorf("error: Wrong number of items received, expected 1, got %d", len(items))
	}

	itemBytes, err := json.Marshal(items[0])
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	retVal := &operator.DeleteClusterResponse{}

	err = json.Unmarshal(itemBytes, retVal)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	return retVal, nil
}

func GetClusterManifests(_ context.Context) (*operator.ClusterManifestsMetadataResponse, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/cluster/manifest/metadata"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		Get(apiUrl)

	if err != nil {
		return nil, err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return nil, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to parse response - %w", err)
	}

	if len(items) != 1 {
		return nil, fmt.Errorf("error: Wrong number of items received, expected 1, got %d", len(items))
	}

	itemBytes, err := json.Marshal(items[0])
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	retVal := &operator.ClusterManifestsMetadataResponse{}

	err = json.Unmarshal(itemBytes, retVal)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	return retVal, nil
}

func CreateClusterToken(_ context.Context) (*operator.ClusterTokenResponse, error) {
	token := viper.GetString("token")
	baseUrl := viper.GetString("url")
	clusterId := viper.GetString("clusterId")

	client := resty.New()
	apiPrefixTemplate := "%v/ocean/cd/cluster/token"
	apiUrl := fmt.Sprintf(apiPrefixTemplate, baseUrl)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetQueryParams(map[string]string{"clusterId": clusterId}).
		Post(apiUrl)

	if err != nil {
		return nil, err
	}

	if status := response.StatusCode(); status != 200 {
		err = parseErrorFromResponse(response.Body())
		return nil, err
	}

	items, err := unmarshalEntityResponse(response.Body())
	if err != nil {
		return nil, fmt.Errorf("error: Failed to parse response - %w", err)
	}

	if len(items) != 1 {
		return nil, fmt.Errorf("error: Wrong number of items received, expected 1, got %d", len(items))
	}

	itemBytes, err := json.Marshal(items[0])
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	output := &operator.ClusterTokenResponse{}

	err = json.Unmarshal(itemBytes, output)
	if err != nil {
		return nil, fmt.Errorf("error: Failed to load response - %w", err)
	}

	return output, nil
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
