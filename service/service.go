package service

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	Middleware []*MiddlewareConfig
	Data       interface{}
}

type MiddlewareConfig struct {
	FileName    string
	ServerName  string
	HandlerName string
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
	}{}

	//load controller config
	if controllerConfig != "" {
		log.Debugf("load config '%s'", controllerConfig)
		yamlFile, err := ioutil.ReadFile(controllerConfig)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return nil, err
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
