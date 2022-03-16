package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	mu          = sync.RWMutex{}
	mux         = &Server{}
	middlewares = make(middlewareCollector) // TODO consider remove last map
)

type middlewareCollector map[string]map[string]map[int][]Middleware

type endpoint struct {
	path   string
	method string
}

func Register[Request RequestConstraint, Response ResponseConstraint, Function ControllerFuncConstraint[Request, Response]](path string, method string, handler Function) *endpoint {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("%s %s", method, path)
	mux.HandleFunc(path, method, methodWrapper[Request, Response, Function](path, method, handler))
	return &endpoint{
		path:   path,
		method: method,
	}
}

type route struct {
	pattern []string
	method  string
	params  map[int]routeParam
	handler http.Handler
}

type routeParam struct {
	name    string
	pattern string
}

type Server struct {
	mu     sync.RWMutex
	routes []*route
}

func (s *Server) Handler(pattern, method string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	splitPattern := strings.Split(pattern, "/")
	s.routes = append(s.routes, &route{pattern: splitPattern, method: method, params: buildParams(splitPattern), handler: handler})
}

func (r route) parse(ctx context.Context, urlPath, method string) (context.Context, error) {
	if r.method != method {
		return nil, errors.New("methods not match")
	}

	splitURL := strings.Split(urlPath, "/")
	if len(splitURL) != len(r.pattern) {
		return nil, errors.New("urls not match")
	}

	var matched bool
	for i := 0; i < len(splitURL); i++ {
		if p, ok := r.params[i]; ok {
			if p.pattern != "" {
				matched, _ = regexp.MatchString(p.pattern, splitURL[i])
				if matched {
					ctx = context.WithValue(ctx, ctxPathVariablePrefix+p.name, splitURL[i])
				}
			} else {
				ctx = context.WithValue(ctx, ctxPathVariablePrefix+p.name, splitURL[i])
			}
		} else {
			matched = r.pattern[i] == splitURL[i]
		}
	}
	if !matched {
		return nil, errors.New("urls not match")
	}

	return ctx, nil
}

func buildParams(splitPattern []string) map[int]routeParam {
	params := make(map[int]routeParam)
	for i, p := range splitPattern {
		if matches := regexp.MustCompile("^{([a-zA-Z0-9-_]+)(:(.*)|)}$").FindStringSubmatch(p); len(matches) > 1 {
			r := routeParam{
				name:    matches[1],
				pattern: "",
			}
			if len(matches) > 2 && matches[2] != "" {
				r.pattern = matches[2][1:]
			}
			params[i] = r
		}
	}
	return params
}

func (s *Server) HandleFunc(pattern, method string, handler func(http.ResponseWriter, *http.Request)) {
	s.Handler(pattern, method, http.HandlerFunc(handler))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range s.routes {
		ctx, err := route.parse(r.Context(), r.URL.Path, r.Method)
		if err == nil {
			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}

	http.NotFound(w, r)
}

func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, mux)
}

func ListenAndServeTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, mux)
}
