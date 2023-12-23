package middleware

import (
	"context"
	"meet/internal/app/server"
	"meet/internal/app/service"
	"net/http"
	"strings"
)

const (
	authMethodBearer = "Bearer"
)

type authMiddleware struct {
	authService *service.AuthService
}

func newAuthMiddleware(authService *service.AuthService) *authMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

func (am *authMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		if a == "" {
			server.ResponseEmpty(w, http.StatusUnauthorized)

			return
		}
		values := strings.Split(a, " ")
		method, tokenString := values[0], values[1]
		if method != authMethodBearer {
			server.ResponseEmpty(w, http.StatusUnauthorized)

			return
		}

		u, err := am.authService.Authorize(tokenString)
		if err != nil {
			server.ResponseError(w, err, http.StatusUnauthorized)

			return
		}

		if u.IsBlocked {
			server.ResponseEmpty(w, http.StatusForbidden)

			return
		}

		ctx := context.WithValue(r.Context(), server.CtxKeyUser, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
