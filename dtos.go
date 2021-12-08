package weoscontroller

import (
	"encoding/json"

	"github.com/wepala/weos"
	"google.golang.org/grpc"
)

type APIConfig struct {
	*weos.ApplicationConfig
	BasePath            string `json:"basePath" ,yaml:"basePath"`
	RecordingBaseFolder string
	Middleware          []string         `json:"middleware"`
	PreMiddleware       []string         `json:"pre-middleware"`
	JWTConfig           *JWTConfig       `json:"jwtConfig"`
	Config              json.RawMessage  `json:"config"`
	Grpc                *GrpcMiddlewares `json:"grpc"`
}

type PathConfig struct {
	Handler        string          `json:"handler" ,yaml:"handler"`
	Group          bool            `json:"group" ,yaml:"group"`
	Middleware     []string        `json:"middleware"`
	Config         json.RawMessage `json:"config"`
	DisableCors    bool            `json:"disable-cors"`
	AllowedHeaders []string        `json:"allowed-headers" ,yaml:"allowed-headers"`
	AllowedOrigins []string        `json:"allowed-origins" ,yaml:"allowed-origins"`
}

type JWTConfig struct {
	Key             string                 `json:"key"`         //Signing key needed for validating token
	SigningKeys     map[string]interface{} `json:"signingKeys"` //Key map used for validating token.  Can be used in place of a single key
	Certificate     []byte                 `json:"certificate"`
	CertificatePath string                 `json:"certificatePath"` //Path the signing certificate used to validate token.  Can  be used in place of a key
	JWKSUrl         string                 `json:"jwksUrl"`         //URL to JSON Web Key set.  Can be used in place of a Key
	TokenLookup     string                 `json:"tokenLookup"`
	Claims          map[string]interface{} `json:"claims"`
	AuthScheme      string                 `json:"authScheme"`
	ContextKey      string                 `json:"contextKey"`
	SigningMethod   string                 `json:"signingMethod"`
}

type GRPCServerOptions struct {
	StreamMiddleware grpc.ServerOption
	UnaryMiddleware  grpc.ServerOption
}

type GrpcMiddlewares struct {
	Middlewares *Middlewares
}

type Middlewares struct {
	Stream *Stream `json:"stream"`
	Unary  *Unary  `json:"unary"`
}

type Stream struct {
	Middleware []string `json:"middleware"`
}

type Unary struct {
	Middleware []string `json:"middleware"`
}

type GRPCAPIConfig struct {
	*weos.ApplicationConfig
	BasePath            string `json:"basePath" ,yaml:"basePath"`
	RecordingBaseFolder string
	Middleware          []string         `json:"middleware"`
	PreMiddleware       []string         `json:"pre-middleware"`
	JWTConfig           *GRPCJWTConfig   `json:"jwtConfig"`
	Config              json.RawMessage  `json:"config"`
	Grpc                *GrpcMiddlewares `json:"grpc"`
}

type GRPCJWTConfig struct {
	Key             interface{}            `json:"key"`         //Signing key needed for validating token
	SigningKeys     map[string]interface{} `json:"signingKeys"` //Key map used for validating token.  Can be used in place of a single key
	Certificate     []byte                 `json:"certificate"`
	CertificatePath string                 `json:"certificatePath"` //Path the signing certificate used to validate token.  Can  be used in place of a key
	JWKSUrl         string                 `json:"jwksUrl"`         //URL to JSON Web Key set.  Can be used in place of a Key
	TokenLookup     string                 `json:"tokenLookup"`
	Claims          map[string]interface{} `json:"claims"`
	AuthScheme      string                 `json:"authScheme"`
	ContextKey      string                 `json:"contextKey"`
	SigningMethod   string                 `json:"signingMethod"`
}
