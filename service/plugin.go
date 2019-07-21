package service

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"plugin"
)

//define an interface that all plugins must implement
type PluginInterface interface {
	GetHandlerByName(name string) *http.HandlerFunc
}

var plugins = make(map[string]PluginInterface)

func GetPlugin(fileName string) (PluginInterface, error) {
	if plugins[fileName] == nil {
		// Open - Loads the plugin
		log.Debugf("Loading plugin %s", fileName)
		p, err := plugin.Open(fileName)
		if err != nil {
			panic(err)
		}

		//load the middleware object
		symbol, err := p.Lookup("Plugin")
		if err != nil {
			log.Errorf("could not load plugin")
		}
		// symbol - Checks the function signature
		weosPlugin, ok := symbol.(PluginInterface)
		if !ok {
			return nil, errors.New("plugin does not implement PluginInterface")
		}
		plugins[fileName] = weosPlugin
	}

	return plugins[fileName], nil
}
