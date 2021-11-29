//go:generate moq -out mocks_test.go -pkg integration_test . TestAPI
//generating did not work in this package so generated the mocks outside and then brought them back into the integration package
package integration_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	weoscontroller "github.com/wepala/weos-controller"
)

type TestAPI interface {
	weoscontroller.APIInterface
	GlobalMiddleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc
	PreGlobalMiddleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc
	Middleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc
	PreMiddleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc
	FooBar(c echo.Context) error
	HelloWorld(c echo.Context) error
	Context(handlerFunc echo.HandlerFunc) echo.HandlerFunc
}

//loadHttpRequestFixture wrapper around the test helper to make it easier to use it with test table
func loadHttpRequestFixture(filename string, t *testing.T) *http.Request {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	reader := bufio.NewReader(bytes.NewReader(data))
	request, err := http.ReadRequest(reader)
	if err == io.EOF {
		return request
	}

	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	actualRequest, err := http.NewRequest(request.Method, request.URL.String(), reader)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}
	return actualRequest
}

func TestMiddlware(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user/1/2", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	var echoInstance *echo.Echo
	var middlewareAndHandlersCalled []string
	//setup a mock api with handlers and middleware
	api := &TestAPIMock{
		InitializeFunc: func() error {
			return nil
		},
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		SetEchoInstanceFunc: func(e *echo.Echo) {
			echoInstance = e
		},
		EchoInstanceFunc: func() *echo.Echo {
			return echoInstance
		},
		FooBarFunc: func(c echo.Context) error {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "fooBarHandler")
			return nil
		},
		ContextFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "contextMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		ZapLoggerFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "zapLoggerMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		GlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "globalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		PreGlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preGlobalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		MiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "middleware")
				return nil
			}
		},
		PreMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc { //run the middleware before calling the handler
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
	}
	weoscontroller.Initialize(e, api, "../fixtures/api/integration.yaml")
	e.ServeHTTP(rec, req)

	//check that the expected handlers and middleware are called
	if len(middlewareAndHandlersCalled) != 7 {
		t.Fatalf("expected %d middlware and handers to be called, got %d", 6, len(middlewareAndHandlersCalled))
	}

	//check the order in which the middleware and handlers are called
	if middlewareAndHandlersCalled[0] != "preGlobalMiddleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 0, "preGlobalMiddleware", middlewareAndHandlersCalled[0])
	}

	if middlewareAndHandlersCalled[1] != "contextMiddleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 1, "contextMiddleware", middlewareAndHandlersCalled[1])
	}

	if middlewareAndHandlersCalled[2] != "globalMiddleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 2, "globalMiddleware", middlewareAndHandlersCalled[2])
	}

	if middlewareAndHandlersCalled[3] != "zapLoggerMiddleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 3, "zapLoggerMiddleware", middlewareAndHandlersCalled[3])
	}

	if middlewareAndHandlersCalled[4] != "preMiddleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 4, "preMiddleware", middlewareAndHandlersCalled[4])
	}

	if middlewareAndHandlersCalled[5] != "fooBarHandler" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 5, "fooBarHandler", middlewareAndHandlersCalled[5])
	}

	if middlewareAndHandlersCalled[6] != "middleware" {
		t.Errorf("expected middleware or handler in position %d to be '%s', got '%s'", 6, "middleware", middlewareAndHandlersCalled[6])
	}

	if len(api.GlobalMiddlewareCalls()) != 1 {
		t.Errorf("expected %d call to global middleware, got '%d", 1, len(api.GlobalMiddlewareCalls()))
	}

	if len(api.FooBarCalls()) != 1 {
		t.Errorf("expected %d call, got %d", 1, len(api.FooBarCalls()))
	}

	if len(api.PreGlobalMiddlewareCalls()) != 1 {
		t.Errorf("expected %d call to global pre middleware, got %d", 1, len(api.PreGlobalMiddlewareCalls()))
	}

	if len(api.ContextCalls()) != 1 {
		t.Errorf("expected %d call to global pre middleware, got %d", 1, len(api.ContextCalls()))
	}

	if len(api.AddPathConfigCalls()) < 1 {
		t.Error("expected the path config to be called")
	}
}
func TestMiddleware_CORSTest(t *testing.T) {
	e := echo.New()
	var echoInstance *echo.Echo
	var middlewareAndHandlersCalled []string
	//setup a mock api with handlers and middleware
	api := &TestAPIMock{
		InitializeFunc: func() error {
			return nil
		},
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		SetEchoInstanceFunc: func(e *echo.Echo) {
			echoInstance = e
		},
		EchoInstanceFunc: func() *echo.Echo {
			return echoInstance
		},
		FooBarFunc: func(c echo.Context) error {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "fooBarHandler")
			return nil
		},
		ContextFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "contextMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		ZapLoggerFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "zapLoggerMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		GlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "globalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		PreGlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preGlobalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		MiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "middleware")
				return nil
			}
		},
		PreMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc { //run the middleware before calling the handler
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
	}
	weoscontroller.Initialize(e, api, "../fixtures/api/integration.yaml")

	t.Run("test cors put", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/putpoint/1/2", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Access-Control-Request-Method", "OPTIONS")
		req.Header.Set(echo.HeaderOrigin, "http://localhost:8682")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		response := rec.Result()

		//check response code
		if response.StatusCode != 204 {
			t.Errorf("expected response code to be %d, got %d", 204, response.StatusCode)
		}

		if !strings.Contains(response.Header.Get(echo.HeaderAccessControlAllowMethods), http.MethodPut) {
			t.Errorf("expected '%s' to be in the allowed methods, got '%s'", http.MethodPut, response.Header.Get(echo.HeaderAccessControlAllowMethods))
		}
	})
}

func TestErrorResponse(t *testing.T) {
	e := echo.New()
	var echoInstance *echo.Echo
	var middlewareAndHandlersCalled []string
	//setup a mock api with handlers and middleware
	api := &TestAPIMock{
		InitializeFunc: func() error {
			return nil
		},
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		SetEchoInstanceFunc: func(e *echo.Echo) {
			echoInstance = e
		},
		EchoInstanceFunc: func() *echo.Echo {
			return echoInstance
		},
		FooBarFunc: func(c echo.Context) error {
			return weoscontroller.NewControllerError("some error", errors.New("Some Detailed Error"), 405)
		},
		ContextFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "contextMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		ZapLoggerFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "zapLoggerMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		GlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "globalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		PreGlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preGlobalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		MiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "middleware")
				return nil
			}
		},
		PreMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc { //run the middleware before calling the handler
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
	}
	weoscontroller.Initialize(e, api, "../fixtures/api/integration.yaml")

	t.Run("test error response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/putpoint/1/2", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		response := rec.Result()

		//check response code
		if response.StatusCode != 405 {
			t.Errorf("expected response code to be %d, got %d", 405, response.StatusCode)
		}
	})
}

func TestZapLogger(t *testing.T) {
	e := echo.New()
	var echoInstance *echo.Echo
	var middlewareAndHandlersCalled []string
	//setup a mock api with handlers and middleware
	api := &TestAPIMock{
		InitializeFunc: func() error {
			return nil
		},
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		SetEchoInstanceFunc: func(e *echo.Echo) {
			echoInstance = e
		},
		EchoInstanceFunc: func() *echo.Echo {
			return echoInstance
		},
		ZapLoggerFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "zapLoggerMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		FooBarFunc: func(c echo.Context) error {
			return weoscontroller.NewControllerError("some error", errors.New("Some Detailed Error"), 405)
		},
		ContextFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "contextMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},

		GlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "globalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		PreGlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preGlobalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		MiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "middleware")
				return nil
			}
		},
		PreMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc { //run the middleware before calling the handler
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
	}
	weoscontroller.Initialize(e, api, "../fixtures/api/integration.yaml")
	prefix := api.EchoInstance().Logger.Prefix()
	level := api.EchoInstance().Logger.Level()

	if prefix != "zap" {
		t.Errorf("expected default logger to be zap but got %s ", prefix)
	}

	if level != log.ERROR {
		t.Errorf("expected default logger level to be erro but got %d ", level)
	}
}

func TestContextZapLogger(t *testing.T) {
	e := echo.New()
	var echoInstance *echo.Echo
	var middlewareAndHandlersCalled []string
	//setup a mock api with handlers and middleware
	api := &TestAPIMock{
		InitializeFunc: func() error {
			return nil
		},
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		SetEchoInstanceFunc: func(e *echo.Echo) {
			echoInstance = e
		},
		EchoInstanceFunc: func() *echo.Echo {
			return echoInstance
		},
		ZapLoggerFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "zapLoggerMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		FooBarFunc: func(c echo.Context) error {
			return weoscontroller.NewControllerError("some error", errors.New("Some Detailed Error"), 405)
		},
		ContextFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "contextMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},

		GlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "globalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		PreGlobalMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preGlobalMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
		MiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "middleware")
				return nil
			}
		},
		PreMiddlewareFunc: func(handlerFunc echo.HandlerFunc) echo.HandlerFunc { //run the middleware before calling the handler
			return func(c echo.Context) error {
				middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "preMiddleware")
				if err := handlerFunc(c); err != nil {
					c.Error(err)
				}
				return nil
			}
		},
	}
	weoscontroller.Initialize(e, api, "../fixtures/api/integration.yaml")

	req := httptest.NewRequest(http.MethodGet, "/user/1/2", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

}
