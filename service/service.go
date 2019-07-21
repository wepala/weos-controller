package service

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface

type Config struct {
	ApiConfig *openapi3.Swagger
	Paths     Paths
}

type Paths map[string]PathItem
type PathItem map[string]*PathConfig

type PathConfig struct {
	Templates  []string
	Middleware []*MiddlewareConfig `yaml:"middleware"`
	Data       interface{}
}

func (config *PathConfig) getHandlers() []*http.HandlerFunc {
	handlers := make([]*http.HandlerFunc, len(config.Middleware))
	for _, mc := range config.Middleware {
		plugin, _ := GetPlugin(mc.File)
		handlers = append(handlers, plugin.GetHandlerByName(mc.Handler))
	}
	return handlers
}

type MiddlewareConfig struct {
	File    string `yaml:"file"`
	Handler string `yaml:"handler"`
}

type controllerService struct {
	config *Config
}

func (s *controllerService) GetPathConfig(path string, operation string) (*PathConfig, error) {
	return s.config.Paths[path][operation], nil
}

func (s *controllerService) GetConfig() *Config {
	return s.config
}

var api openapi3.Swagger

type ServiceInterface interface {
	GetPathConfig(path string, operation string) (*PathConfig, error)
	GetConfig() *Config
}

func NewControllerService(apiConfig string, controllerConfig string) (ServiceInterface, error) {

	loader := openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromFile(apiConfig)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			if strings.Contains(err.Error(), "MiddlewareConfig") {
				return nil, errors.New("the list of middlewares must be an array in the config")
			}

			if strings.Contains(err.Error(), "Template") {
				return nil, errors.New("the list of templates must be an array in the config")
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
	}

	return svc, nil
}
