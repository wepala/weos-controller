package service

import "github.com/getkin/kin-openapi/openapi3"

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface

type Config struct {
	ApiConfig   *openapi3.Swagger
	PathConfigs map[string]map[string]PathConfig
}

type PathConfig struct {
	Templates  []string
	Middleware []*MiddlewareConfig
}

type MiddlewareConfig struct {
	FileName    string
	ServerName  string
	HandlerName string
}

type controllerService struct {
	config *Config
}

func (*controllerService) GetPathConfig(path string, operation string) (*PathConfig, error) {
	return nil, nil
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

	svc := &controllerService{
		config: &Config{
			ApiConfig:   swagger,
			PathConfigs: map[string]map[string]PathConfig{},
		},
	}

	//TODO loop through the paths and build the path configs

	return svc, nil
}
