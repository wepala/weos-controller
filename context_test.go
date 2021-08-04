package weoscontroller_test

import (
	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
	"testing"
)

func TestContext_WithValue(t *testing.T) {
	t.Run("adding variable to context", func(t *testing.T) {
		e := echo.New()
		parentContext := &weoscontroller.Context{
			Context: e.AcquireContext(),
		}

		context := parentContext.WithValue(parentContext, "test", "ing")
		if context.RequestContext().Value("test") != "ing" {
			t.Errorf("expected the context to have key '%s' with value '%s', got '%s'", "test", "ing", context.RequestContext().Value("test"))
		}
	})
}
