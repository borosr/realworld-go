package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/auth"
	"github.com/borosr/realworld/lib/broken"
)

const tokenPrefix = "Token "

var (
	ErrNotAuthenticated = broken.Internal("not authenticated")
)

var _ api.Middleware = TokenAuthentication

func TokenAuthentication(next api.MiddlewareFunc) api.MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		token := r.Header.Get("Authorization")
		if token == "" {
			return r.Context(), ErrNotAuthenticated
		}

		if !strings.HasPrefix(token, tokenPrefix) {
			return r.Context(), ErrNotAuthenticated
		}

		token = strings.ReplaceAll(token, tokenPrefix, "")
		if token == "" {
			return r.Context(), ErrNotAuthenticated
		}

		claims, err := auth.Verify(token)
		if err != nil {
			return r.Context(), err
		}

		return context.WithValue(r.Context(), "email", claims["email"]), nil
	}
}
