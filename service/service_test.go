package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
)

type Config struct {
	Mysql struct {
		Host     string `json:"host" yaml:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"mysql"`
}

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
			t.Fatalf("expected 1 middleware to be configured, got %d", len(aboutPathConfig.Middleware))
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
	var handlerNames []string
	config := Config{}
	//setup mock
	weosPluginMock := &PluginInterfaceMock{
		GetHandlerByNameFunc: func(name string) http.HandlerFunc {
			return func(writer http.ResponseWriter, request *http.Request) {
				handlerNames = append(handlerNames, name)
			}
		},
		AddConfigFunc: func(tconfig json.RawMessage) error {
			//check the config on the middleware
			tbytes, err := tconfig.MarshalJSON()
			if err != nil {
				t.Fatalf("encountered error marshaling json for config")
			}
			if err = json.Unmarshal(tbytes, &config); err != nil {
				t.Fatalf("encountered error unmarshaling json for config")
			}

			return nil
		},
		AddPathConfigFunc: func(handler string, config json.RawMessage) error {
			return nil
		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

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

	if config.Mysql.Host != "localhost" {
		t.Errorf("exepcted mysql host to be %s", "localhost")
	}

	if config.Mysql.User != "root" {
		t.Errorf("exepcted mysql user to be %s", "root")
	}

	if config.Mysql.Password != "root" {
		t.Errorf("exepcted mysql password to be %s", "root")
	}

	if len(weosPluginMock.AddPathConfigCalls()) != 1 {
		t.Errorf("expected add handler config to be called %d time, called %d times", 1, len(weosPluginMock.AddPathConfigCalls()))
	}

}

func TestControllerService_HandlerPriority(t *testing.T) {
	apiYaml := "testdata/api/basic-site-api." + runtime.GOOS + ".yml"
	var handlerNames []string
	config := Config{}
	//setup mock
	weosPluginMock1 := &PluginInterfaceMock{
		GetHandlerByNameFunc: func(name string) http.HandlerFunc {
			return func(writer http.ResponseWriter, request *http.Request) {
				handlerNames = append(handlerNames, name)
			}
		},
		AddConfigFunc: func(tconfig json.RawMessage) error {
			//check the config on the middleware
			tbytes, err := tconfig.MarshalJSON()
			if err != nil {
				t.Fatalf("encountered error marshaling json for config")
			}
			if err = json.Unmarshal(tbytes, &config); err != nil {
				t.Fatalf("encountered error unmarshaling json for config")
			}

			return nil
		},
		AddPathConfigFunc: func(handler string, config json.RawMessage) error {
			return nil
		},
	}

	weosPluginMock2 := &PluginInterfaceMock{
		GetHandlerByNameFunc: func(name string) http.HandlerFunc {
			return func(writer http.ResponseWriter, request *http.Request) {
				handlerNames = append(handlerNames, name)
			}
		},
		AddConfigFunc: func(tconfig json.RawMessage) error {
			//check the config on the middleware
			tbytes, err := tconfig.MarshalJSON()
			if err != nil {
				t.Fatalf("encountered error marshaling json for config")
			}
			if err = json.Unmarshal(tbytes, &config); err != nil {
				t.Fatalf("encountered error unmarshaling json for config")
			}

			return nil
		},
		AddPathConfigFunc: func(handler string, config json.RawMessage) error {
			return nil
		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			if strings.Contains(fileName, "testdata/plugins/test2") {
				return weosPluginMock2, nil
			}
			return weosPluginMock1, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

	//get path config
	pathConfig, err := s.GetPathConfig("/multiple-handlers", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig)
	if len(handlers) != 2 {
		t.Errorf("expected %d handlers to be loaded: got %d [%s]", 2, len(handlers), strings.Join(handlerNames, ","))
	}

	rw := httptest.NewRecorder()
	r := httptest.NewRequest("get", "/foo", nil)
	handlers[0].ServeHTTP(rw, r)

	if len(handlerNames) != 1 {
		t.Fatalf("handlers were not called")
	}

	if handlerNames[0] != "FooBar" {
		t.Errorf("expected the first handler to be %s, got %s", "FooBar", handlerNames[0])
	}
}
