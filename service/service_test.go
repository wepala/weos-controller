package service

import (
	"testing"
)

func TestNewControllerService(t *testing.T) {
	t.Run("test basic yaml loaded", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := "testdata/api/basic-site-config.yml"
		service, err := NewControllerService(apiYaml, configYaml)
		if err != nil {
			t.Fatalf("there was an error setting up service: %v", err)
		}

		if service.GetConfig() == nil {
			t.Fatalf("failed to load config: '%s'", apiYaml)
		}

		//test loading the swagger file
		if service.GetConfig().ApiConfig.Info.Title != "Basic Site" {
			t.Errorf("expected the api title to be: '%s', got: '%s", "Basic Site", service.GetConfig().ApiConfig.Info.Title)
		}

		//check that the path is parsed. Note it was decided that the casing must match what is in the config. This can (should) be fixed in the future
		pathConfig, err := service.GetPathConfig("/", "get")
		if err != nil {
			t.Fatalf("issue getting path config: '%v", err)
		}

		if pathConfig == nil {
			t.Fatalf("pathconfig for path '/' not loaded")
		}

		if len(pathConfig.Templates) != 2 {
			t.Errorf("expected 2 templates to be configured, got %d", len(pathConfig.Templates))
		}

		if pathConfig.Data == nil {
			t.Errorf("expected data to be loaded")
		}

		aboutPathConfig, err := service.GetPathConfig("/about", "get")
		if len(aboutPathConfig.Templates) != 2 {
			t.Errorf("expected 2 templates to be configured, got %d", len(pathConfig.Templates))
		}

		if len(aboutPathConfig.Middleware) != 1 {
			t.Errorf("expected 1 middleware to be configured, got %d", len(pathConfig.Middleware))
		}

	})
	t.Run("test templates must be an array", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := "testdata/api/basic-site-template-error-config.yml"
		_, err := NewControllerService(apiYaml, configYaml)
		if err == nil || err.Error() != "the list of templates must be an array in the config" {
			t.Fatalf("expected an error 'the list of templates must be an array in the config' got: %v", err)
		}
	})
	t.Run("test middleware must be an array", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := "testdata/api/basic-site-middleware-error-config.yml"
		_, err := NewControllerService(apiYaml, configYaml)
		if err == nil || err.Error() != "the list of middlewares must be an array in the config" {
			t.Fatalf("expected an error 'the list of templates must be an array in the config' got: %v", err)
		}
	})
	t.Run("test loading api config only", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := ""
		service, err := NewControllerService(apiYaml, configYaml)
		if err != nil {
			t.Fatalf("there was an error setting up service: %v", err)
		}

		if service.GetConfig() == nil {
			t.Fatalf("failed to load config: '%s'", apiYaml)
		}

		//test loading the swagger file
		if service.GetConfig().ApiConfig.Info.Title != "Basic Site" {
			t.Errorf("expected the api title to be: '%s', got: '%s", "Basic Site", service.GetConfig().ApiConfig.Info.Title)
		}

		//check that the path is parsed. Note it was decided that the casing must match what is in the config. This can (should) be fixed in the future
		pathConfig, err := service.GetPathConfig("/", "get")
		if err != nil {
			t.Fatalf("issue getting path config: '%v", err)
		}

		if pathConfig == nil {
			t.Fatalf("pathconfig for path '/' not loaded")
		}
	})
}
