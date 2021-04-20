package weoscontroller

import (
	"github.com/wepala/weos"
)

type APIConfig struct {
	*weos.ApplicationConfig
	RecordingBaseFolder string
	Middleware          []string   `json:"middleware"`
	PreMiddleware       []string   `json:"pre-middleware"`
	JWTConfig           *JWTConfig `json:"jwtConfig"`
}

type PathConfig struct {
	Handler    string   `json:"handler" ,yaml:"handler"`
	Group      bool     `json:"group" ,yaml:"group"`
	Middleware []string `json:"middleware"`
}

type JWTConfig struct {
	Key           string                 `json:"key"`
	SigningKeys   map[string]interface{} `json:"signingKeys"`
	KeyPath       string                 `json:"keyPath"`
	TokenLookup   string                 `json:"tokenLookup"`
	Claims        map[string]interface{} `json:"claims"`
	AuthScheme    string                 `json:"authScheme"`
	ContextKey    string                 `json:"contextKey"`
	SigningMethod string                 `json:"signingMethod"`
}
