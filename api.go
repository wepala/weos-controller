//go:generate moq -out mocks_test.go -pkg weoscontroller_test . APIInterface
package weoscontroller

import "github.com/labstack/echo/v4"

//define an interface that all plugins must implement
type APIInterface interface {
	AddConfig(config *APIConfig) error
	Initialize() error
	EchoInstance() *echo.Echo
	SetEchoInstance(e *echo.Echo)
}
