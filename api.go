//go:generate moq -out mocks_test.go -pkg weoscontroller . APIInterface
package weoscontroller

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	AddConfig(config *APIConfig) error
	Initialize() error
	Context() context.Context
	SetContext(c context.Context)
	SetAllMiddleware()
	GetStreamMiddleware() grpc.ServerOption
	GetUnaryMiddleware() grpc.ServerOption
}
