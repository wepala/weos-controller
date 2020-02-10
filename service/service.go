package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface PluginInterface PluginLoaderInterface
type PathConfig struct {
	Mock       bool                `yaml:"mock"`
	Middleware []*MiddlewareConfig `yaml:"middleware"`
	Data       interface{}
}

type MiddlewareConfig struct {
	Plugin struct {
		FileName string           `yaml:"filename"`
		Config   *json.RawMessage `yaml:"config"`
	} `yaml:"plugin"`
	Handler  string                 `yaml:"handler"`
	Context  map[string]interface{} `yaml:"context"`
	Priority int                    `yaml:"priority"`
}

type controllerService struct {
	config       *openapi3.Swagger
	pluginLoader PluginLoaderInterface
}

func (s *controllerService) GetGlobalMiddlewareConfig() ([]*MiddlewareConfig, error) {
	if s.config.ExtensionProps.Extensions["x-weos-config"] != nil {
		globalConfigBytes, err := s.config.ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
		if err != nil {
			return nil, err
		}
		var globalConfig PathConfig
		err = json.Unmarshal(globalConfigBytes, &globalConfig)
		if err != nil {
			return nil, err
		}

		return globalConfig.Middleware, nil
	}
	return nil, nil
}

func (s *controllerService) GetPathConfig(path string, operation string) (*PathConfig, error) {
	weosConfig := s.config.Paths[path].GetOperation(strings.ToUpper(operation)).ExtensionProps.Extensions["x-weos-config"]
	if weosConfig == nil {
		return nil, nil
	}
	bytes, err := s.config.Paths[path].GetOperation(strings.ToUpper(operation)).ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
	if err != nil {
		return nil, err
	}
	pathConfig := PathConfig{}
	if err = json.Unmarshal(bytes, &pathConfig); err != nil {
		return nil, err
	}
	return &pathConfig, nil
}

func (s *controllerService) GetConfig() *openapi3.Swagger {
	return s.config
}

func (s *controllerService) GetHandlers(path string, config *PathConfig, pathInfo *openapi3.PathItem) ([]http.HandlerFunc, error) {
	globalHandlers, err := s.GetGlobalMiddlewareConfig()
	var middlewareConfig []*MiddlewareConfig

	if err != nil {
		log.Debug("there was an issue loading global handlers")
		return nil, err
	}

	if config != nil {
		middlewareConfig = config.Middleware
	}
	middlewareConfig = append(middlewareConfig, globalHandlers...)
	handlers := make([]http.HandlerFunc, len(middlewareConfig))

	//WEOS-168 if there are no handlers or the config has mock set to true return mock handlers
	if config == nil || len(handlers) == 0 || config.Mock {
		mockHandlers := make([]http.HandlerFunc, 1)
		for method, operation := range pathInfo.Operations() {
			var responseContent *openapi3.Content
			var statusCode int
			var err error

			for statusCodeString, responseRef := range operation.Responses {
				statusCode, err = strconv.Atoi(statusCodeString)
				if err != nil {
					log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the code statusCode %s could not be converted to an integer", path, method, statusCodeString)
				} else {
					responseContent = &responseRef.Value.Content
				}
			}

			if responseContent != nil {
				mh, err := NewMockExampleHandler(statusCode, responseContent)
				if err != nil {
					log.Errorf("could not mock the response for the path '%s' for the operation '%s' because the mock handler could not be created because '%s'", path, method, err)
				}
				mockHandlers[0] = mh.ServeHTTP
			}
		}
		return mockHandlers, nil
	} else { // otherwise let's load the plugins
		sort.Sort(NewMiddlewareConfigSorter(middlewareConfig))
		for key, mc := range middlewareConfig {
			log.Debugf("loading plugin %s", mc.Plugin.FileName)
			plugin, err := s.pluginLoader.GetPlugin(mc.Plugin.FileName)
			if err != nil {
				log.Errorf("error loading plugin %s", err)
				return nil, err
			}

			if mc.Plugin.Config != nil {
				err = plugin.AddConfig(*mc.Plugin.Config) //pass the raw json message that is loaded to the plugin
				if err != nil {
					log.Error("error loading plugin config", err)
					return nil, err
				}
			}

			log.Debugf("retrieving handler %s", mc.Handler)
			handlers[key] = plugin.GetHandlerByName(mc.Handler)
		}
	}

	return handlers, nil
}

var api openapi3.Swagger

type ServiceInterface interface {
	GetPathConfig(path string, operation string) (*PathConfig, error)
	GetConfig() *openapi3.Swagger
	GetHandlers(path string, config *PathConfig, pathInfo *openapi3.PathItem) ([]http.HandlerFunc, error)
	GetGlobalMiddlewareConfig() ([]*MiddlewareConfig, error)
}

func NewControllerService(apiConfig string, pluginLoader PluginLoaderInterface) (ServiceInterface, error) {

	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromFile(apiConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error loading %s: %s", apiConfig, err.Error()))
	}

	svc := &controllerService{
		config:       swagger,
		pluginLoader: pluginLoader,
	}

	return svc, nil
}
