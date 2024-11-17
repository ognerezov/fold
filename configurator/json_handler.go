package configurator

import (
	"encoding/json"
	"fold/console"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetJsonHandlers(route string, mux *goji.Mux, api *openapi.ApiDescription) {
	console.BluePrintln("Registering GET " + route)
	noSql, _ := (*mem.TheStore).NoSql(route)
	schema := noSql.Schema()
	entitySchema := noSql.Entity()

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
	patchQuery := openapi.AnObject
	if noSql.IsCollection() {
		patchQuery = entitySchema
	}
	api.Path(route).GetJson(
		"Get or query json",
		entitySchema,
		true).PatchQuery(
		"Update json data",
		patchQuery,
		schema,
		schema).PostJson(
		"Override document or add document to collection",
		openapi.AnObject,
		openapi.AnObject).DeleteQuery(
		"Delete objects from collection or clear document",
		entitySchema,
		openapi.AnObject)
}
