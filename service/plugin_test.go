package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"plugin"
	"testing"
)

func TestGetPlugin(t *testing.T) {
	t.SkipNow()
	plugin, err := service.NewPluginLoader().GetPlugin("testdata/plugins/test.so")
	if plugin == nil {
		t.Errorf("expected plugin to be loaded, got error '%s'", err)
	}
}

func TestPluginLoadedOnce(t *testing.T) {
	t.SkipNow()
	timesPluginLoaded := 0

	//monkey patch the function "OpenPlugin" to increment the counter so we can confirm that it runs only once
	service.OpenPlugin = func(path string) (*plugin.Plugin, error) {
		timesPluginLoaded = timesPluginLoaded + 1
		return plugin.Open(path)
	}
	pluginLoader := service.NewPluginLoader()

	plugin, err := pluginLoader.GetPlugin("testdata/plugins/test.so")
	if plugin == nil {
		t.Fatalf("expected plugin to be loaded the first time, got error '%s'", err)
	}
	plugin, _ = pluginLoader.GetPlugin("testdata/plugins/test.so")
	if timesPluginLoaded != 1 {
		t.Errorf("expected plugin to be loaded once")
	}
}

func TestInvalidPluginNotLoaded(t *testing.T) {
	_, err := service.NewPluginLoader().GetPlugin("testdata/plugins/invalid_test.so")
	if err == nil {
		t.Errorf("expected error loading plugin")
	}
}

func TestPluginLoader_GetRepository(t *testing.T) {
	t.SkipNow()
	plugin, err := service.NewPluginLoader().GetRepository("testdata/plugins/test.so")
	if plugin == nil {
		t.Errorf("expected plugin to be loaded, got error '%s'", err)
	}
}
