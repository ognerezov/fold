package configurator

import (
	"fmt"
	"fold/console"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetCSVHandlers(route string, mux *goji.Mux) {
	console.BluePrintln("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		err := router.ProcessSearch(route, w, r)
		if err != nil {
			console.RedPrint(err.Error())
		}
	})
	paramLiteral := "/:id"
	if route == "/" {
		paramLiteral = ":id"
	}
	mux.HandleFunc(pat.Get(fmt.Sprintf("%s%s", route, paramLiteral)), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		err := router.ProcessGet(route, id, w)
		if err != nil {
			console.RedPrint(err.Error())
		}
	})
}
