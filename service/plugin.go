package service

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"plugin"
	"reflect"
)

//define an interface that all plugins must implement
type PluginInterface interface {
	GetHandlerByName(name string) http.HandlerFunc
}

var plugins = make(map[string]PluginInterface)

func GetPlugin(fileName string) (PluginInterface, error) {
	if plugins[fileName] == nil {
		// Open - Loads the plugin
		log.Debugf("Loading plugin %s", fileName)
		p, err := plugin.Open(fileName)
		if err != nil {
			log.Errorf("Unable to log plugin '%s' because of error '%s'", fileName, err)
			return nil, err
		}

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
		plugins[fileName] = weosPlugin
	}

	return plugins[fileName], nil
}
