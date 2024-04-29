package middleware

import (
	"context"
	"meet/internal/pkg/api"
	"meet/internal/pkg/app"
	"meet/internal/pkg/app/service"
	"net/http"
	"strings"
)

const authMethodBearer = "Bearer"

type AuthorizeMiddleware interface {
	Authorize(next http.Handler) http.Handler
}

type authorizeMiddleware struct {
	authService service.AuthService
}

func NewAuthorizeMiddleware(authService service.AuthService) AuthorizeMiddleware {
	return &authorizeMiddleware{authService}
}

func (am *authorizeMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		if a == "" {
			api.ResponseEmpty(w, http.StatusUnauthorized)

			return
		}
		values := strings.Split(a, " ")
		method, tokenString := values[0], values[1]
		if method != authMethodBearer {
			api.ResponseEmpty(w, http.StatusUnauthorized)

			return
		}

		u, err := am.authService.Authorize(tokenString)
		if err != nil {
			api.ResponseError(w, err, http.StatusUnauthorized)

			return
		}

		if u.IsBlocked {
			api.ResponseEmpty(w, http.StatusForbidden)

			return
		}

		ctx := context.WithValue(r.Context(), app.CtxKeyUser, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
