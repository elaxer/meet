package middleware

import (
	"meet/internal/pkg/api"
	"net/http"
	"strconv"
)

const maxUploadSize = 10 * 1024 * 1024

func ContentLength(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cl := r.Header.Get("Content-Length")
		if cl == "" {
			next.ServeHTTP(w, r)

			return
		}

		len, err := strconv.Atoi(cl)
		if err != nil {
			api.ResponseError(w, err, http.StatusBadRequest)

			return
		}

		if len > maxUploadSize {
			api.ResponseEmpty(w, http.StatusRequestEntityTooLarge)

			return
		}

		next.ServeHTTP(w, r)
	})
}
