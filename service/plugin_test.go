package service

import "testing"

func TestGetPlugin(t *testing.T) {
	plugin, _ := GetPlugin("testdata/plugins/test.so")
	if plugin == nil {
		t.Errorf("expected plugin to be loaded")
	}
}

func TestPluginLoadedOnce(t *testing.T) {
	//TODO figure out how to mock the glob call
	timesPluginLoaded := 0
	plugin, _ := GetPlugin("testdata/test.so")
	if plugin == nil {
		t.Fatalf("expected plugin to be loaded the first time")
	}
	plugin, _ = GetPlugin("testdata/test.so")
	if timesPluginLoaded != 1 {
		t.Errorf("expected plugin to be loaded once")
	}
}

func TestInvalidPluginNotLoaded(t *testing.T) {
	_, err := GetPlugin("testdata/plugins/invalid_test.so")
	if err == nil {
		t.Errorf("expected error loading plugin")
	}
}
