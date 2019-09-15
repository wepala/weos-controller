package service

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface PluginInterface PluginLoaderInterface

type Config struct {
	ApiConfig *openapi3.Swagger
	Paths     Paths
}

type Paths map[string]PathItem
type PathItem map[string]*PathConfig

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
	config       *Config
	pluginLoader PluginLoaderInterface
}

func (s *controllerService) GetPathConfig(path string, operation string) (*PathConfig, error) {
	return s.config.Paths[path][operation], nil
}

func (s *controllerService) GetConfig() *Config {
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
	GetConfig() *Config
	GetHandlers(config *PathConfig) ([]http.HandlerFunc, error)
}

func NewControllerService(apiConfig string, controllerConfig string, pluginLoader PluginLoaderInterface) (ServiceInterface, error) {

	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromFile(apiConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error loading %s: %s", apiConfig, err.Error()))
	}

	config := &struct {
		Paths Paths
	}{
		Paths: Paths{},
	}

	//load controller config
	if controllerConfig != "" {
		log.Debugf("load config '%s'", controllerConfig)
		yamlFile, err := ioutil.ReadFile(controllerConfig)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error loading %s: %s", controllerConfig, err.Error()))
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			if strings.Contains(err.Error(), "MiddlewareConfig") {
				return nil, errors.New(err.Error())
			}
			return nil, err
		}
	}

	for pathName, path := range swagger.Paths {
		if config.Paths[pathName] == nil {
			config.Paths[pathName] = make(PathItem, 6)
		}

		if path.Get != nil && config.Paths[pathName]["get"] == nil {
			config.Paths[pathName]["get"] = &PathConfig{}
		}
	}

	svc := &controllerService{
		config: &Config{
			ApiConfig: swagger,
			Paths:     config.Paths,
		},
		pluginLoader: pluginLoader,
	}

	return svc, nil
}
