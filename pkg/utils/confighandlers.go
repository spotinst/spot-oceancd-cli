package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	fp "path/filepath"
)

type Config struct {
	Filename   string
	SingleOnly bool
}

type Options struct {
	SingleOnly bool
	PathToConfig string
	DefaultData map[string]interface{}
}

type commandHandler func(ctx context.Context, resource map[string]interface{}) error

type ConfigHandler interface {
	Handle(context context.Context, commandHandler commandHandler) error
}

func NewConfigHandler(options Options) (ConfigHandler, error) {
	if options.PathToConfig != "" {
		fileExtension := fp.Ext(options.PathToConfig)[1:]

		switch fileExtension {
		case "json":
			return &JsonConfigHandler{options}, nil
		case "yaml", "yml":
			return &YamlConfigHandler{options}, nil
		default:
			return nil, errors.New("wrong file extension: Only Json and Yaml formats are supported")
		}
	}

	return &DefaultConfigHandler{config}, nil
}

type JsonConfigHandler struct {
	Options
}

func (h *JsonConfigHandler) Handle(ctx context.Context, commandHandler commandHandler) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = h.ToArrayOfMaps()
	if err != nil {
		resource, err = h.ToMap()
		if err != nil {
			return err
		}

		return commandHandler(ctx, resource)
	}

	if h.Options.SingleOnly && len(resources) > 1 {
		return errors.New("expected a single config but got more")
	}

	for _, resource = range resources {
		err = commandHandler(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *JsonConfigHandler) ToMap() (map[string]interface{}, error) {
	retVal := h.Options.DefaultData

	bytesContent, err := ioutil.ReadFile(h.Options.PathToConfig)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesContent, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}

func (h *JsonConfigHandler) ToArrayOfMaps() ([]map[string]interface{}, error) {
	var retVal []map[string]interface{}

	bytesContent, err := ioutil.ReadFile(h.Options.PathToConfig)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesContent, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}

type YamlConfigHandler struct {
	Options
}

func (h *YamlConfigHandler) Handle(ctx context.Context, commandHandler commandHandler) error {
	var resources []map[string]interface{}
	var resource map[string]interface{}
	var err error

	resources, err = h.ToArrayOfMaps(h.Filename)
	if err != nil {
		resources, err = h.ToMap(h.Filename)
		if err != nil {
			return err
		}
	}

	if h.SingleOnly && len(resources) > 1 {
		return errors.New("expected a single config but got more")
	}

	for _, resource = range resources {
		err = commandHandler(ctx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *YamlConfigHandler) ToMap(fileName string) ([]map[string]interface{}, error) {
	retVal := make([]map[string]interface{}, 0)

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(fileBytes))

	for {
		var resource map[string]interface{}
		if err = dec.Decode(&resource); err != nil {
			break
		}

		retVal = append(retVal, resource)
	}

	if err != io.EOF {
		return nil, err
	}

	return retVal, nil
}

func (h *YamlConfigHandler) ToArrayOfMaps(fileName string) ([]map[string]interface{}, error) {
	var retVal []map[string]interface{}

	bytesContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytesContent, &retVal)
	if err != nil {
		return nil, err
	}

	return retVal, err
}
