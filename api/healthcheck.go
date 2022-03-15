package api

import (
	"context"
	"fmt"
	"go/types"
	"net/http"

	"github.com/borosr/realworld/lib/api"
)

var _ api.ControllerSimpleFunc[types.Nil, HealthCheckResponse] = HealthCheck

func init() {
	api.Register[types.Nil, HealthCheckResponse, api.ControllerSimpleFunc[types.Nil, HealthCheckResponse]]("/hc", http.MethodGet, HealthCheck).
		PreProcess(func(next api.MiddlewareFunc) api.MiddlewareFunc {
			return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
				fmt.Println("do stuff")
				return next(w, r)
			}
		}).
		PostProcess(func(next api.MiddlewareFunc) api.MiddlewareFunc {
			return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
				fmt.Println("after things...")
				return next(w, r)
			}
		})
}

type HealthCheckResponse struct {
	Msg string `json:"msg"`
}

func HealthCheck(_ context.Context, _ types.Nil) (HealthCheckResponse, error) {
	return HealthCheckResponse{Msg: "ok"}, nil
}
