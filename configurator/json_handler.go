package configurator

import (
	"encoding/json"
	"fold/console"
	"fold/mem"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetJsonHandlers(route string, mux *goji.Mux) {
	console.BluePrintln("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		q, _ := util.MapQuery(r)
		router.WriteResponse(n.RawSearch(&q), w)
	})
	console.BluePrintln("Registering PATCH " + route)
	mux.HandleFunc(pat.Patch(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		q, _ := util.MapQuery(r)
		decoder := json.NewDecoder(r.Body)
		var record map[string]any
		err := decoder.Decode(&record)
		if err != nil {
			router.ServerError(err, w)
			return
		}
		router.WriteResponse(n.Patch(&q, &record), w)
	})
	console.BluePrintln("Registering POST " + route)
	mux.HandleFunc(pat.Post(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		decoder := json.NewDecoder(r.Body)
		var record map[string]any
		err := decoder.Decode(&record)
		if err != nil {
			router.ServerError(err, w)
			return
		}
		router.WriteResponse(n.Post(&record), w)
	})
	console.BluePrintln("Registering DELETE " + route)
	mux.HandleFunc(pat.Delete(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		q, _ := util.MapQuery(r)
		router.WriteResponse(n.Delete(&q), w)
	})
}
