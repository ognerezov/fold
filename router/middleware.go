package router

import (
	"fmt"
	"fold/console"
	"net/http"
)

func LogRequest(f http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			console.BluePrintln(fmt.Sprintf("Incoming Request %s: %s", r.Method, r.RequestURI))
			f.ServeHTTP(w, r)
		})
}
