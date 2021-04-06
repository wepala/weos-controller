package weoscontroller_test

import (
	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
	"net/http"
	"os"
	"testing"
)

type TestAPI struct {
	*weoscontroller.API
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

	weoscontroller.Configure(e, "../fixtures/api/api.yaml", testPlugin)

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitializeCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitializeCalls()))
	}

	if len(e.Routes()) != 23 {
		t.Errorf("expected %d route, got %d", 23, len(e.Routes()))
	}
}
