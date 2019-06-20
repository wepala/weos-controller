package service

import (
	"testing"
)

func TestNewControllerService(t *testing.T) {
	t.Run("test basic yaml loaded", func(t *testing.T) {
		configYaml := "testdata/api/basic-site.yml"
		service, err := NewControllerService(configYaml, "")
		if err != nil {
			t.Fatalf("there was an error setting up service: %v", err)
		}

		if service.GetConfig() == nil {
			t.Fatalf("failed to load config: '%s'", configYaml)
		}

		//test loading the swagger file
		if service.GetConfig().ApiConfig.Info.Title != "Basic Site" {
			t.Errorf("expected the api title to be: '%s', got: '%s", "Basic Site", service.GetConfig().ApiConfig.Info.Title)
		}

		//check that the path is parsed
		pathConfig, err := service.GetPathConfig("/", "GET")
		if err != nil {
			t.Fatalf("issue getting path config: '%v", err)
		}

		if pathConfig == nil {
			t.Fatalf("pathconfig for path '/' not loaded")
		}

	})
}
