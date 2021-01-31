package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/ksuid"
	"github.com/wepala/weos-controller/core"
	weosmodule "github.com/wepala/weos/module"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

//Handlers container for all handlers

func NewAPIPlugin(e *echo.Echo) *APIPlugin {
	return &APIPlugin{
		HTTPClient:        http.DefaultClient,
		CommandDispatcher: &weosmodule.DefaultDispatcher{},
		e:                 e,
	}
}

type APIPlugin struct {
	HTTPClient        *http.Client
	CommandDispatcher weosmodule.Dispatcher
	Config            *core.APIConfig
	e                 *echo.Echo
}

func (p *APIPlugin) AddConfig(config *core.APIConfig) error {
	p.Config = config
	return nil
}

func (p *APIPlugin) InitModules(mod *weosmodule.WeOSMod) {

}

//Common Middleware

func (p *APIPlugin) RequestID(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return ksuid.New().String()
		},
	})(handlerFunc)
}

func (p *APIPlugin) Static(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Static("/static")(handlerFunc)
}

func (p *APIPlugin) Logger(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Logger()(handlerFunc)
}

func (p *APIPlugin) Recover(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return middleware.Recover()(handlerFunc)
}

func (p *APIPlugin) RequestRecording(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
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

		return nil
	}
}

func (p *APIPlugin) HealthChecker(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
