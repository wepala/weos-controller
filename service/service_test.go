package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"
	"testing"
)

func TestNewControllerService(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Run("test basic yaml loaded", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := "testdata/api/basic-site-config." + runtime.GOOS + ".yml"
		service, err := service.NewControllerService(apiYaml, configYaml, nil)
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

		aboutPathConfig, err := service.GetPathConfig("/about", "get")

		if len(aboutPathConfig.Middleware) != 1 {
			t.Errorf("expected 1 middleware to be configured, got %d", len(aboutPathConfig.Middleware))
		}

	})
	t.Run("test middleware must be an array", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := "testdata/api/basic-site-middleware-error-config.yml"
		_, err := service.NewControllerService(apiYaml, configYaml, nil)
		if err == nil || !strings.Contains(err.Error(), "Middleware") {
			t.Fatalf("expected an error 'the list of templates must be an array in the config' got: %v", err)
		}
	})
	t.Run("test loading api config only", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api.yml"
		configYaml := ""
		service, err := service.NewControllerService(apiYaml, configYaml, nil)
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
			t.Fatalf("pathconfig for path '/about' not loaded")
		}
	})
}

func TestControllerService_GetHandlers(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	apiYaml := "testdata/api/basic-site-api.yml"
	configYaml := "testdata/api/basic-site-config." + runtime.GOOS + ".yml"
	handlerNames := make([]string, 1)
	//setup mock
	weosPluginMock := &PluginInterfaceMock{
		GetHandlerByNameFunc: func(name string) http.HandlerFunc {
			return func(writer http.ResponseWriter, request *http.Request) {
				handlerNames = append(handlerNames, name)
			}
		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, configYaml, pluginLoaderMock)

	//get path config
	pathConfig, err := s.GetPathConfig("/about", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig)
	if len(handlers) != 1 {
		t.Errorf("expected %d handlers to be loaded: got %d [%s]", 1, len(handlers), strings.Join(handlerNames, ","))
	}

}
