package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

const token = "79b8b542e613a96ae282c2e10cc328ef98afd89bd5a778078605e7808b8892ec"

func GetSpotContext(ctx context.Context) context.Context {
	testCtx := context.WithValue(ctx, "spottoken", token)

	return testCtx
}

func CreateServiceFromFile(ctx context.Context, file string) error {
	return nil
}
func CreateService(ctx context.Context, service *model.Service) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/microservice"

	//resourceUrl := fmt.Sprintf("%s/%s", api, service.Microservice.Name)

	bodyBytes, err := json.Marshal(service)
	serviceRequest := model.ServiceRequest{}
	serviceRequest.Microservice = service.ServiceMetadata
	if err != nil {
		return err
	}

	file := "../../yamls/body.json"
	err = os.WriteFile(file, bodyBytes, 0644)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(serviceRequest).
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
func ListServices(ctx context.Context) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/microservice"

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").Get(api)

	//	fmt.Printf("%v", string(response.Body()))

	if err != nil {
		return err
	}
	type MarshalHelper struct {
		Request  map[string]interface{} `json:"request"`
		Response struct {
			Items []model.Service `json:"items"`
		} `json:"response"`
	}
	helper := MarshalHelper{}
	err = json.Unmarshal(response.Body(), &helper)

	if err != nil {
		return err
	}
	utils.OutputServicesTable(helper.Response.Items)

	return nil
}
func PrintServices(items []model.Service) {
	for i, v := range items {
		fmt.Println("============================")
		color.Green("[%v]service %v\n", i, v.Name)
		fmt.Printf("[%v]service  labels %v\n", i, v.K8sResources.Labels)
		fmt.Printf("[%v]service  workload  type %v\n", i, v.K8sResources.Type)
		//fmt.Printf("\trollout  %v\n", v.Rollouts)
		fmt.Println("============================")

		//bytes, _ := json.Marshal(v)
		//	file := fmt.Sprintf("./yamls/services_%v.json", v.Name)
		//s.WriteFile(file, bytes, 0644)

	}
}
