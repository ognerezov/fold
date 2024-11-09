package configurator

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetTableHandlers(route string, mux *goji.Mux) {
	console.BluePrintln("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		router.ProcessSearch(route, w, r)
	})
	paramLiteral := "/:id"
	if route == "/" {
		paramLiteral = ":id"
	}
	mux.HandleFunc(pat.Get(fmt.Sprintf("%s%s", route, paramLiteral)), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		router.ProcessGet(route, id, w)
	})

	mux.HandleFunc(pat.Patch(fmt.Sprintf("%s%s", route, paramLiteral)), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		decoder := json.NewDecoder(r.Body)
		var record map[string]string
		err := decoder.Decode(&record)
		if err != nil {
			router.ServerError(err, w)
			return
		}
		router.ProcessPatch(route, id, record, w)
	})

	mux.HandleFunc(pat.Post(route), func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var record map[string]string
		err := decoder.Decode(&record)
		if err != nil {
			router.ServerError(err, w)
			return
		}
		router.ProcessPost(route, record, w)
	})

	mux.HandleFunc(pat.Delete(fmt.Sprintf("%s%s", route, paramLiteral)), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		router.ProcessDelete(route, id, w)
	})
}
