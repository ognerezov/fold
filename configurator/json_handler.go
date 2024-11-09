package configurator

import (
	"encoding/json"
	"fold/mem"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetJsonHandlers(route string, mux *goji.Mux) {
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		q, _ := util.MapQuery(r)
		router.WriteResponse(n.RawSearch(&q), w)
	})

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

	mux.HandleFunc(pat.Delete(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		q, _ := util.MapQuery(r)
		router.WriteResponse(n.Delete(&q), w)
	})
}
