package weosgrpc

import (
	"context"

	weoscontroller "github.com/wepala/weos-controller"
)

type GRPCAPI struct {
	Config      *weoscontroller.APIConfig
	c           *context.Context
	PathConfigs map[string]*weoscontroller.PathConfig
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

/*
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

func (p *GRPCAPI) Authenticate(ctx context.Context) context.Context {
	//Remove all middleware. usage as this is related to echo. An alternative is required
	var config weoscontroller.JWTConfig
	if p.Config.JWTConfig.JWKSUrl != "" {
		config := weoscontroller.JWTConfig{
			KeyFunc: p.getKey,
		}

		return ctx.WithValue("JWTConfig", config)
	}
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
			config.SigningKey = publicKey
		} else if config.SigningMethod == "EC256" || config.SigningMethod == "EC384" || config.SigningMethod == "EC512" {
			publicKey, err := crypto.ParseECPublicKeyFromPEM(p.Config.JWTConfig.Certificate)
			if err != nil {
				//a.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
			}
			config.SigningKey = publicKey
		}
	}
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
	return middleware.JWTWithConfig(config)(handlerFunc)

}


func (p *GRPCAPI) Context(ctx) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			Context: c,
		}
		return next(cc)
	}
}*/
