package weoscontroller

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

func Initialize(e *echo.Echo, api APIInterface, apiConfig string) *echo.Echo {
	e.HideBanner = true
	if apiConfig == "" {
		apiConfig = "./api.yaml"
	}

	//set echo instance because the instance may not already be in the api that is passed in but the handlers must have access to it
	api.SetEchoInstance(e)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		CustomErrorHandler(err, c)
		e.DefaultHTTPErrorHandler(err, c)
	}

	var content []byte
	var err error
	//try load file if it's a yaml file otherwise it's the contents of a yaml file WEOS-1009
	if strings.Contains(apiConfig, ".yaml") || strings.Contains(apiConfig, "/yml") {
		content, err = ioutil.ReadFile(apiConfig)
		if err != nil {
			e.Logger.Fatalf("error loading api specification '%s'", err)
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
		e.Logger.Fatalf("error loading api specification '%s'", err)
	}

	//parse the main config
	var config *APIConfig
	if swagger.ExtensionProps.Extensions["x-weos-config"] != nil {

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
		//setup middleware  - https://echo.labstack.com/middleware/

		//setup global pre middleware
		var preMiddlewares []echo.MiddlewareFunc
		for _, middlewareName := range config.PreMiddleware {
			t := reflect.ValueOf(api)
			m := t.MethodByName(middlewareName)
			if !m.IsValid() {
				e.Logger.Fatalf("invalid handler set '%s'", middlewareName)
			}
			preMiddlewares = append(preMiddlewares, m.Interface().(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc))
		}
		//all routes setup after this will use this middleware
		e.Pre(preMiddlewares...)

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
		//all routes setup after this will use this middleware
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

					if weosConfig.Group { //TODO move this form here because it creates weird behavior
						group := e.Group(config.BasePath + path)
						err = api.AddPathConfig(config.BasePath+path, weosConfig)
						if err != nil {
							e.Logger.Fatalf("error adding path config '%s' '%s'", config.BasePath+path, err)
						}
						group.Use(middlewares...)
					} else {
						//TODO make it so that it automatically matches the paths to a group based on the prefix
						//update path so that the open api way of specifying url parameters is change to the echo style of url parameters
						re := regexp.MustCompile(`\{([a-zA-Z0-9\-_]+?)\}`)
						echoPath := re.ReplaceAllString(path, `:$1`)
						err = api.AddPathConfig(config.BasePath+echoPath, weosConfig)
						if err != nil {
							e.Logger.Fatalf("error adding path config '%s' '%s'", echoPath, err)
						}
						switch method {
						case "GET":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.GET(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "POST":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.POST(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "PUT":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.PUT(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "PATCH":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.PATCH(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "DELETE":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.DELETE(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "HEAD":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.HEAD(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "OPTIONS":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.OPTIONS(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "TRACE":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.TRACE(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)
						case "CONNECT":
							if !weosConfig.DisableCors {
								e.Use(EnableCORS(method))
							}
							e.CONNECT(config.BasePath+echoPath, handler.Interface().(func(ctx echo.Context) error), middlewares...)

						}
					}

				}
			}
		}

	}
	return e
}
