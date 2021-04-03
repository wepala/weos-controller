package echo_framework_test

import (
	"github.com/labstack/echo/v4"
	echo2 "github.com/wepala/weos-controller/echo_framework"
	"github.com/wepala/weos/module"
	"net/http"
	"os"
	"testing"
)

type TestPlugin struct {
	*echo2.APIPlugin
	plugin *PluginInterfaceMock
}

func (t *TestPlugin) InitModules(mod *module.WeOSMod) {
	t.plugin.InitModules(mod)
}

func (*TestPlugin) HealthChecker(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func TestStart(t *testing.T) {
	e := echo.New()
	err := os.Setenv("POSTGRES_HOST", "localhost")
	if err != nil {
		t.Fatalf("error setting up environment variables '%s'", err)
	}
	plugin := &PluginInterfaceMock{
		InitModulesFunc: func(mod *module.WeOSMod) {

		},
	}

	apiPlugin := echo2.NewAPIPlugin(e)

	testPlugin := &TestPlugin{
		apiPlugin,
		plugin,
	}

	echo2.Configure(e, "../fixtures/api/api.yaml", testPlugin)

	if testPlugin.Config.Database.Host != "localhost" {
		t.Errorf("expected the database host to be '%s', got '%s'", "localhost", testPlugin.Config.Database.Host)
	}

	if len(plugin.InitModulesCalls()) != 1 {
		t.Errorf("expected init modules to be called %d time, called %d times", 1, len(plugin.InitModulesCalls()))
	}

	if len(e.Routes()) != 23 {
		t.Errorf("expected %d route, got %d", 23, len(e.Routes()))
	}
}
