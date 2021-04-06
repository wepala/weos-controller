package weoscontroller

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func Configure(e *echo.Echo, apiConfigPath string, api APIInterface) *echo.Echo {
	if apiConfigPath == "" {
		apiConfigPath = "./api.yaml"
	}
	content, err := ioutil.ReadFile(apiConfigPath)
	if err != nil {
		e.Logger.Fatalf("error loading api specification '%s'", err)
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
		e.Logger.Fatalf("error loading api specification '%s'", err)
	}

	//parse the main config
	if swagger.ExtensionProps.Extensions["x-weos-config"] != nil {
		var config *APIConfig
		data, err := swagger.ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
		if err != nil {
			e.Logger.Fatalf("error loading api config '%s", err)
			return e
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			e.Logger.Fatalf("error loading api config '%s", err)
			return e
		}

		err = api.AddConfig(config)
		if err != nil {
			e.Logger.Fatalf("error setting up module '%s", err)
			return e
		}

		//setup global middleware
		var middlewares []echo.MiddlewareFunc
		for _, middlewareName := range config.Middleware {
			t := reflect.ValueOf(api)
			m := t.MethodByName(middlewareName)
			if !m.IsValid() {
				e.Logger.Fatalf("invalid handler set '%s'", middlewareName)
			}
			middlewares = append(middlewares, m.Interface().(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc))
		}
		e.Use(middlewares...)
		err = api.Initialize()
		if err != nil {
			e.Logger.Fatalf("error initializing application '%s'", err)
			return e
		}
	}

	knownActions := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE", "CONNECT"}

	for path, pathData := range swagger.Paths {
		for _, method := range knownActions {
			var weosConfig *PathConfig
			//get the handler using reflection. This should be fine because this is only on startup
			if pathData.GetOperation(strings.ToUpper(method)) != nil {
				weosConfigData := pathData.GetOperation(strings.ToUpper(method)).ExtensionProps.Extensions["x-weos-config"]
				if weosConfigData != nil {
					bytes, err := weosConfigData.(json.RawMessage).MarshalJSON()
					if err != nil {
						e.Logger.Fatalf("error reading weos config on the path '%s' '%s'", path, err)
					}

					if err = json.Unmarshal(bytes, &weosConfig); err != nil {
						e.Logger.Fatalf("error reading weos config on the path '%s' '%s'", path, err)
						return e
					}

					t := reflect.ValueOf(api)
					handler := t.MethodByName(weosConfig.Handler)
					//only show error if handler was set
					if weosConfig.Handler != "" && !handler.IsValid() {
						e.Logger.Fatalf("invalid handler set '%s'", weosConfig.Handler)
					}

					var middlewares []echo.MiddlewareFunc
					for _, middlewareName := range weosConfig.Middleware {
						m := t.MethodByName(middlewareName)
						if !m.IsValid() {
							e.Logger.Fatalf("invalid handler set '%s'", middlewareName)
						}
						middlewares = append(middlewares, m.Interface().(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc))
					}

					if weosConfig.Group {
						group := e.Group(path)
						group.Use(middlewares...)
					} else {
						//TODO make it so that it automatically matches the paths to a group based on the prefix
						switch method {
						case "GET":
							e.GET(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "POST":
							e.POST(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "PUT":
							e.PUT(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "PATCH":
							e.PATCH(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "DELETE":
							e.DELETE(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "HEAD":
							e.HEAD(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "OPTIONS":
							e.OPTIONS(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "TRACE":
							e.TRACE(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "CONNECT":
							e.CONNECT(path, handler.Interface().(func(ctx echo.Context) error), middlewares...)

						}
					}

				}
			}
		}

	}
	return e
}
