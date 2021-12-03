package weosgrpc

import (
	"golang.org/x/net/context"
)

type Context struct {
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
