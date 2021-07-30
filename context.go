package weoscontroller

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

type Context struct {
	echo.Context
	requestContext context.Context
}

func (c *Context) WithValue(parent *Context, key, val interface{}) *Context {
	if parent.requestContext != nil {
		parent.requestContext = context.WithValue(parent.requestContext, key, val)
	} else {
		parent.requestContext = context.WithValue(context.TODO(), key, val)
	}
	return parent
}

func (c *Context) RequestContext() context.Context {
	return c.requestContext
}
