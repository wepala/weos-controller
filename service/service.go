package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boj/redistore"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface PluginInterface PluginLoaderInterface
type PathConfig struct {
	Mock       bool                `yaml:"mock"`
	Middleware []*MiddlewareConfig `yaml:"middleware"`
	Session    *SessionConfig      `yaml:"session"`
	Data       interface{}
}

type SessionStoreConfig struct {
	Type     string `yaml:"type"`
	Size     int    `yaml:"size"`
	Network  string `yaml:"network"`
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type SessionConfig struct {
	Key         string              `yaml:"key"`
	StoreConfig *SessionStoreConfig `json:"store"`
	MaxAge      int                 `json:"max-age"`
	Path        string              `json:"path"`
	Domain      string              `json:"domain"`
	Secure      bool                `json:"secure"`
	HttpOnly    bool                `json:"http-only"`
	SameSite    string              `json:"same-site"`
}

type MiddlewareConfig struct {
	Plugin struct {
		FileName string           `yaml:"filename"`
		Config   *json.RawMessage `yaml:"config"`
	} `yaml:"plugin"`
	Handler  string                 `yaml:"handler"`
	Context  map[string]interface{} `yaml:"context"`
	Priority int                    `yaml:"priority"`
}

type controllerService struct {
	config       *openapi3.Swagger
	pluginLoader PluginLoaderInterface
	logger       log.Ext1FieldLogger
	globalConfig *PathConfig
	session      sessions.Store
}

func (s *controllerService) GetGlobalConfig() (*PathConfig, error) {
	if s.globalConfig == nil {
		if s.config.ExtensionProps.Extensions["x-weos-config"] != nil {
			globalConfigBytes, err := s.config.ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
			if err != nil {
				return nil, err
			}
			var globalConfig *PathConfig
			err = json.Unmarshal(globalConfigBytes, &globalConfig)
			if err != nil {
				return nil, err
			}

			return globalConfig, nil
		}
	}

	return nil, nil
}

func (s *controllerService) GetGlobalMiddlewareConfig() ([]*MiddlewareConfig, error) {
	globalConfig, err := s.GetGlobalConfig()
	if err != nil {
		return nil, err
	}
	if globalConfig != nil {
		return globalConfig.Middleware, nil
	}
	return nil, nil
}

func (s *controllerService) GetPathConfig(path string, operation string) (*PathConfig, error) {
	weosConfig := s.config.Paths[path].GetOperation(strings.ToUpper(operation)).ExtensionProps.Extensions["x-weos-config"]
	if weosConfig == nil {
		return nil, nil
	}
	bytes, err := s.config.Paths[path].GetOperation(strings.ToUpper(operation)).ExtensionProps.Extensions["x-weos-config"].(json.RawMessage).MarshalJSON()
	if err != nil {
		return nil, err
	}
	pathConfig := PathConfig{}
	if err = json.Unmarshal(bytes, &pathConfig); err != nil {
		return nil, err
	}
	return &pathConfig, nil
}

func (s *controllerService) GetConfig() *openapi3.Swagger {
	return s.config
}

func (s *controllerService) GetHandlers(config *PathConfig, mockHandler http.Handler) ([]http.HandlerFunc, error) {
	globalHandlers, err := s.GetGlobalMiddlewareConfig()
	var middlewareConfig []*MiddlewareConfig

	if err != nil {
		s.logger.Debug("there was an issue loading global handlers")
		return nil, err
	}

	if config != nil {
		middlewareConfig = config.Middleware
	}
	middlewareConfig = append(middlewareConfig, globalHandlers...)
	handlers := make([]http.HandlerFunc, len(middlewareConfig))

	//WEOS-168 if there are no handlers or the config has mock set to true return mock handlers
	if config == nil || len(handlers) == 0 || config.Mock {
		return []http.HandlerFunc{mockHandler.ServeHTTP}, nil
	} else { // otherwise let's load the plugins
		sort.Sort(NewMiddlewareConfigSorter(middlewareConfig))
		for key, mc := range middlewareConfig {
			s.logger.Debugf("loading plugin %s", mc.Plugin.FileName)
			plugin, err := s.pluginLoader.GetPlugin(mc.Plugin.FileName)
			if err != nil {
				s.logger.Errorf("error loading plugin %s", err)
				return nil, err
			}

			if mc.Plugin.Config != nil {
				err = plugin.AddConfig(*mc.Plugin.Config) //pass the raw json message that is loaded to the plugin
				plugin.AddLogger(s.logger)
				if s.session != nil {
					plugin.AddSession(s.session)
				}
				if err != nil {
					log.Error("error loading plugin config", err)
					return nil, err
				}
			}

			s.logger.Debugf("retrieving handler %s", mc.Handler)
			handlers[key] = plugin.GetHandlerByName(mc.Handler)
		}
	}

	return handlers, nil
}

var api openapi3.Swagger

type ServiceInterface interface {
	GetPathConfig(path string, operation string) (*PathConfig, error)
	GetConfig() *openapi3.Swagger
	GetHandlers(config *PathConfig, mockHandler http.Handler) ([]http.HandlerFunc, error)
	GetGlobalMiddlewareConfig() ([]*MiddlewareConfig, error)
}

func NewControllerService(apiConfig string, pluginLoader PluginLoaderInterface) (ServiceInterface, error) {

	loader := openapi3.NewSwaggerLoader()
	file, err := ioutil.ReadFile(apiConfig)
	//replace environment variables in file
	file = []byte(os.ExpandEnv(string(file)))
	swagger, err := loader.LoadSwaggerFromData(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error loading %s: %s", apiConfig, err.Error()))
	}

	svc := &controllerService{
		config:       swagger,
		pluginLoader: pluginLoader,
		logger:       log.WithField("application", "weos-controller"),
	}

	//check if session is configured and then setup session
	if globalConfig, err := svc.GetGlobalConfig(); err == nil && globalConfig != nil {
		if globalConfig.Session != nil {
			if globalConfig.Session.StoreConfig != nil {
				switch globalConfig.Session.StoreConfig.Type {
				case "redis":
					svc.session, err = redistore.NewRediStore(globalConfig.Session.StoreConfig.Size, globalConfig.Session.StoreConfig.Network, globalConfig.Session.StoreConfig.Address, globalConfig.Session.StoreConfig.Password, []byte(globalConfig.Session.Key))
					if err != nil {
						svc.logger.Errorf("encountered error starting redis session %s", err)
						svc.logger.Info("Falling back to default cookie store")
						svc.session = setupDefaultSession(globalConfig.Session)
					} else {
						svc.session.(*redistore.RediStore).Options = &sessions.Options{
							Path:     globalConfig.Session.Path,
							Domain:   globalConfig.Session.Domain,
							MaxAge:   globalConfig.Session.MaxAge,
							Secure:   globalConfig.Session.Secure,
							HttpOnly: globalConfig.Session.HttpOnly,
						}
						switch globalConfig.Session.SameSite {
						case "SameSiteNoneMode":
							svc.session.(*redistore.RediStore).Options.SameSite = http.SameSiteNoneMode
							break
						case "SameSiteLaxMode":
							svc.session.(*redistore.RediStore).Options.SameSite = http.SameSiteLaxMode
							break
						case "SameSiteStrictMode":
							svc.session.(*redistore.RediStore).Options.SameSite = http.SameSiteStrictMode
							break
						default:
							svc.session.(*redistore.RediStore).Options.SameSite = http.SameSiteDefaultMode
							break
						}
					}

					break
				default:
					svc.session = sessions.NewCookieStore([]byte(globalConfig.Session.Key))
				}

			} else {
				svc.session = setupDefaultSession(globalConfig.Session)
			}
		}
	}

	return svc, nil
}

func setupDefaultSession(config *SessionConfig) *sessions.CookieStore {
	session := sessions.NewCookieStore([]byte(config.Key))
	session.MaxAge(config.MaxAge)
	return session
}
