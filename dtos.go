package weoscontroller

import "github.com/wepala/weos"

type APIConfig struct {
	*weos.ApplicationConfig
	RecordingBaseFolder string
	Middleware          []string `json:"middleware"`
	PreMiddleware       []string `json:"pre-middleware"`
}

type PathConfig struct {
	Handler    string   `json:"handler" ,yaml:"handler"`
	Group      bool     `json:"group" ,yaml:"group"`
	Middleware []string `json:"middleware"`
}
