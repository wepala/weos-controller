package weoscontroller

import (
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/ksuid"
)

//Handlers container for all handlers

type API struct {
	Config *APIConfig
	e      *echo.Echo
}

func (p *API) AddConfig(config *APIConfig) error {
	p.Config = config
	return nil
}

func (p *API) EchoInstance() *echo.Echo {
	return p.e
}

func (p *API) SetEchoInstance(e *echo.Echo) {
	p.e = e
}

//Common Middleware

//Functionality to check claims will be added here
func (a *API) Authenticate(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(a.Config.JWTConfig.Key),
		TokenLookup: a.Config.JWTConfig.TokenLookup,
	})(handlerFunc)
}

func (p *API) RequestID(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return ksuid.New().String()
		},
	})(handlerFunc)
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

func (p *API) HealthChecker(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
