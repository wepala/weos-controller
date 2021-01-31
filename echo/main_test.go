package echo_test

import (
	"github.com/labstack/echo/v4"
	"github.com/wepala/weos-controller/core"
	echo2 "github.com/wepala/weos-controller/echo"
	"github.com/wepala/weos/module"
	"net/http"
	"os"
	"testing"
)

type TestPlugin struct {
	*PluginInterfaceMock
}

func (TestPlugin) HealthChecker(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TestStart(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &PluginInterfaceMock{
		AddConfigFunc: func(config *core.APIConfig) error {
			if config.Database.Host != "localhost" {
				t.Errorf("expected the database host to be '%s', got '%s'", "localhost", config.Database.Host)
			}
			return nil
		},
		InitModulesFunc: func(mod *module.WeOSMod) {

		},
	}

	testPlugin := &TestPlugin{
		plugin,
	}

	echo2.Configure(e, "../fixtures/api/api.yaml", testPlugin)

	if len(plugin.AddConfigCalls()) != 1 {
		t.Errorf("expected add config to be called %d time, called %d times", 1, len(plugin.AddConfigCalls()))
	}

	if len(plugin.InitModulesCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitModulesCalls()))
	}

	if len(e.Routes()) != 1 {
		t.Errorf("expected %d route, got %d", 1, len(e.Routes()))
	}
}
