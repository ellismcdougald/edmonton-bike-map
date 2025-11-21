package middleware

import (
	"net/http"
	"strings"
)

// CorsMiddleware returns a middleware that sets CORS headers allowing only the specified methods.
func CorsMiddleware(allowedMethods ...string) func(http.Handler) http.Handler {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:4173",
		"http://localhost:5173",
		"https://edmonton-bike-map-frontend.vercel.app",
		"https://www.edmontonbikemap.xyz",
	}

	methods := strings.Join(allowedMethods, ", ")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, o := range allowedOrigins {
				if o == origin {
					w.Header().Set("Access-Control-Allow-Origin", o)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
