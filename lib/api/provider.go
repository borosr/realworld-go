package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

type RequestConstraint interface {
	any
}

type ResponseConstraint interface {
	any
}

type Meta struct {
	Headers http.Header
	Params  url.Values
}

type ControllerFuncConstraint[Request RequestConstraint, Response ResponseConstraint] interface {
	ControllerSimpleFunc[Request, Response] | ControllerFunc[Request, Response]
}

type ControllerFunc[Request RequestConstraint, Response ResponseConstraint] func(ctx context.Context, request Request, meta Meta) (Response, error)
type ControllerSimpleFunc[Request RequestConstraint, Response ResponseConstraint] func(ctx context.Context, request Request) (Response, error)

func methodWrapper[Request RequestConstraint, Response ResponseConstraint, Function ControllerFuncConstraint[Request, Response]](path, method string, f Function) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			http.NotFound(w, r)
			return
		}
		ctx, err := processMiddlewares(middlewares[path][r.Method][pre])(w, r)
		if err != nil {
			handleResponse(w, err)
			return
		}
		req, err := processRequest[Request](r.Body)
		if err != nil {
			handleResponse(w, err)
			return
		}

		if _, ok := middlewares[path][r.Method][validate]; ok {
			if v, ok := (interface{})(req).(Validator); ok {
				if err := v.Validate(r.Context()); err != nil {
					handleResponse(w, err)
					return
				}
			}
		}

		var resp Response
		switch ft := (interface{})(f).(type) {
		case ControllerSimpleFunc[Request, Response]:
			resp, err = ft(ctx, req)
		case ControllerFunc[Request, Response]:
			resp, err = ft(ctx, req, Meta{
				Headers: r.Header,
				Params:  r.URL.Query(),
			})
		}
		if err != nil {
			handleResponse(w, err)
			return
		}
		provideResponse[Response](w, resp)

		if _, err := processMiddlewares(middlewares[path][r.Method][post])(w, r); err != nil {
			handleResponse(w, err)
			return
		}
	}
}

func processMiddlewares(mws []Middleware) MiddlewareFunc {
	doneFunc := func(_ http.ResponseWriter, r *http.Request) (context.Context, error) {
		// Done
		return r.Context(), nil
	}
	if len(mws) == 0 {
		return doneFunc
	}
	h := mws[0](doneFunc)
	for i := 1; i < len(mws); i++ {
		h = mws[i](h)
	}
	return h
}

func provideResponse[Response ResponseConstraint](w http.ResponseWriter, resp Response) {
	rawResponse, err := json.Marshal(resp)
	if err != nil {
		handleResponse(w, err)
		return
	}
	log.Println(string(rawResponse))
	w.Header().Set("Content-type", "application/json")
	if _, err := w.Write(rawResponse); err != nil {
		handleResponse(w, err)
		return
	}
}

func processRequest[Request RequestConstraint](body io.Reader) (Request, error) {
	var r Request
	if err := json.NewDecoder(body).Decode(&r); err != nil && !errors.Is(err, io.EOF) {
		log.Printf("decode error: %v", err)
		return r, err
	}
	return r, nil
}

func handleResponse(w http.ResponseWriter, err error) {
	// TODO create custom error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
