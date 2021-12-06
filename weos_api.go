package weoscontroller

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/wepala/weos"
	weosLogs "github.com/wepala/weos-controller/log"

	"github.com/SermoDigital/jose/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/segmentio/ksuid"
)

//Handlers container for all handlers
const HeaderXAccountID = "X-Account-ID"

type API struct {
	Config      *APIConfig
	e           *echo.Echo
	PathConfigs map[string]*PathConfig
}

func (p *API) AddConfig(config *APIConfig) error {
	p.Config = config
	return nil
}

func (p *API) AddPathConfig(path string, config *PathConfig) error {
	if p.PathConfigs == nil {
		p.PathConfigs = make(map[string]*PathConfig)
	}
	p.PathConfigs[path] = config
	return nil
}

func (p *API) EchoInstance() *echo.Echo {
	return p.e
}

func (p *API) SetEchoInstance(e *echo.Echo) {
	p.e = e
}

//Common Middleware

func (p *API) HTTPSRedirect(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.HTTPSRedirect()(handlerFunc)
}

func (p *API) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*Context)
		req := cc.Request()
		res := cc.Response()
		rid := req.Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = ksuid.New().String()
		}
		res.Header().Set(echo.HeaderXRequestID, rid)

		return next(cc.WithValue(cc, weos.REQUEST_ID, rid))
	}
}

func (p *API) LogLevel(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*Context)
		req := cc.Request()
		logLevel := req.Header.Get(string(weos.LOG_LEVEL))
		if logLevel != "" {
			return next(cc.WithValue(cc, weos.LOG_LEVEL, logLevel))
		}
		return next(c)
	}
}

func (p *API) AccountID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*Context)
		req := cc.Request()
		accountID := req.Header.Get(HeaderXAccountID)
		if accountID != "" {
			return next(cc.WithValue(cc, weos.ACCOUNT_ID, accountID))
		}
		return next(c)
	}
}

func (p *API) UserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if validUser, ok := user.(*jwt.Token); ok {
			cc := c.(*Context)
			claims := validUser.Claims.(jwt.MapClaims)
			return next(cc.WithValue(cc, weos.USER_ID, claims["sub"].(string)))
		}

		return next(c)
	}
}

func (p *API) Static(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Static("/static")(handlerFunc)
}

func (p *API) Logger(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Logger()(handlerFunc)
}

func (p *API) Recover(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Recover()(handlerFunc)
}

func (a *API) getKey(token *jwt.Token) (interface{}, error) {

	keySet, err := jwk.Fetch(context.Background(), a.Config.JWTConfig.JWKSUrl)
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

//Functionality to check claims will be added here
func (a *API) Authenticate(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	var config middleware.JWTConfig
	if a.Config.JWTConfig.JWKSUrl != "" {
		config := middleware.JWTConfig{
			KeyFunc: a.getKey,
		}

		return middleware.JWTWithConfig(config)(handlerFunc)
	}
	if a.Config.JWTConfig.Key != "" {
		config.SigningKey = []byte(a.Config.JWTConfig.Key)
	}
	if len(a.Config.JWTConfig.SigningKeys) > 0 {
		config.SigningKeys = a.Config.JWTConfig.SigningKeys
	}
	if a.Config.JWTConfig.SigningMethod != "" {
		config.SigningMethod = a.Config.JWTConfig.SigningMethod
	}
	if a.Config.JWTConfig.CertificatePath != "" && a.Config.JWTConfig.Certificate == nil {
		bytes, err := ioutil.ReadFile(a.Config.JWTConfig.CertificatePath)
		a.Config.JWTConfig.Certificate = bytes
		if err != nil {
			a.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
		}
	}
	if a.Config.JWTConfig.Certificate != nil {
		if config.SigningMethod == "RS256" || config.SigningMethod == "RS384" || config.SigningMethod == "RS512" {
			publicKey, err := crypto.ParseRSAPublicKeyFromPEM(a.Config.JWTConfig.Certificate)
			if err != nil {
				a.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
			}
			config.SigningKey = publicKey
		} else if config.SigningMethod == "EC256" || config.SigningMethod == "EC384" || config.SigningMethod == "EC512" {
			publicKey, err := crypto.ParseECPublicKeyFromPEM(a.Config.JWTConfig.Certificate)
			if err != nil {
				a.e.Logger.Fatalf("unable to read the jwt certificate, got error '%s'", err)
			}
			config.SigningKey = publicKey
		}
	}
	if config.SigningKey == nil && config.SigningKeys == nil {
		a.e.Logger.Fatalf("no jwt secret was configured.")
	}
	if a.Config.JWTConfig.TokenLookup != "" {
		config.TokenLookup = a.Config.JWTConfig.TokenLookup
	}
	if a.Config.JWTConfig.AuthScheme != "" {
		config.AuthScheme = a.Config.JWTConfig.AuthScheme
	}
	if a.Config.JWTConfig.ContextKey != "" {
		config.ContextKey = a.Config.JWTConfig.ContextKey
	}
	return middleware.JWTWithConfig(config)(handlerFunc)
}

func (p *API) RequestRecording(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := strings.Replace(c.Path(), "/", "_", -1)
		baseFolder := p.Config.RecordingBaseFolder
		if baseFolder == "" {
			baseFolder = "testdata/http"
		}

		p.e.Logger.Infof("Record request to %s", baseFolder+"/"+name+".input.http")

		reqf, err := os.Create(baseFolder + "/" + name + ".input.http")
		if err == nil {
			//record request
			requestBytes, _ := httputil.DumpRequest(c.Request(), true)
			_, err := reqf.Write(requestBytes)
			if err != nil {
				return err
			}
		} else {
			return err
		}

		defer func() {
			reqf.Close()
			if r := recover(); r != nil {
				p.e.Logger.Errorf("Recording failed with errors: %s", r)
			}
		}()

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

func (p *API) ResponseRecording(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := strings.Replace(c.Path(), "/", "_", -1)
		baseFolder := p.Config.RecordingBaseFolder
		if baseFolder == "" {
			baseFolder = "testdata/http"
		}

		responseRecorder := httptest.NewRecorder()
		c.Response().Writer = MultiWriter(c.Response().Writer, responseRecorder)

		if err := next(c); err != nil {
			c.Error(err)
		}

		p.e.Logger.Infof("Record response to %s", baseFolder+"/"+name+".golden.http")
		respf, err := os.Create(baseFolder + "/" + name + ".golden.http")
		if err == nil {
			//record response

			responseBytes, _ := httputil.DumpResponse(responseRecorder.Result(), true)
			_, err = respf.Write(responseBytes)
			if err != nil {
				c.Error(err)
			}

		} else {
			return err
		}

		defer func() {
			respf.Close()
			if r := recover(); r != nil {
				p.e.Logger.Errorf("Recording failed with errors: %s", r)
			}
		}()

		return nil
	}
}

func (p *API) HealthChecker(c echo.Context) error {
	response := &HealthCheckResponse{
		Version: p.Config.Version,
	}
	return c.JSON(http.StatusOK, response)
}

func (p *API) Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			Context: c,
		}
		return next(cc)
	}
}

func (p *API) ZapLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//setting the default logger in the context as zap with the default mode being error
		zapLogger, err := weosLogs.NewZap("error")
		if err != nil {
			p.e.Logger.Errorf("Unexpected error setting the context logger : %s", err)
		}
		zapLogger.SetPrefix("zap")
		c.SetLogger(zapLogger)
		cc := &Context{
			Context: c,
		}
		return next(cc)
	}
}

func (p *API) Initialize() error {
	panic("implement me")
}
