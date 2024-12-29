package configurator

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetTableHandlers(route string, mux *goji.Mux, api *openapi.ApiDescription) {
	console.BluePrintln("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		console.BluePrintln("Incoming GET request to " + route)
		router.ProcessSearch(route, w, r)
	})
	paramLiteral := ":id"
	paramBaseRoute := route + "/"
	if route == "/" {
		paramBaseRoute = "/"
	}
	paramRoute := fmt.Sprintf("%s%s", paramBaseRoute, paramLiteral)
	table, _ := mem.TheStore.GetTable(route)
	schema := table.Schema()
	console.BluePrintln("Registering GET " + paramRoute)
	mux.HandleFunc(pat.Get(paramRoute), func(w http.ResponseWriter, r *http.Request) {
		console.BluePrintln("Incoming GET request to " + paramRoute)
		id := pat.Param(r, "id")
		router.ProcessGet(route, id, w)
	})
	mux.HandleFunc(pat.Patch(paramRoute), func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc(pat.Delete(paramRoute), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		router.ProcessDelete(route, id, w)
	})

	api.Path(route).GetJson(
		"Get table "+route,
		openapi.AnArray,
		false).PostJson(
		"Post new entity of "+route,
		schema,
		openapi.IdObject)

	paramApiRoute := fmt.Sprintf("%s{id}", paramBaseRoute)

	api.Path(paramApiRoute).GetJson(
		"Get entity by id from "+route,
		schema, true).PatchJson(
		"Update entity of "+route,
		schema,
		schema).DeleteById()

}
