//go:generate moq -out plugin_mocks_test.go -pkg core_test . PluginInterface
package core

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wepala/weos/module"
	"plugin"
	"reflect"
)

//define an interface that all plugins must implement
type PluginInterface interface {
	AddConfig(config *APIConfig) error
	InitModules(mod *module.WeOSMod)
}

type RepositoryInterface interface {
	Get(name string) []interface{}
}

type PluginLoaderInterface interface {
	GetPlugin(fileName string) (PluginInterface, error)
	GetRepository(fileName string) (RepositoryInterface, error)
}

//monkey patch for opening plugin so testing is easier
var OpenPlugin = plugin.Open

//setup a login loader
type PluginLoader struct {
	plugins      map[string]PluginInterface
	repositories map[string]RepositoryInterface
}

func (loader *PluginLoader) GetRepository(fileName string) (RepositoryInterface, error) {
	var p *plugin.Plugin
	var err error

	//if the so hasn't been loaded for plugins or repositories then let's load the file
	if loader.plugins[fileName] == nil && loader.repositories[fileName] == nil {
		// Open - Loads the plugin
		log.Debugf("Loading plugin %s", fileName)
		p, err = OpenPlugin(fileName)
		if err != nil {
			log.Errorf("Unable to log plugin '%s' because of error '%s'", fileName, err)
			return nil, err
		}
	}

	if loader.repositories[fileName] == nil {
		//load the middleware object
		symbol, err := p.Lookup("WeRepository")
		if err != nil {
			log.Errorf("could not load repository")
			return nil, err
		}
		// symbol - Checks the function signature
		weosRepository, ok := symbol.(RepositoryInterface)
		if !ok {
			v := reflect.ValueOf(symbol)
			return nil, errors.New(fmt.Sprintf("plugin does not implement PluginInterface, it is type '%s'", v.Kind().String()))
		}
		loader.repositories[fileName] = weosRepository
	}

	return loader.repositories[fileName], nil
}

var NewPluginLoader = func() *PluginLoader {
	return &PluginLoader{plugins: make(map[string]PluginInterface)}
}

func (loader *PluginLoader) GetPlugin(fileName string) (PluginInterface, error) {
	var p *plugin.Plugin
	var err error

	//if the so hasn't been loaded for plugins or repositories then let's load the file
	if loader.plugins[fileName] == nil && loader.repositories[fileName] == nil {
		// Open - Loads the plugin
		log.Debugf("Loading plugin %s", fileName)
		p, err = OpenPlugin(fileName)
		if err != nil {
			log.Errorf("Unable to log plugin '%s' because of error '%s'", fileName, err)
			return nil, err
		}
	}

	if loader.plugins[fileName] == nil {
		//load the middleware object
		symbol, err := p.Lookup("WePlugin")
		if err != nil {
			log.Errorf("could not load plugin")
			return nil, err
		}
		// symbol - Checks the function signature
		weosPlugin, ok := symbol.(PluginInterface)
		if !ok {
			v := reflect.ValueOf(symbol)
			return nil, errors.New(fmt.Sprintf("plugin does not implement PluginInterface, it is type '%s'", v.Kind().String()))
		}
		loader.plugins[fileName] = weosPlugin
	}

	return loader.plugins[fileName], nil
}
