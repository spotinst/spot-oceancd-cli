package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

func ListRolloutSpecs(ctx context.Context) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/rolloutSpec"

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").Get(api)

	//fmt.Printf("%v", string(response.Body()))

	if err != nil {
		return err
	}
	type MarshalHelper struct {
		Request  map[string]interface{} `json:"request"`
		Response struct {
			Items []model.RolloutSpec `json:"items"`
		} `json:"response"`
	}
	helper := MarshalHelper{}
	err = json.Unmarshal(response.Body(), &helper)

	if err != nil {
		return err
	}

	utils.OutputSRolloutsTable(helper.Response.Items)

	return nil
}
func CreateRollout(ctx context.Context, rolloutSpec *model.RolloutSpec) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	api := "https://api.spotinst.io/ocean/cd/rolloutSpec"

	//resourceUrl := fmt.Sprintf("%s/%s", api, service.Microservice.Name)

	specRequest := model.RolloutSpecRequest{}
	specRequest.Spec = rolloutSpec

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(specRequest).
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
func PrintRollouts(items []model.RolloutSpec) {
	for i, v := range items {

		fmt.Println("============================")
		color.Green("[%v]rollout name %v\n", i, v.Name)
		fmt.Printf("\tenvironment %v\n", v.Environment)
		//		fmt.Printf("[\trollout  namespace %v\n", v.Namespace)
		fmt.Println("\trollout type is rolling")
		for i1, v1 := range v.Strategy.Rolling.Verification.Phases {
			fmt.Printf("[%v]\t\t phase   %v\n", i1, v1)
		}
		fmt.Printf("\tnotification  %v\n", v.Notification)
		fmt.Printf("\tfailure policy  %v\n", v.FailurePolicy)

		bytes, err := json.Marshal(v)
		file := fmt.Sprintf("./yamls/rolloutSpec_%v.json", i)
		err = os.WriteFile(file, bytes, 0644)

		if err != nil {
			log.Fatalf("can't save fle %v", file)
		}
	}
}
