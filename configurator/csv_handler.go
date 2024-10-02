package configurator

import (
	"fmt"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetCSVHandlers(route string, filename string, mux *goji.Mux) {
	fmt.Println("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Find file, %s!", filename)
	})
}
