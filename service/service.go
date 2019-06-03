package service

import "github.com/getkin/kin-openapi/openapi3"

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface

type Config struct {
	ApiConfig   openapi3.Swagger
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
}

func (*controllerService) GetPathConfig(path string) (map[string]*PathConfig, error) {
	panic("implement me")
}

func (*controllerService) GetConfig() (*Config, error) {
	panic("implement me")
}

var api openapi3.Swagger

type ServiceInterface interface {
	GetPathConfig(path string) (map[string]*PathConfig, error)
	GetConfig() (*Config, error)
}

func NewControllerService() ServiceInterface {
	return &controllerService{}
}
