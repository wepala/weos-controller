package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"encoding/json"
	"github.com/boj/redistore"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
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
		apiYaml := "testdata/api/basic-site-api.yml"
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
		apiYaml := "testdata/api/basic-site-api.yml"
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
	apiYaml := "testdata/api/basic-site-api.yml"
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

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
	handlers, _ := s.GetHandlers(pathConfig, nil)
	if len(handlers) != 2 {
		t.Errorf("expected %d handlers to be loaded: got %d [%s]", 2, len(handlers), strings.Join(handlerNames, ","))
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

	//if len(weosPluginMock.AddPathConfigCalls()) != 1 {
	//	t.Errorf("expected add handler config to be called %d time, called %d times", 1, len(weosPluginMock.AddPathConfigCalls()))
	//}

}

func TestControllerService_HandlerPriority(t *testing.T) {
	apiYaml := "testdata/api/basic-site-api.yml"
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

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
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/"]})
	if len(handlers) != 3 {
		t.Errorf("expected %d handlers to be loaded: got %d [%s]", 3, len(handlers), strings.Join(handlerNames, ","))
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

func TestControllerService_GlobalHandlers(t *testing.T) {
	apiYaml := "testdata/api/basic-site-api.yml"
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
	pathConfig, err := s.GetPathConfig("/", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/"]})
	if len(handlers) != 1 {
		t.Errorf("expected %d handlers to be loaded: got %d [%s]", 1, len(handlers), strings.Join(handlerNames, ","))
	}
}

func Test_WEOS_168(t *testing.T) {
	t.Run("test mock when no plugins associated", func(t *testing.T) {
		apiYaml := "testdata/api/mock-api.yml"

		pluginLoaderMock := &PluginLoaderInterfaceMock{
			GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
				return nil, nil
			},
		}

		s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
		if err != nil {
			t.Fatalf("got an error while create new controller service %s", err)
		}

		//get path config
		pathConfig, err := s.GetPathConfig("/", "get")
		if err != nil {
			t.Fatalf("issue getting path config: '%v", err)
		}

		if len(pluginLoaderMock.GetPluginCalls()) > 0 {
			t.Fatalf("didn't expect any plugin to be loaded")
		}

		//use path config to get handlers
		handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/"]})

		if len(handlers) != 1 {
			t.Errorf("expected %d handlers to be loaded, got %d", 1, len(handlers))
		}
	})

	t.Run("test mock when config is set to true", func(t *testing.T) {
		apiYaml := "testdata/api/mock-api.yml"
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

		pluginLoaderMock := &PluginLoaderInterfaceMock{
			GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
				return weosPluginMock1, nil
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
		handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/about"]})

		if len(handlers) != 1 {
			t.Errorf("expected %d handlers to be loaded, got %d", 1, len(handlers))
		}

		if len(pluginLoaderMock.GetPluginCalls()) > 0 {
			t.Errorf("didn't expect the plugin to be loaded")
		}
	})
}

//Test_WEOS_482 make it so that environment variables can be passed for routes in weos controller
func Test_WEOS_482(t *testing.T) {
	apiYaml := "testdata/api/basic-site-api.yml"
	var handlerNames []string
	config := Config{}
	os.Setenv("ACCOUNT", "wepala")
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock1, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

	//get path config
	pathConfig, err := s.GetPathConfig("/wepala/users", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/wepala/users"]})

	if len(handlers) != 2 {
		t.Errorf("expected %d handlers to be loaded, got %d", 2, len(handlers))
	}

}

func Test_WECON_1(t *testing.T) {
	apiYaml := "testdata/api/wetutor-api.yaml"
	var handlerNames []string
	config := Config{}
	os.Setenv("ACCOUNT", "wepala")
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock1, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

	//get path config
	pathConfig, err := s.GetPathConfig("/health", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/health"]})

	if len(handlers) != 1 {
		t.Errorf("expected %d handlers to be loaded, got %d", 1, len(handlers))
	}

}

func Test_AddSession(t *testing.T) {
	apiYaml := "testdata/api/session-api.yml"
	var handlerNames []string
	config := Config{}
	os.Setenv("SESSION_KEY", "wepala")
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

		},
		AddSessionFunc: func(session sessions.Store) {
			//expect default key to be cookie store
			sessionStore, ok := session.(*sessions.CookieStore)
			if !ok {
				t.Errorf("expected the default session store to be a cookie store")
			}

			//expect default max age to be as long as the browser session
			if sessionStore.Options.MaxAge != 0 {
				t.Errorf("expected the default max age to be %d, got %d", 0, sessionStore.Options.MaxAge)
			}

			if sessionStore.Options.Path != "/some-path" {
				t.Errorf("expected the max age to be %s, got %s", "/some-path", sessionStore.Options.Path)
			}

			if sessionStore.Options.Domain != "http://weos.cloud" {
				t.Errorf("expected the max age to be %s, got %s", "http://weos.cloud", sessionStore.Options.Domain)
			}

			if !sessionStore.Options.Secure {
				t.Errorf("expected the secure option to be true")
			}

			if !sessionStore.Options.HttpOnly {
				t.Errorf("expected the http-only option to be true")
			}

			if sessionStore.Options.SameSite != http.SameSiteNoneMode {
				t.Errorf("expected the same site option to be set to SameSiteNoneMode")
			}
		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock1, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

	//get path config
	pathConfig, err := s.GetPathConfig("/wepala/users", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/wepala/users"]})

	if len(handlers) != 2 {
		t.Errorf("expected %d handlers to be loaded, got %d", 2, len(handlers))
	}

	//confirm session is setup
	if len(weosPluginMock1.AddSessionCalls()) == 0 {
		t.Errorf("expected session to be created and passed to plugin")
	}
}

//Test_RedisSession this test only work when run with docker-compose run redis test (we need redis running for this to work)
func Test_RedisSession(t *testing.T) {
	apiYaml := "testdata/api/session-redis-api.yml"
	var handlerNames []string
	config := Config{}
	os.Setenv("SESSION_KEY", "wepala")
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
		AddLoggerFunc: func(logger log.Ext1FieldLogger) {

		},
		AddSessionFunc: func(session sessions.Store) {
			//expect default key to be cookie store
			sessionStore, ok := session.(*redistore.RediStore)
			if !ok {
				t.Fatalf("expected the session store to be a redis store")
			}

			if sessionStore.Options.MaxAge != 86400 {
				t.Errorf("expected the max age to be %d, got %d", 86400, sessionStore.Options.MaxAge)
			}

			if sessionStore.Options.Path != "/some-path" {
				t.Errorf("expected the max age to be %s, got %s", "/some-path", sessionStore.Options.Path)
			}

			if sessionStore.Options.Domain != "http://weos.cloud" {
				t.Errorf("expected the max age to be %s, got %s", "http://weos.cloud", sessionStore.Options.Domain)
			}

			if !sessionStore.Options.Secure {
				t.Errorf("expected the secure option to be true")
			}

			if !sessionStore.Options.HttpOnly {
				t.Errorf("expected the http-only option to be true")
			}

			if sessionStore.Options.SameSite != http.SameSiteNoneMode {
				t.Errorf("expected the same site option to be set to SameSiteNoneMode")
			}
		},
	}

	pluginLoaderMock := &PluginLoaderInterfaceMock{
		GetPluginFunc: func(fileName string) (pluginInterface service.PluginInterface, e error) {
			return weosPluginMock1, nil
		},
	}

	s, err := service.NewControllerService(apiYaml, pluginLoaderMock)
	if err != nil {
		t.Fatalf("got an error while create new controller service %s", err)
	}

	//get path config
	pathConfig, err := s.GetPathConfig("/wepala/users", "get")
	if err != nil {
		t.Fatalf("issue getting path config: '%v", err)
	}

	//use path config to get handlers
	handlers, _ := s.GetHandlers(pathConfig, &service.MockHandler{PathInfo: s.GetConfig().Paths["/wepala/users"]})

	if len(handlers) != 2 {
		t.Errorf("expected %d handlers to be loaded, got %d", 2, len(handlers))
	}

	//confirm session is setup
	if len(weosPluginMock1.AddSessionCalls()) == 0 {
		t.Errorf("expected session to be created and passed to plugin")
	}
}
