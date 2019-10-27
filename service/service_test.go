package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"net/http"
	"runtime"
	"strings"
	"testing"
)

func TestNewControllerService(t *testing.T) {
	t.Run("test basic yaml loaded", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api." + runtime.GOOS + ".yml"
		testService, err := service.NewControllerService(apiYaml, nil)
		if err != nil {
			t.Fatalf("there was an error setting up testService: %v", err)
		}

		if testService.GetConfig() == nil {
			t.Fatalf("failed to load config: '%s'", apiYaml)
		}

		//test loading the swagger file
		if testService.GetConfig().Info.Title != "Basic Site" {
			t.Errorf("expected the api title to be: '%s', got: '%s", "Basic Site", testService.GetConfig().Info.Title)
		}

		aboutPathConfig, err := testService.GetPathConfig("/about", "get")

		if len(aboutPathConfig.Middleware) != 1 {
			t.Errorf("expected 1 middleware to be configured, got %d", len(aboutPathConfig.Middleware))
		}

	})
	t.Run("test middleware must be an array", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api-error.yml"
		testService, err := service.NewControllerService(apiYaml, nil)
		if err != nil {
			t.Fatalf("unable to instantiate new service, got error %s", err)
		}
		_, err = testService.GetPathConfig("/about", "get")
		if err == nil || !strings.Contains(err.Error(), "Middleware") {
			t.Fatalf("expected an error 'the list of templates must be an array in the config' got: %v", err)
		}
	})
	t.Run("test loading api config only", func(t *testing.T) {
		apiYaml := "testdata/api/basic-site-api." + runtime.GOOS + ".yml"
		service, err := service.NewControllerService(apiYaml, nil)
		if err != nil {
			t.Fatalf("there was an error setting up service: %v", err)
		}

		if service.GetConfig() == nil {
			t.Fatalf("failed to load config: '%s'", apiYaml)
		}

		//test loading the swagger file
		if service.GetConfig().Info.Title != "Basic Site" {
			t.Errorf("expected the api title to be: '%s', got: '%s", "Basic Site", service.GetConfig().Info.Title)
		}

		pathConfig, err := service.GetPathConfig("/about", "get")
		if err != nil {
			t.Fatalf("issue getting path config: '%v", err)
		}

		if pathConfig == nil {
			t.Fatalf("pathconfig for path '/about' not loaded")
		}
	})
}

func TestControllerService_GetHandlers(t *testing.T) {
	apiYaml := "testdata/api/basic-site-api." + runtime.GOOS + ".yml"
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

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)

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
