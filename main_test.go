package weoscontroller_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
)

type TestAPI struct {
	weoscontroller.API
	plugin *APIInterfaceMock
}

func (t *TestAPI) Initialize() error {
	return t.plugin.Initialize()
}

func (*TestAPI) HealthChecker(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TestStart(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &APIInterfaceMock{
		InitializeFunc: func() error {
			return nil
		},
	}

	//we're only nesting the plugin interface for testing
	testPlugin := &TestAPI{
		plugin: plugin,
	}

	weoscontroller.Initialize(e, testPlugin, "./fixtures/api/api.yaml")

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitializeCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitializeCalls()))
	}

	if len(e.Routes()) != 23 {
		t.Errorf("expected %d route, got %d", 23, len(e.Routes()))
	}

	if testPlugin.API.EchoInstance() == nil {
		t.Errorf("expected echo instance to be set")
	}
}

func TestParsingJWTConfigurations(t *testing.T) {

	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	err = os.Setenv("JWT_KEY", "littleBooPeep")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	err = os.Setenv("JWT_LOOKUP", "bigtoe")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &APIInterfaceMock{
		InitializeFunc: func() error {
			return nil
		},
	}

	//we're only nesting the plugin interface for testing
	testPlugin := &TestAPI{
		plugin: plugin,
	}

	weoscontroller.Initialize(e, testPlugin, "./fixtures/api/api.yaml")

	if testPlugin.API.Config.JWTConfig.Key != "littleBooPeep" {
		t.Errorf("expected the jwt config key to be '%s', got '%s'", "littleBooPeep", testPlugin.Config.JWTConfig.Key)
	}

	if testPlugin.API.Config.JWTConfig.TokenLookup != "bigtoe" {
		t.Errorf("expected the jwt config key to be '%s', got '%s'", "bigtoe", testPlugin.Config.JWTConfig.TokenLookup)
	}

	if testPlugin.API.Config.JWTConfig.Claims["email"] != "string" {
		t.Errorf("expected the jwt config email claims to be '%s', got '%s'", "email@mail", testPlugin.Config.JWTConfig.Claims["email"])
	}

	if testPlugin.API.Config.JWTConfig.Claims["real"] != "bool" {
		t.Errorf("expected the jwt config real claims to be '%s', got '%s'", "real", testPlugin.Config.JWTConfig.Claims["real"])
	}

}
func TestParsingRoutesWithParams(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &APIInterfaceMock{
		InitializeFunc: func() error {
			return nil
		},
	}

	//we're only nesting the plugin interface for testing
	testPlugin := &TestAPI{
		plugin: plugin,
	}

	weoscontroller.Initialize(e, testPlugin, "./fixtures/api/routestest.yaml")

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitializeCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitializeCalls()))
	}

	if len(e.Routes()) != 2 {
		t.Errorf("expected %d route, got %d", 2, len(e.Routes()))
	}

	if e.Routes()[1].Path != "/user/:id/:contentID" {
		t.Errorf("expected the path to be '%s', got '%s'", "/user/:id/:contentID", e.Routes()[1].Path)
	}
}

func TestParsingRoutesWithBasePath(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &APIInterfaceMock{
		InitializeFunc: func() error {
			return nil
		},
	}

	//we're only nesting the plugin interface for testing
	testPlugin := &TestAPI{
		plugin: plugin,
	}

	weoscontroller.Initialize(e, testPlugin, "./fixtures/api/basepathtest.yaml")

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitializeCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitializeCalls()))
	}

	if len(e.Routes()) != 2 {
		t.Errorf("expected %d route, got %d", 2, len(e.Routes()))
	}

	if e.Routes()[0].Path != "/weos/health" {
		t.Errorf("expected the path to be '%s', got '%s'", "/weos/health", e.Routes()[0].Path)
	}
}

func TestSendingConfigContents(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &APIInterfaceMock{
		InitializeFunc: func() error {
			return nil
		},
	}

	//we're only nesting the plugin interface for testing
	testPlugin := &TestAPI{
		plugin: plugin,
	}

	weoscontroller.Initialize(e, testPlugin, `openapi: 3.0.2
info:
  title: WeOS REST API
  version: 1.0.0
  description:  REST API for passing information into WeOS

x-weos-config:
  basePath: /weos
  session:
    key: "${SESSION_KEY}"
    path: ""
  logger:
    level: ${LOG_LEVEL}
    report-caller: true
    formatter: ${LOG_FORMAT}
  applicationId: ${APPLICATION_ID}
  applicationTitle: ${APPLICATION_TITLE}
  accountId: ${ACCOUNT_ID}
  database:
    host: ${POSTGRES_HOST}
    database: ${POSTGRES_DB}
    username: ${POSTGRES_USER}
    password: ${POSTGRES_PASSWORD}
    port: ${POSTGRES_PORT}
  middleware:
    - RequestID
    - Recover
    - Static

paths:
  /health:
    summary: Health Check
    get:
      x-weos-config:
        handler: HealthChecker
      responses:
        "200":
          description: Health Response
        "500":
          description: API Internal Error
  /user/{id}/{contentID}:
    summary: Some user endpoint
    get:
      x-weos-config:
        handler: HealthChecker
        pre-middlware:
          - RequestRecording
      responses:
        200:
          description: Admin Endpoint`)

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitializeCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitializeCalls()))
	}

	if len(e.Routes()) != 2 {
		t.Errorf("expected %d route, got %d", 2, len(e.Routes()))
	}

	if e.Routes()[0].Path != "/weos/health" {
		t.Errorf("expected the path to be '%s', got '%s'", "/weos/health", e.Routes()[0].Path)
	}
}
