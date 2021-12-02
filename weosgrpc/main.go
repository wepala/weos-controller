package weosgrpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	weoscontroller "github.com/wepala/weos-controller"
)

func InitalizeGrpc(ctx *context.Context, api weoscontroller.GRPCAPIInterface, apiConfig string) *context.Context {
	var content []byte
	var err error
	//try load file if it's a yaml file otherwise it's the contents of a yaml file WEOS-1009
	if strings.Contains(apiConfig, ".yaml") || strings.Contains(apiConfig, "/yml") {
		content, err = ioutil.ReadFile(apiConfig)
		if err != nil {
			//e.Logger.Fatalf("error loading api specification '%s'", err)
		}
	} else {
		content = []byte(apiConfig)
	}

	//change the $ref to another marker so that it doesn't get considered an environment variable WECON-1
	tempFile := strings.ReplaceAll(string(content), "$ref", "__ref__")
	//replace environment variables in file
	tempFile = os.ExpandEnv(string(tempFile))
	tempFile = strings.ReplaceAll(string(tempFile), "__ref__", "$ref")
	content = []byte(tempFile)
	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromData(content)
	if err != nil {
		//e.Logger.Fatalf("error loading api specification '%s'", err)
	}

	//parse the main config
	var config *weoscontroller.GRPCAPIConfig
	if swagger.ExtensionProps.Extensions["x-weos-config"] != nil {

		data, err := swagger.ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
		if err != nil {
			//e.Logger.Fatalf("error loading api config '%s", err)
			return ctx
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			//e.Logger.Fatalf("error loading api config '%s", err)
			return ctx
		}

		err = api.AddConfig(config)
		if err != nil {
			//e.Logger.Fatalf("error setting up module '%s", err)
			return ctx
		}

		//TODO intialize the grpcmiddleware here

		err = api.Initialize()
		if err != nil {
			//e.Logger.Fatalf("error initializing application '%s'", err)
			return ctx
		}
	}
	return ctx
}
