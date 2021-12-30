package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/verchol/applier/pkg/model"
	"github.com/verchol/applier/pkg/utils"
)

//like -f ./service_newservices.json
func ServiceSpecFromFile(file string) (model.Service, error) {
	if !strings.Contains(file, "service") {
		return model.Service{}, errors.New(fmt.Sprintf("file name should have service postfix %v", file))
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
func makeRequestObject(entity interface{}) interface{} {
	meta, ok := entity.(model.EntityMeta)
	if !ok {
		return nil
	}
	entityType := meta.GetEntityKind()
	var requestObj interface{}
	switch entityType {
	case model.EnvEntity:
		envRequest := &model.EnvironmentRequest{}
		envRequest.Envrionment = entity.(*model.EnvironmentSpec)
		requestObj = envRequest
	case model.ServiceEntity:

		svcRequest := &model.ServiceRequest{}
		svcRequest.Microservice = entity.(*model.Service)
		requestObj = svcRequest
	case model.RolloutSpecEntity:
		rSpecRequest := &model.RolloutSpecRequest{}
		rSpecRequest.Spec = entity.(*model.RolloutSpec)
		requestObj = rSpecRequest

	default:
		return nil
	}

	return requestObj
}
func EntitySpecFromFile(file string) (interface{}, error) {

	entityType := utils.GetEntityKindByFilename(file)
	var obj interface{}
	switch entityType {
	case model.EnvEntity:
		obj = &model.EnvironmentRequest{}
	case model.ServiceEntity:
		obj = &model.ServiceRequest{}
	case model.RolloutSpecEntity:
		obj = &model.RolloutSpecRequest{}

	default:
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(bytes, obj)
	if err != nil {
		return nil, err
	}

	return obj, err

}

//TODO create metadata
func GetEntityMetatData(obj interface{}) (*model.EntityMeta, error) {
	meta, ok := obj.(*model.EntityMeta)

	if !ok {
		return nil, errors.New("can't retrieve object metadata ")
	}

	return meta, nil

}

func GetEntitySpec(obj interface{}) (interface{}, error) {
	spec, ok := obj.(model.EntitySpec)

	if !ok {
		return nil, errors.New("can't retrieve object spec ")
	}

	return spec.GetEntitySpec(), nil
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
	b, _ := response.Request.RawRequest.GetBody()
	reqBody, _ := ioutil.ReadAll(b)
	fmt.Printf("[%v\n]", string(reqBody))
	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		errorMsg := color.New(color.FgRed).Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}
func UpdateEntity(ctx context.Context, obj interface{}, entityType string, entityName string) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	apiTemplate := "https://api.spotinst.io/ocean/cd/%v/%v"
	api := fmt.Sprintf(apiTemplate, entityType, entityName)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		SetBody(obj).
		//	SetResult(model.OperationResponse{}).
		Put(api)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		errorMsg := color.New(color.FgRed).Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}
func DeleteEntity(ctx context.Context, entityType string, entityName string) error {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	apiTemplate := "https://api.spotinst.io/ocean/cd/%v/%v"
	api := fmt.Sprintf(apiTemplate, entityType, entityName)

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		//	SetResult(model.OperationResponse{}).
		Delete(api)

	if err != nil {
		return err
	}

	if status := response.StatusCode(); status != 200 {
		errorMsg := color.New(color.FgRed).Sprintf("response status is invalid : %v\n", string(response.Body()))
		return errors.New(errorMsg)
	}

	return nil
}
func GetEntity(ctx context.Context, entityType string, entityName string, format string) (string, error) {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	apiTemplate := "https://api.spotinst.io/ocean/cd/%v/%v"
	api := fmt.Sprintf(apiTemplate, entityType, url.QueryEscape(entityName))

	response, err := client.R().
		SetAuthToken(token).
		ForceContentType("application/json").
		//	SetResult(model.OperationResponse{}).
		Get(api)

	if err != nil {
		return "", err
	}

	if status := response.StatusCode(); status != 200 {
		errorMsg := color.New(color.FgRed).Sprintf("response status is invalid : %v\n", string(response.Body()))
		return "", errors.New(errorMsg)
	}
	items, err := unmarshalEntityResponse(entityType, response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New("entity does not exist")
	}
	objToReturn := makeRequestObject(items[0])
	bytes, err := json.Marshal(objToReturn)
	if err != nil {
		return "", err
	}
	output := string(bytes)
	return output, nil
}
func ListEntities(ctx context.Context, entityType string) ([]interface{}, error) {
	testCtx := GetSpotContext(context.Background())
	token := testCtx.Value("spottoken").(string)

	client := resty.New()
	apiTemplate := "https://api.spotinst.io/ocean/cd/%v"
	api := fmt.Sprintf(apiTemplate, entityType)

	response, err := client.R().
		SetAuthToken(token).
		//	SetResult(model.OperationResponse{}).
		ForceContentType("application/json").Get(api)

	//	fmt.Printf("%v", string(response.Body()))

	if err != nil {
		return nil, err
	}

	items, err := unmarshalEntityResponse(entityType, response.Body()) //getListMarshallHelper(entityType)
	if err != nil {
		return nil, err
	}

	//OutputEntities(entityType, items)

	return items, nil
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
func OutputEntitiesWide(entityType string, items []interface{}, more interface{}) error {
	Headers := map[string][]string{}
	rsWide := fmt.Sprintf("%s_wide", model.RolloutSpecEntity)
	RolloutSpecHeaderWide := []string{"Name", "Environment", "Service", "Selector"}

	Headers[rsWide] = RolloutSpecHeaderWide
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(Headers[entityType])

	if len(items) == 0 {
		return errors.New("emplty list nothing ot show  ")
	}

	for _, item := range items {
		printer := item.(model.EntityPrinter)
		row := printer.Format("", more)
		table.Append(row)
	}

	table.Render() // Send output

	return nil

}
func OutputEntities(entityType string, items []interface{}) error {
	Headers := map[string][]string{}

	ServiceHeader := []string{"Name", "Labels", "Wokload Type"}
	RolloutSpecHeader := []string{"Name", "Environment", "Service"}

	EnvHeader := []string{"Name", "Cluster", "Namespace"}
	ClusterHeader := []string{"Name", "KubeVersion", "CtlVersion", "Node", "Pod"}

	Headers[model.ServiceEntity] = ServiceHeader
	Headers[model.RolloutSpecEntity] = RolloutSpecHeader
	Headers[model.EnvEntity] = EnvHeader
	Headers[model.ClusterEntity] = ClusterHeader

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(Headers[entityType])

	if len(items) == 0 {
		return errors.New("emplty list nothing ot show  ")
	}

	for _, item := range items {
		printer := item.(model.EntityPrinter)
		row := printer.Format("", "")
		table.Append(row)
	}

	table.Render() // Send output

	return nil

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
				Items []*model.Service `json:"items"`
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

	}

	errorMsg := fmt.Sprintf("unsupported entity %v", entityType)
	return nil, errors.New(errorMsg)

}
func WaitSpinner() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner
	time.Sleep(4 * time.Second)                                 // Run for some time to simulate work
	s.Stop()
}
