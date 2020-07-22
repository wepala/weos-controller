package test

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PluginInterface interface {
	GetHandlerByName(name string) http.HandlerFunc
	AddConfig(config json.RawMessage) error
	AddPathConfig(handler string, config json.RawMessage) error
	AddLogger(logger log.Ext1FieldLogger)
}

type WeOSPlugin struct {
}

func (m WeOSPlugin) AddLogger(logger log.Ext1FieldLogger) {

}

//func (m WeOSPlugin) AddConfig(config interface{}) error {
//	return nil
//}

func (m WeOSPlugin) GetHandlerByName(name string) http.HandlerFunc {
	if name == "HelloWorld" {
		return m.HelloWorld
	}

	if name == "FooBar" {
		return m.FooBar
	}

	return nil
}

func (m WeOSPlugin) AddConfig(config json.RawMessage) error {
	return nil
}

func (m WeOSPlugin) AddPathConfig(handler string, config json.RawMessage) error {
	return nil
}

func (m WeOSPlugin) HelloWorld(rw http.ResponseWriter, r *http.Request) {
	io := bytes.NewBufferString("Hello World")
	_, err := rw.Write(io.Bytes())
	r = r.WithContext(context.WithValue(r.Context(), "title", "Hello World"))
	if err != nil {
		//TODO log with fields so that we know which plugin, file and handler to look for to debug
	}
}

func (m WeOSPlugin) FooBar(rw http.ResponseWriter, r *http.Request) {

	if title := r.Context().Value("title"); title != nil {
		rw.WriteHeader(http.StatusOK)
		_, err := rw.Write([]byte(title.(string)))
		if err != nil {
			//TODO log with fields so that we know which plugin, file and handler to look for to debug
		}
	} else {
		io := bytes.NewBufferString("Foobar")
		_, err := rw.Write(io.Bytes())
		if err != nil {
			//TODO log with fields so that we know which plugin, file and handler to look for to debug
		}
	}

}

var WePlugin = WeOSPlugin{}
