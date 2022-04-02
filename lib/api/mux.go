package api

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/borosr/realworld/lib/broken"
)

var (
	mu          = sync.RWMutex{}
	mux         = &Server{}
	middlewares = make(middlewareCollector) // TODO consider remove last map

	pathVariableRegex = regexp.MustCompile("^{([a-zA-Z0-9-_]+)(:(.*)|)}$")
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
	pattern *regexp.Regexp
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
		return nil, broken.Internal("methods not match")
	}

	if strings.Count(urlPath, "/")+1 != len(r.pattern) {
		return nil, broken.Internal("urls not match")
	}
	splitURL := strings.Split(urlPath, "/")

	var matched bool
	for i := 0; i < len(splitURL); i++ {
		if p, ok := r.params[i]; ok {
			if p.pattern != nil {
				matched = p.pattern.MatchString(splitURL[i])
				if matched {
					ctx = SetPathVariable(ctx, p.name, splitURL[i])
				}
			} else {
				ctx = SetPathVariable(ctx, p.name, splitURL[i])
			}
		} else {
			matched = r.pattern[i] == splitURL[i]
		}
	}
	if !matched {
		return nil, broken.Internal("urls not match")
	}

	return ctx, nil
}

func buildParams(splitPattern []string) map[int]routeParam {
	params := make(map[int]routeParam)
	for i, p := range splitPattern {
		if matches := pathVariableRegex.FindStringSubmatch(p); len(matches) > 1 {
			r := routeParam{
				name:    matches[1],
				pattern: nil,
			}
			if len(matches) > 2 && matches[2] != "" {
				r.pattern = regexp.MustCompile(matches[2][1:])
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
	for _, route := range s.routes { // TODO make it faster with tree
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
