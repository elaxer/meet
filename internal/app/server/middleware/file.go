package middleware

import (
	"meet/internal/app/server"
	"net/http"
)

func FileSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-type")
		if ct != "multipart/form-data" {
			next.ServeHTTP(w, r)

			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			server.ResponseError(w, err, http.StatusRequestEntityTooLarge)

			return
		}

		next.ServeHTTP(w, r)
	})
}
