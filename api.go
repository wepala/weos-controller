//go:generate moq -out mocks_test.go -pkg weoscontroller . APIInterface
package weoscontroller

import (
	"github.com/labstack/echo/v4"
)

//define an interface that all plugins must implement
type APIInterface interface {
	AddPathConfig(path string, config *PathConfig) error
	AddConfig(config *APIConfig) error
	Initialize() error
	EchoInstance() *echo.Echo
	SetEchoInstance(e *echo.Echo)
}

type GRPCAPIInterface interface {
	AddPathConfig(path string, config *PathConfig) error
	AddConfig(config *GRPCAPIConfig) error
	Initialize() error
}
