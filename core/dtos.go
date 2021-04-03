package core

import weosmodule "github.com/wepala/weos/module"

type APIConfig struct {
	*weosmodule.WeOSModuleConfig
	RecordingBaseFolder string
	Middleware          []string `json:"middleware"`
}

type PathConfig struct {
	Handler    string   `json:"handler" ,yaml:"handler"`
	Group      bool     `json:"group" ,yaml:"group"`
	Middleware []string `json:"middleware"`
}
