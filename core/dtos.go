package core

import weosmodule "github.com/wepala/weos/module"

type APIConfig struct {
	*weosmodule.WeOSModuleConfig
}

type PathConfig struct {
	Handler string `json:"handler" ,yaml:"handler"`
}
