package server

import (
	"net/http"
	"strings"
)

// corsMiddleware returns a middleware that sets CORS headers allowing only the specified methods.
func corsMiddleware(allowedMethods ...string) func(http.Handler) http.Handler {
	methods := strings.Join(allowedMethods, ", ")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
