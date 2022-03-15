package api

import (
	"context"
	"net/http"
)

const (
	pre = iota + 1
	post
	validate
)

type Validator interface {
	Validate(ctx context.Context) error
}

type Middleware func(next MiddlewareFunc) MiddlewareFunc

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request) (context.Context, error)

func (e *endpoint) PreProcess(mws ...Middleware) *endpoint {
	methods := e.buildMiddlewareEnvironment()
	for i := range mws {
		methods[pre] = append(methods[pre], mws[i])
	}

	return e
}

func (e *endpoint) PostProcess(mws ...Middleware) *endpoint {
	methods := e.buildMiddlewareEnvironment()
	for i := range mws {
		methods[post] = append(methods[post], mws[i])
	}

	return e
}

func (e *endpoint) Validated() *endpoint {
	methods := e.buildMiddlewareEnvironment()
	methods[validate] = make([]Middleware, 0, 0)
	return e
}

func (e *endpoint) buildMiddlewareEnvironment() map[int][]Middleware {
	paths := middlewares[e.path]
	if paths == nil {
		paths = make(map[string]map[int][]Middleware)
		middlewares[e.path] = paths
	}
	methods := paths[e.method]
	if methods == nil {
		methods = make(map[int][]Middleware)
		middlewares[e.path][e.method] = methods
	}
	return methods
}
