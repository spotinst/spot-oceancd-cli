package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/verchol/applier/pkg/model"
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

	fmt.Printf("%v", string(response.Body()))

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

	for i, v := range helper.Response.Items {

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

	return nil
}
