package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface PluginInterface PluginLoaderInterface
type PathConfig struct {
	Middleware []*MiddlewareConfig `yaml:"middleware"`
	Data       interface{}
}

type MiddlewareConfig struct {
	Plugin struct {
		FileName string                 `yaml:"filename"`
		Config   map[string]interface{} `yaml:"config"`
	} `yaml:"plugin"`
	Handler string                 `yaml:"handler"`
	Context map[string]interface{} `yaml:"context"`
}

type controllerService struct {
	config       *openapi3.Swagger
	pluginLoader PluginLoaderInterface
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

func (s *controllerService) GetHandlers(config *PathConfig) ([]http.HandlerFunc, error) {
	if config == nil {
		return nil, errors.New("path config cannot be empty")
	}
	handlers := make([]http.HandlerFunc, len(config.Middleware))
	for key, mc := range config.Middleware {
		log.Debugf("loading plugin %s", mc.Plugin.FileName)
		plugin, err := s.pluginLoader.GetPlugin(mc.Plugin.FileName)
		if err != nil {
			log.Errorf("error loading plugin %s", err)
			return nil, err
		}
		log.Debugf("retrieving handler %s", mc.Handler)
		handlers[key] = plugin.GetHandlerByName(mc.Handler)
	}
	return handlers, nil
}

var api openapi3.Swagger

type ServiceInterface interface {
	GetPathConfig(path string, operation string) (*PathConfig, error)
	GetConfig() *openapi3.Swagger
	GetHandlers(config *PathConfig) ([]http.HandlerFunc, error)
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
