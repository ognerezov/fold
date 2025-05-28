package configurator

import (
	"net/http"
	"slices"
)

func AddHeaders(f http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Origin, Access-Control-Allow-Origin, Access-Control-Allow-Headers, X-Requested-With, Access-Control-Request-Method, Access-Control-Request-Headers")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			allowedOrigins := []string{config.AllowOrigin, "http://localhost:3000", "http://localhost:3333"}
			origin := r.Header.Get("Origin")
			if slices.Contains(allowedOrigins, "*") || slices.Contains(allowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			f.ServeHTTP(w, r)
		})
}
