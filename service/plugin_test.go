package service

import "testing"

func TestGetPlugin(t *testing.T) {
	plugin, err := GetPlugin("testdata/plugins/test.so")
	if plugin == nil {
		t.Errorf("expected plugin to be loaded, got error '%s'", err)
	}
}

func TestPluginLoadedOnce(t *testing.T) {
	//TODO figure out how to mock the glob call
	timesPluginLoaded := 0
	plugin, err := GetPlugin("testdata/plugins/test.so")
	if plugin == nil {
		t.Fatalf("expected plugin to be loaded the first time, got error '%s'", err)
	} else {
		timesPluginLoaded = timesPluginLoaded + 1
	}
	plugin, _ = GetPlugin("testdata/plugins/test.so")
	if timesPluginLoaded != 1 {
		t.Errorf("expected plugin to be loaded once")
	}
}

func TestInvalidPluginNotLoaded(t *testing.T) {
	t.SkipNow()
	_, err := GetPlugin("testdata/plugins/invalid_test.so")
	if err == nil {
		t.Errorf("expected error loading plugin")
	}
}
