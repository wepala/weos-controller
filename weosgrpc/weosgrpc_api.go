package weosgrpc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/SermoDigital/jose/crypto"
	"github.com/golang-jwt/jwt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/lestrrat-go/jwx/jwk"
	weoscontroller "github.com/wepala/weos-controller"
	"google.golang.org/grpc"
)

type GRPCAPI struct {
	Config        *weoscontroller.APIConfig
	c             *context.Context
	PathConfigs   map[string]*weoscontroller.PathConfig
	ServerOptions *weoscontroller.GRPCServerOptions
}

func (p *GRPCAPI) AddConfig(config *weoscontroller.APIConfig) error {
	p.Config = config
	return nil
}

func (p *GRPCAPI) AddPathConfig(path string, config *weoscontroller.PathConfig) error {
	if p.PathConfigs == nil {
		p.PathConfigs = make(map[string]*weoscontroller.PathConfig)
	}
	p.PathConfigs[path] = config
	return nil
}

func (p *GRPCAPI) Context() *context.Context {
	return p.c
}

func (p *GRPCAPI) SetContext(c *context.Context) {
	p.c = c
}

func (p *GRPCAPI) GetStreamMiddleware() grpc.ServerOption {
	return p.ServerOptions.StreamMiddleware
}

func (p *GRPCAPI) GetUnaryMiddleware() grpc.ServerOption {
	return p.ServerOptions.UnaryMiddleware
}

func (p *GRPCAPI) SetAllMiddleware() {
	grpcStream := make([]grpc.StreamServerInterceptor, 2)
	grpcUnary := make([]grpc.UnaryServerInterceptor, 2)
	//TODO call the functions to convert the middleware to the interceptors and append to array
	//call setUnaryMiddleware and setStreamMiddleware with the array

	grpcMiddlewareConfig := p.Config.Grpc.Middlewares

	for _, streamMiddleware := range grpcMiddlewareConfig.Stream.Middleware {
		switch streamMiddleware {
		case "Authenticate":
			grpcStream = append(grpcStream, grpc_auth.StreamServerInterceptor(p.AuthFunc)) //Not sure how to properly pass the auth function into this
		case "Recovery":
			grpcStream = append(grpcStream, grpc_recovery.StreamServerInterceptor())
		}
	}

	for _, UnaryMiddleware := range grpcMiddlewareConfig.Unary.Middleware {
		switch UnaryMiddleware {
		case "Authenticate":
			grpcUnary = append(grpcUnary, grpc_auth.UnaryServerInterceptor(p.AuthFunc)) //Not sure how to properly pass the auth function into this
		case "Recovery":
			grpcUnary = append(grpcUnary, grpc_recovery.UnaryServerInterceptor())
		}
	}

	chainStream := grpc_middleware.ChainStreamServer(grpcStream...)
	p.ServerOptions.StreamMiddleware = grpc.StreamInterceptor(chainStream)

	chainUnary := grpc_middleware.ChainUnaryServer(grpcUnary...)
	p.ServerOptions.UnaryMiddleware = grpc.UnaryInterceptor(chainUnary)
}

func (p *GRPCAPI) getKey(token *jwt.Token) (interface{}, error) {

	keySet, err := jwk.Fetch(context.Background(), p.Config.JWTConfig.JWKSUrl)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have a key ID in the kid field")
	}

	key, found := keySet.LookupKeyID(keyID)

	if !found {
		return nil, fmt.Errorf("unable to find key %q", keyID)
	}

	var pubkey interface{}
	if err := key.Raw(&pubkey); err != nil {
		return nil, fmt.Errorf("unable to get the public key. error: %s", err.Error())
	}

	return pubkey, nil
}

func Authenticate(ctx context.Context) (context.Context, error) {
	//Remove all middleware. usage as this is related to echo. An alternative is required
	var config weoscontroller.JWTConfig
	if p.Config.JWTConfig.JWKSUrl != "" {
		config := weoscontroller.JWTConfig{
			KeyFunc: p.getKey, //KeyFunc returns a interface, maybe use key in jwtconfig - this takes string
		}

		context := context.WithValue(*ctx, "grpcServerOptions", config)
		return &context
	}
	//RD: Not sure about this one
	if p.Config.JWTConfig.Key != "" {
		config.SigningKey = []byte(p.Config.JWTConfig.Key)
	}
	if len(p.Config.JWTConfig.SigningKeys) > 0 {
		config.SigningKeys = p.Config.JWTConfig.SigningKeys
	}
	if p.Config.JWTConfig.SigningMethod != "" {
		config.SigningMethod = p.Config.JWTConfig.SigningMethod
	}
	if p.Config.JWTConfig.CertificatePath != "" && p.Config.JWTConfig.Certificate == nil {
		bytes, err := ioutil.ReadFile(p.Config.JWTConfig.CertificatePath)
		p.Config.JWTConfig.Certificate = bytes
		if err != nil {
			//p.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
		}
	}
	if p.Config.JWTConfig.Certificate != nil {
		if config.SigningMethod == "RS256" || config.SigningMethod == "RS384" || config.SigningMethod == "RS512" {
			publicKey, err := crypto.ParseRSAPublicKeyFromPEM(p.Config.JWTConfig.Certificate)
			if err != nil {
				//p.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
			}

			// RD: Not sure about this one
			config.SigningKey = publicKey
		} else if config.SigningMethod == "EC256" || config.SigningMethod == "EC384" || config.SigningMethod == "EC512" {
			publicKey, err := crypto.ParseECPublicKeyFromPEM(p.Config.JWTConfig.Certificate)
			if err != nil {
				//a.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
			}
			// RD: Not sure about this one
			config.SigningKey = publicKey
		}
	}

	// RD: Not sure about this one
	if config.SigningKey == nil && config.SigningKeys == nil {
		//p.e.Logger.Fatalf("no jwt secret was configured.")
	}
	if p.Config.JWTConfig.TokenLookup != "" {
		config.TokenLookup = p.Config.JWTConfig.TokenLookup
	}
	if p.Config.JWTConfig.AuthScheme != "" {
		config.AuthScheme = p.Config.JWTConfig.AuthScheme
	}
	if p.Config.JWTConfig.ContextKey != "" {
		config.ContextKey = p.Config.JWTConfig.ContextKey
	}

	context := context.WithValue(*ctx, "grpcServerOptions", config)
	return &context
}

func (p *GRPCAPI) AuthFunc(ctx context.Context) (context.Context, error) { return nil, nil }
