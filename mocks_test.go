// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package weoscontroller_test

import (
	"github.com/labstack/echo/v4"
	"github.com/wepala/weos-controller"
	"sync"
)

// Ensure, that APIInterfaceMock does implement APIInterface.
// If this is not the case, regenerate this file with moq.
var _ weoscontroller.APIInterface = &APIInterfaceMock{}

// APIInterfaceMock is a mock implementation of APIInterface.
//
// 	func TestSomethingThatUsesAPIInterface(t *testing.T) {
//
// 		// make and configure a mocked APIInterface
// 		mockedAPIInterface := &APIInterfaceMock{
// 			AddConfigFunc: func(config *weoscontroller.APIConfig) error {
// 				panic("mock out the AddConfig method")
// 			},
// 			AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
// 				panic("mock out the AddPathConfig method")
// 			},
// 			EchoInstanceFunc: func() *echo.Echo {
// 				panic("mock out the EchoInstance method")
// 			},
// 			InitializeFunc: func() error {
// 				panic("mock out the Initialize method")
// 			},
// 			SetEchoInstanceFunc: func(e *echo.Echo)  {
// 				panic("mock out the SetEchoInstance method")
// 			},
// 		}
//
// 		// use mockedAPIInterface in code that requires APIInterface
// 		// and then make assertions.
//
// 	}
type APIInterfaceMock struct {
	// AddConfigFunc mocks the AddConfig method.
	AddConfigFunc func(config *weoscontroller.APIConfig) error

	// AddPathConfigFunc mocks the AddPathConfig method.
	AddPathConfigFunc func(path string, config *weoscontroller.PathConfig) error

	// EchoInstanceFunc mocks the EchoInstance method.
	EchoInstanceFunc func() *echo.Echo

	// InitializeFunc mocks the Initialize method.
	InitializeFunc func() error

	// SetEchoInstanceFunc mocks the SetEchoInstance method.
	SetEchoInstanceFunc func(e *echo.Echo)

	// calls tracks calls to the methods.
	calls struct {
		// AddConfig holds details about calls to the AddConfig method.
		AddConfig []struct {
			// Config is the config argument value.
			Config *weoscontroller.APIConfig
		}
		// AddPathConfig holds details about calls to the AddPathConfig method.
		AddPathConfig []struct {
			// Path is the path argument value.
			Path string
			// Config is the config argument value.
			Config *weoscontroller.PathConfig
		}
		// EchoInstance holds details about calls to the EchoInstance method.
		EchoInstance []struct {
		}
		// Initialize holds details about calls to the Initialize method.
		Initialize []struct {
		}
		// SetEchoInstance holds details about calls to the SetEchoInstance method.
		SetEchoInstance []struct {
			// E is the e argument value.
			E *echo.Echo
		}
	}
	lockAddConfig       sync.RWMutex
	lockAddPathConfig   sync.RWMutex
	lockEchoInstance    sync.RWMutex
	lockInitialize      sync.RWMutex
	lockSetEchoInstance sync.RWMutex
}

// AddConfig calls AddConfigFunc.
func (mock *APIInterfaceMock) AddConfig(config *weoscontroller.APIConfig) error {
	if mock.AddConfigFunc == nil {
		panic("APIInterfaceMock.AddConfigFunc: method is nil but APIInterface.AddConfig was just called")
	}
	callInfo := struct {
		Config *weoscontroller.APIConfig
	}{
		Config: config,
	}
	mock.lockAddConfig.Lock()
	mock.calls.AddConfig = append(mock.calls.AddConfig, callInfo)
	mock.lockAddConfig.Unlock()
	return mock.AddConfigFunc(config)
}

// AddConfigCalls gets all the calls that were made to AddConfig.
// Check the length with:
//     len(mockedAPIInterface.AddConfigCalls())
func (mock *APIInterfaceMock) AddConfigCalls() []struct {
	Config *weoscontroller.APIConfig
} {
	var calls []struct {
		Config *weoscontroller.APIConfig
	}
	mock.lockAddConfig.RLock()
	calls = mock.calls.AddConfig
	mock.lockAddConfig.RUnlock()
	return calls
}

// AddPathConfig calls AddPathConfigFunc.
func (mock *APIInterfaceMock) AddPathConfig(path string, config *weoscontroller.PathConfig) error {
	if mock.AddPathConfigFunc == nil {
		panic("APIInterfaceMock.AddPathConfigFunc: method is nil but APIInterface.AddPathConfig was just called")
	}
	callInfo := struct {
		Path   string
		Config *weoscontroller.PathConfig
	}{
		Path:   path,
		Config: config,
	}
	mock.lockAddPathConfig.Lock()
	mock.calls.AddPathConfig = append(mock.calls.AddPathConfig, callInfo)
	mock.lockAddPathConfig.Unlock()
	return mock.AddPathConfigFunc(path, config)
}

// AddPathConfigCalls gets all the calls that were made to AddPathConfig.
// Check the length with:
//     len(mockedAPIInterface.AddPathConfigCalls())
func (mock *APIInterfaceMock) AddPathConfigCalls() []struct {
	Path   string
	Config *weoscontroller.PathConfig
} {
	var calls []struct {
		Path   string
		Config *weoscontroller.PathConfig
	}
	mock.lockAddPathConfig.RLock()
	calls = mock.calls.AddPathConfig
	mock.lockAddPathConfig.RUnlock()
	return calls
}

// EchoInstance calls EchoInstanceFunc.
func (mock *APIInterfaceMock) EchoInstance() *echo.Echo {
	if mock.EchoInstanceFunc == nil {
		panic("APIInterfaceMock.EchoInstanceFunc: method is nil but APIInterface.EchoInstance was just called")
	}
	callInfo := struct {
	}{}
	mock.lockEchoInstance.Lock()
	mock.calls.EchoInstance = append(mock.calls.EchoInstance, callInfo)
	mock.lockEchoInstance.Unlock()
	return mock.EchoInstanceFunc()
}

// EchoInstanceCalls gets all the calls that were made to EchoInstance.
// Check the length with:
//     len(mockedAPIInterface.EchoInstanceCalls())
func (mock *APIInterfaceMock) EchoInstanceCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockEchoInstance.RLock()
	calls = mock.calls.EchoInstance
	mock.lockEchoInstance.RUnlock()
	return calls
}

// Initialize calls InitializeFunc.
func (mock *APIInterfaceMock) Initialize() error {
	if mock.InitializeFunc == nil {
		panic("APIInterfaceMock.InitializeFunc: method is nil but APIInterface.Initialize was just called")
	}
	callInfo := struct {
	}{}
	mock.lockInitialize.Lock()
	mock.calls.Initialize = append(mock.calls.Initialize, callInfo)
	mock.lockInitialize.Unlock()
	return mock.InitializeFunc()
}

// InitializeCalls gets all the calls that were made to Initialize.
// Check the length with:
//     len(mockedAPIInterface.InitializeCalls())
func (mock *APIInterfaceMock) InitializeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockInitialize.RLock()
	calls = mock.calls.Initialize
	mock.lockInitialize.RUnlock()
	return calls
}

// SetEchoInstance calls SetEchoInstanceFunc.
func (mock *APIInterfaceMock) SetEchoInstance(e *echo.Echo) {
	if mock.SetEchoInstanceFunc == nil {
		panic("APIInterfaceMock.SetEchoInstanceFunc: method is nil but APIInterface.SetEchoInstance was just called")
	}
	callInfo := struct {
		E *echo.Echo
	}{
		E: e,
	}
	mock.lockSetEchoInstance.Lock()
	mock.calls.SetEchoInstance = append(mock.calls.SetEchoInstance, callInfo)
	mock.lockSetEchoInstance.Unlock()
	mock.SetEchoInstanceFunc(e)
}

// SetEchoInstanceCalls gets all the calls that were made to SetEchoInstance.
// Check the length with:
//     len(mockedAPIInterface.SetEchoInstanceCalls())
func (mock *APIInterfaceMock) SetEchoInstanceCalls() []struct {
	E *echo.Echo
} {
	var calls []struct {
		E *echo.Echo
	}
	mock.lockSetEchoInstance.RLock()
	calls = mock.calls.SetEchoInstance
	mock.lockSetEchoInstance.RUnlock()
	return calls
}
