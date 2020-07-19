// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var (
	lockServiceInterfaceMockGetConfig                 sync.RWMutex
	lockServiceInterfaceMockGetGlobalMiddlewareConfig sync.RWMutex
	lockServiceInterfaceMockGetHandlers               sync.RWMutex
	lockServiceInterfaceMockGetPathConfig             sync.RWMutex
)

// Ensure, that ServiceInterfaceMock does implement service.ServiceInterface.
// If this is not the case, regenerate this file with moq.
var _ service.ServiceInterface = &ServiceInterfaceMock{}

// ServiceInterfaceMock is a mock implementation of service.ServiceInterface.
//
//     func TestSomethingThatUsesServiceInterface(t *testing.T) {
//
//         // make and configure a mocked service.ServiceInterface
//         mockedServiceInterface := &ServiceInterfaceMock{
//             GetConfigFunc: func() *openapi3.Swagger {
// 	               panic("mock out the GetConfig method")
//             },
//             GetGlobalMiddlewareConfigFunc: func() ([]*service.MiddlewareConfig, error) {
// 	               panic("mock out the GetGlobalMiddlewareConfig method")
//             },
//             GetHandlersFunc: func(config *service.PathConfig, mockHandler http.Handler) ([]http.HandlerFunc, error) {
// 	               panic("mock out the GetHandlers method")
//             },
//             GetPathConfigFunc: func(path string, operation string) (*service.PathConfig, error) {
// 	               panic("mock out the GetPathConfig method")
//             },
//         }
//
//         // use mockedServiceInterface in code that requires service.ServiceInterface
//         // and then make assertions.
//
//     }
type ServiceInterfaceMock struct {
	// GetConfigFunc mocks the GetConfig method.
	GetConfigFunc func() *openapi3.Swagger

	// GetGlobalMiddlewareConfigFunc mocks the GetGlobalMiddlewareConfig method.
	GetGlobalMiddlewareConfigFunc func() ([]*service.MiddlewareConfig, error)

	// GetHandlersFunc mocks the GetHandlers method.
	GetHandlersFunc func(config *service.PathConfig, mockHandler http.Handler) ([]http.HandlerFunc, error)

	// GetPathConfigFunc mocks the GetPathConfig method.
	GetPathConfigFunc func(path string, operation string) (*service.PathConfig, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetConfig holds details about calls to the GetConfig method.
		GetConfig []struct {
		}
		// GetGlobalMiddlewareConfig holds details about calls to the GetGlobalMiddlewareConfig method.
		GetGlobalMiddlewareConfig []struct {
		}
		// GetHandlers holds details about calls to the GetHandlers method.
		GetHandlers []struct {
			// Config is the config argument value.
			Config *service.PathConfig
			// MockHandler is the mockHandler argument value.
			MockHandler http.Handler
		}
		// GetPathConfig holds details about calls to the GetPathConfig method.
		GetPathConfig []struct {
			// Path is the path argument value.
			Path string
			// Operation is the operation argument value.
			Operation string
		}
	}
}

// GetConfig calls GetConfigFunc.
func (mock *ServiceInterfaceMock) GetConfig() *openapi3.Swagger {
	if mock.GetConfigFunc == nil {
		panic("ServiceInterfaceMock.GetConfigFunc: method is nil but ServiceInterface.GetConfig was just called")
	}
	callInfo := struct {
	}{}
	lockServiceInterfaceMockGetConfig.Lock()
	mock.calls.GetConfig = append(mock.calls.GetConfig, callInfo)
	lockServiceInterfaceMockGetConfig.Unlock()
	return mock.GetConfigFunc()
}

// GetConfigCalls gets all the calls that were made to GetConfig.
// Check the length with:
//     len(mockedServiceInterface.GetConfigCalls())
func (mock *ServiceInterfaceMock) GetConfigCalls() []struct {
} {
	var calls []struct {
	}
	lockServiceInterfaceMockGetConfig.RLock()
	calls = mock.calls.GetConfig
	lockServiceInterfaceMockGetConfig.RUnlock()
	return calls
}

// GetGlobalMiddlewareConfig calls GetGlobalMiddlewareConfigFunc.
func (mock *ServiceInterfaceMock) GetGlobalMiddlewareConfig() ([]*service.MiddlewareConfig, error) {
	if mock.GetGlobalMiddlewareConfigFunc == nil {
		panic("ServiceInterfaceMock.GetGlobalMiddlewareConfigFunc: method is nil but ServiceInterface.GetGlobalMiddlewareConfig was just called")
	}
	callInfo := struct {
	}{}
	lockServiceInterfaceMockGetGlobalMiddlewareConfig.Lock()
	mock.calls.GetGlobalMiddlewareConfig = append(mock.calls.GetGlobalMiddlewareConfig, callInfo)
	lockServiceInterfaceMockGetGlobalMiddlewareConfig.Unlock()
	return mock.GetGlobalMiddlewareConfigFunc()
}

// GetGlobalMiddlewareConfigCalls gets all the calls that were made to GetGlobalMiddlewareConfig.
// Check the length with:
//     len(mockedServiceInterface.GetGlobalMiddlewareConfigCalls())
func (mock *ServiceInterfaceMock) GetGlobalMiddlewareConfigCalls() []struct {
} {
	var calls []struct {
	}
	lockServiceInterfaceMockGetGlobalMiddlewareConfig.RLock()
	calls = mock.calls.GetGlobalMiddlewareConfig
	lockServiceInterfaceMockGetGlobalMiddlewareConfig.RUnlock()
	return calls
}

// GetHandlers calls GetHandlersFunc.
func (mock *ServiceInterfaceMock) GetHandlers(config *service.PathConfig, mockHandler http.Handler) ([]http.HandlerFunc, error) {
	if mock.GetHandlersFunc == nil {
		panic("ServiceInterfaceMock.GetHandlersFunc: method is nil but ServiceInterface.GetHandlers was just called")
	}
	callInfo := struct {
		Config      *service.PathConfig
		MockHandler http.Handler
	}{
		Config:      config,
		MockHandler: mockHandler,
	}
	lockServiceInterfaceMockGetHandlers.Lock()
	mock.calls.GetHandlers = append(mock.calls.GetHandlers, callInfo)
	lockServiceInterfaceMockGetHandlers.Unlock()
	return mock.GetHandlersFunc(config, mockHandler)
}

// GetHandlersCalls gets all the calls that were made to GetHandlers.
// Check the length with:
//     len(mockedServiceInterface.GetHandlersCalls())
func (mock *ServiceInterfaceMock) GetHandlersCalls() []struct {
	Config      *service.PathConfig
	MockHandler http.Handler
} {
	var calls []struct {
		Config      *service.PathConfig
		MockHandler http.Handler
	}
	lockServiceInterfaceMockGetHandlers.RLock()
	calls = mock.calls.GetHandlers
	lockServiceInterfaceMockGetHandlers.RUnlock()
	return calls
}

// GetPathConfig calls GetPathConfigFunc.
func (mock *ServiceInterfaceMock) GetPathConfig(path string, operation string) (*service.PathConfig, error) {
	if mock.GetPathConfigFunc == nil {
		panic("ServiceInterfaceMock.GetPathConfigFunc: method is nil but ServiceInterface.GetPathConfig was just called")
	}
	callInfo := struct {
		Path      string
		Operation string
	}{
		Path:      path,
		Operation: operation,
	}
	lockServiceInterfaceMockGetPathConfig.Lock()
	mock.calls.GetPathConfig = append(mock.calls.GetPathConfig, callInfo)
	lockServiceInterfaceMockGetPathConfig.Unlock()
	return mock.GetPathConfigFunc(path, operation)
}

// GetPathConfigCalls gets all the calls that were made to GetPathConfig.
// Check the length with:
//     len(mockedServiceInterface.GetPathConfigCalls())
func (mock *ServiceInterfaceMock) GetPathConfigCalls() []struct {
	Path      string
	Operation string
} {
	var calls []struct {
		Path      string
		Operation string
	}
	lockServiceInterfaceMockGetPathConfig.RLock()
	calls = mock.calls.GetPathConfig
	lockServiceInterfaceMockGetPathConfig.RUnlock()
	return calls
}

var (
	lockPluginInterfaceMockAddConfig        sync.RWMutex
	lockPluginInterfaceMockAddLogger        sync.RWMutex
	lockPluginInterfaceMockAddPathConfig    sync.RWMutex
	lockPluginInterfaceMockAddSession       sync.RWMutex
	lockPluginInterfaceMockGetHandlerByName sync.RWMutex
)

// Ensure, that PluginInterfaceMock does implement service.PluginInterface.
// If this is not the case, regenerate this file with moq.
var _ service.PluginInterface = &PluginInterfaceMock{}

// PluginInterfaceMock is a mock implementation of service.PluginInterface.
//
//     func TestSomethingThatUsesPluginInterface(t *testing.T) {
//
//         // make and configure a mocked service.PluginInterface
//         mockedPluginInterface := &PluginInterfaceMock{
//             AddConfigFunc: func(config json.RawMessage) error {
// 	               panic("mock out the AddConfig method")
//             },
//             AddLoggerFunc: func(logger logrus.Ext1FieldLogger)  {
// 	               panic("mock out the AddLogger method")
//             },
//             AddPathConfigFunc: func(handler string, config json.RawMessage) error {
// 	               panic("mock out the AddPathConfig method")
//             },
//             AddSessionFunc: func(session sessions.Store)  {
// 	               panic("mock out the AddSession method")
//             },
//             GetHandlerByNameFunc: func(name string) http.HandlerFunc {
// 	               panic("mock out the GetHandlerByName method")
//             },
//         }
//
//         // use mockedPluginInterface in code that requires service.PluginInterface
//         // and then make assertions.
//
//     }
type PluginInterfaceMock struct {
	// AddConfigFunc mocks the AddConfig method.
	AddConfigFunc func(config json.RawMessage) error

	// AddLoggerFunc mocks the AddLogger method.
	AddLoggerFunc func(logger logrus.Ext1FieldLogger)

	// AddPathConfigFunc mocks the AddPathConfig method.
	AddPathConfigFunc func(handler string, config json.RawMessage) error

	// AddSessionFunc mocks the AddSession method.
	AddSessionFunc func(session sessions.Store)

	// GetHandlerByNameFunc mocks the GetHandlerByName method.
	GetHandlerByNameFunc func(name string) http.HandlerFunc

	// calls tracks calls to the methods.
	calls struct {
		// AddConfig holds details about calls to the AddConfig method.
		AddConfig []struct {
			// Config is the config argument value.
			Config json.RawMessage
		}
		// AddLogger holds details about calls to the AddLogger method.
		AddLogger []struct {
			// Logger is the logger argument value.
			Logger logrus.Ext1FieldLogger
		}
		// AddPathConfig holds details about calls to the AddPathConfig method.
		AddPathConfig []struct {
			// Handler is the handler argument value.
			Handler string
			// Config is the config argument value.
			Config json.RawMessage
		}
		// AddSession holds details about calls to the AddSession method.
		AddSession []struct {
			// Session is the session argument value.
			Session sessions.Store
		}
		// GetHandlerByName holds details about calls to the GetHandlerByName method.
		GetHandlerByName []struct {
			// Name is the name argument value.
			Name string
		}
	}
}

// AddConfig calls AddConfigFunc.
func (mock *PluginInterfaceMock) AddConfig(config json.RawMessage) error {
	if mock.AddConfigFunc == nil {
		panic("PluginInterfaceMock.AddConfigFunc: method is nil but PluginInterface.AddConfig was just called")
	}
	callInfo := struct {
		Config json.RawMessage
	}{
		Config: config,
	}
	lockPluginInterfaceMockAddConfig.Lock()
	mock.calls.AddConfig = append(mock.calls.AddConfig, callInfo)
	lockPluginInterfaceMockAddConfig.Unlock()
	return mock.AddConfigFunc(config)
}

// AddConfigCalls gets all the calls that were made to AddConfig.
// Check the length with:
//     len(mockedPluginInterface.AddConfigCalls())
func (mock *PluginInterfaceMock) AddConfigCalls() []struct {
	Config json.RawMessage
} {
	var calls []struct {
		Config json.RawMessage
	}
	lockPluginInterfaceMockAddConfig.RLock()
	calls = mock.calls.AddConfig
	lockPluginInterfaceMockAddConfig.RUnlock()
	return calls
}

// AddLogger calls AddLoggerFunc.
func (mock *PluginInterfaceMock) AddLogger(logger logrus.Ext1FieldLogger) {
	if mock.AddLoggerFunc == nil {
		panic("PluginInterfaceMock.AddLoggerFunc: method is nil but PluginInterface.AddLogger was just called")
	}
	callInfo := struct {
		Logger logrus.Ext1FieldLogger
	}{
		Logger: logger,
	}
	lockPluginInterfaceMockAddLogger.Lock()
	mock.calls.AddLogger = append(mock.calls.AddLogger, callInfo)
	lockPluginInterfaceMockAddLogger.Unlock()
	mock.AddLoggerFunc(logger)
}

// AddLoggerCalls gets all the calls that were made to AddLogger.
// Check the length with:
//     len(mockedPluginInterface.AddLoggerCalls())
func (mock *PluginInterfaceMock) AddLoggerCalls() []struct {
	Logger logrus.Ext1FieldLogger
} {
	var calls []struct {
		Logger logrus.Ext1FieldLogger
	}
	lockPluginInterfaceMockAddLogger.RLock()
	calls = mock.calls.AddLogger
	lockPluginInterfaceMockAddLogger.RUnlock()
	return calls
}

// AddPathConfig calls AddPathConfigFunc.
func (mock *PluginInterfaceMock) AddPathConfig(handler string, config json.RawMessage) error {
	if mock.AddPathConfigFunc == nil {
		panic("PluginInterfaceMock.AddPathConfigFunc: method is nil but PluginInterface.AddPathConfig was just called")
	}
	callInfo := struct {
		Handler string
		Config  json.RawMessage
	}{
		Handler: handler,
		Config:  config,
	}
	lockPluginInterfaceMockAddPathConfig.Lock()
	mock.calls.AddPathConfig = append(mock.calls.AddPathConfig, callInfo)
	lockPluginInterfaceMockAddPathConfig.Unlock()
	return mock.AddPathConfigFunc(handler, config)
}

// AddPathConfigCalls gets all the calls that were made to AddPathConfig.
// Check the length with:
//     len(mockedPluginInterface.AddPathConfigCalls())
func (mock *PluginInterfaceMock) AddPathConfigCalls() []struct {
	Handler string
	Config  json.RawMessage
} {
	var calls []struct {
		Handler string
		Config  json.RawMessage
	}
	lockPluginInterfaceMockAddPathConfig.RLock()
	calls = mock.calls.AddPathConfig
	lockPluginInterfaceMockAddPathConfig.RUnlock()
	return calls
}

// AddSession calls AddSessionFunc.
func (mock *PluginInterfaceMock) AddSession(session sessions.Store) {
	if mock.AddSessionFunc == nil {
		panic("PluginInterfaceMock.AddSessionFunc: method is nil but PluginInterface.AddSession was just called")
	}
	callInfo := struct {
		Session sessions.Store
	}{
		Session: session,
	}
	lockPluginInterfaceMockAddSession.Lock()
	mock.calls.AddSession = append(mock.calls.AddSession, callInfo)
	lockPluginInterfaceMockAddSession.Unlock()
	mock.AddSessionFunc(session)
}

// AddSessionCalls gets all the calls that were made to AddSession.
// Check the length with:
//     len(mockedPluginInterface.AddSessionCalls())
func (mock *PluginInterfaceMock) AddSessionCalls() []struct {
	Session sessions.Store
} {
	var calls []struct {
		Session sessions.Store
	}
	lockPluginInterfaceMockAddSession.RLock()
	calls = mock.calls.AddSession
	lockPluginInterfaceMockAddSession.RUnlock()
	return calls
}

// GetHandlerByName calls GetHandlerByNameFunc.
func (mock *PluginInterfaceMock) GetHandlerByName(name string) http.HandlerFunc {
	if mock.GetHandlerByNameFunc == nil {
		panic("PluginInterfaceMock.GetHandlerByNameFunc: method is nil but PluginInterface.GetHandlerByName was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	lockPluginInterfaceMockGetHandlerByName.Lock()
	mock.calls.GetHandlerByName = append(mock.calls.GetHandlerByName, callInfo)
	lockPluginInterfaceMockGetHandlerByName.Unlock()
	return mock.GetHandlerByNameFunc(name)
}

// GetHandlerByNameCalls gets all the calls that were made to GetHandlerByName.
// Check the length with:
//     len(mockedPluginInterface.GetHandlerByNameCalls())
func (mock *PluginInterfaceMock) GetHandlerByNameCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	lockPluginInterfaceMockGetHandlerByName.RLock()
	calls = mock.calls.GetHandlerByName
	lockPluginInterfaceMockGetHandlerByName.RUnlock()
	return calls
}

var (
	lockPluginLoaderInterfaceMockGetPlugin     sync.RWMutex
	lockPluginLoaderInterfaceMockGetRepository sync.RWMutex
)

// Ensure, that PluginLoaderInterfaceMock does implement service.PluginLoaderInterface.
// If this is not the case, regenerate this file with moq.
var _ service.PluginLoaderInterface = &PluginLoaderInterfaceMock{}

// PluginLoaderInterfaceMock is a mock implementation of service.PluginLoaderInterface.
//
//     func TestSomethingThatUsesPluginLoaderInterface(t *testing.T) {
//
//         // make and configure a mocked service.PluginLoaderInterface
//         mockedPluginLoaderInterface := &PluginLoaderInterfaceMock{
//             GetPluginFunc: func(fileName string) (service.PluginInterface, error) {
// 	               panic("mock out the GetPlugin method")
//             },
//             GetRepositoryFunc: func(fileName string) (service.RepositoryInterface, error) {
// 	               panic("mock out the GetRepository method")
//             },
//         }
//
//         // use mockedPluginLoaderInterface in code that requires service.PluginLoaderInterface
//         // and then make assertions.
//
//     }
type PluginLoaderInterfaceMock struct {
	// GetPluginFunc mocks the GetPlugin method.
	GetPluginFunc func(fileName string) (service.PluginInterface, error)

	// GetRepositoryFunc mocks the GetRepository method.
	GetRepositoryFunc func(fileName string) (service.RepositoryInterface, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetPlugin holds details about calls to the GetPlugin method.
		GetPlugin []struct {
			// FileName is the fileName argument value.
			FileName string
		}
		// GetRepository holds details about calls to the GetRepository method.
		GetRepository []struct {
			// FileName is the fileName argument value.
			FileName string
		}
	}
}

// GetPlugin calls GetPluginFunc.
func (mock *PluginLoaderInterfaceMock) GetPlugin(fileName string) (service.PluginInterface, error) {
	if mock.GetPluginFunc == nil {
		panic("PluginLoaderInterfaceMock.GetPluginFunc: method is nil but PluginLoaderInterface.GetPlugin was just called")
	}
	callInfo := struct {
		FileName string
	}{
		FileName: fileName,
	}
	lockPluginLoaderInterfaceMockGetPlugin.Lock()
	mock.calls.GetPlugin = append(mock.calls.GetPlugin, callInfo)
	lockPluginLoaderInterfaceMockGetPlugin.Unlock()
	return mock.GetPluginFunc(fileName)
}

// GetPluginCalls gets all the calls that were made to GetPlugin.
// Check the length with:
//     len(mockedPluginLoaderInterface.GetPluginCalls())
func (mock *PluginLoaderInterfaceMock) GetPluginCalls() []struct {
	FileName string
} {
	var calls []struct {
		FileName string
	}
	lockPluginLoaderInterfaceMockGetPlugin.RLock()
	calls = mock.calls.GetPlugin
	lockPluginLoaderInterfaceMockGetPlugin.RUnlock()
	return calls
}

// GetRepository calls GetRepositoryFunc.
func (mock *PluginLoaderInterfaceMock) GetRepository(fileName string) (service.RepositoryInterface, error) {
	if mock.GetRepositoryFunc == nil {
		panic("PluginLoaderInterfaceMock.GetRepositoryFunc: method is nil but PluginLoaderInterface.GetRepository was just called")
	}
	callInfo := struct {
		FileName string
	}{
		FileName: fileName,
	}
	lockPluginLoaderInterfaceMockGetRepository.Lock()
	mock.calls.GetRepository = append(mock.calls.GetRepository, callInfo)
	lockPluginLoaderInterfaceMockGetRepository.Unlock()
	return mock.GetRepositoryFunc(fileName)
}

// GetRepositoryCalls gets all the calls that were made to GetRepository.
// Check the length with:
//     len(mockedPluginLoaderInterface.GetRepositoryCalls())
func (mock *PluginLoaderInterfaceMock) GetRepositoryCalls() []struct {
	FileName string
} {
	var calls []struct {
		FileName string
	}
	lockPluginLoaderInterfaceMockGetRepository.RLock()
	calls = mock.calls.GetRepository
	lockPluginLoaderInterfaceMockGetRepository.RUnlock()
	return calls
}
