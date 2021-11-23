package model

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
)

const token = "79b8b542e613a96ae282c2e10cc328ef98afd89bd5a778078605e7808b8892ec"

func InitTest(ctx context.Context) context.Context {
	testCtx := context.WithValue(ctx, "spottoken", token)

	return testCtx
}
func TestEnvironmentCreate(result *testing.T) {

	testCtx := InitTest(context.Background())
	token := testCtx.Value("spottoken").(string)
	client := resty.New()
	//	api := "https://api.spotinst.io/ocean/cd/environment"

	env := Environment{}
	env.Envrionment = EnvironmentSpec{}
	env.Envrionment.ClusterId = "olegv"
	env.Envrionment.Name = "generatedenv"
	env.Envrionment.Namespace = "default"

	body, err := json.Marshal(env)

	if err != nil {
		result.FailNow()
	}

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").
		SetBody(body).
		Put("https://api.spotinst.io/ocean/cd/environment/" + env.Envrionment.Name)

	fmt.Printf("%v", string(response.Body()))
}

func TestServiceCreate(result *testing.T) {
	testCtx := InitTest(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/microservice"

	service := Service{}
	service.Microservice.Name = "service_test_wed"
	//service.Labels = []ServiceLabel{{Key: "app", Value: "test1"}}

	body, err := json.Marshal(service)
	resourceUrl := fmt.Sprintf("%s/%s", api, service.Microservice.Name)

	if err != nil {
		result.FailNow()
	}

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").
		SetBody(body).
		Post(resourceUrl)

	fmt.Printf("%v", string(response.Body()))

	if err != nil {
		result.FailNow()
	}

}

/*
{
  "name": "AmirFirstMicroservice",
  "k8sResources": {
    "workload": {
      "type": "deployment",
      "labels": [
        {
          "key": "app",
          "value": "AmirFirstMicroservice"
        }
      ],
      "versionLabelKey": "ms-version"
    }
  },
  "createdAt": "2021-11-04T10: 06: 13.803Z",
  "updatedAt": "2021-11-04T10: 06: 13.803Z"
}
*/
func TestListServices(result *testing.T) {
	testCtx := InitTest(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/microservice"

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").Get(api)

	fmt.Printf("%v", string(response.Body()))

	if err != nil {
		result.FailNow()
	}
	type MarshalHelper struct {
		Request  map[string]interface{} `json:"request"`
		Response struct {
			Items []Service `json:"items"`
		} `json:"response"`
	}
	helper := MarshalHelper{}
	err = json.Unmarshal(response.Body(), &helper)

	if err != nil {
		result.FailNow()
	}

	for i, v := range helper.Response.Items {
		fmt.Println("============================")
		fmt.Printf("[%v]service %v\n", i, v.Microservice.Name)
		//	fmt.Printf("\trollout  %v\n", v.Rollouts)
		fmt.Println("============================")

	}
}
