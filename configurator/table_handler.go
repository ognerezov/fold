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

type CsvSetup struct {
	mux        *goji.Mux
	baseRoute  string
	paramRoute string
}

func (cs CsvSetup) Act() {
	console.CyanPrintln("Registering GET " + cs.paramRoute)
	cs.mux.HandleFunc(pat.Get(cs.paramRoute), func(w http.ResponseWriter, r *http.Request) {
		console.BluePrintln("Incoming GET request to " + cs.paramRoute)
		id := pat.Param(r, "id")
		router.ProcessGet(cs.baseRoute, id, w)
	})

	console.CyanPrintln("Registering PATCH " + cs.paramRoute)
	cs.mux.HandleFunc(pat.Patch(cs.paramRoute), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		decoder := json.NewDecoder(r.Body)
		var record map[string]string
		err := decoder.Decode(&record)
		if err != nil {
			router.ServerError(err, w)
			return
		}
		router.ProcessPatch(cs.baseRoute, id, record, w)
	})

	console.CyanPrintln("Registering DELETE " + cs.paramRoute)
	cs.mux.HandleFunc(pat.Delete(cs.paramRoute), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		router.ProcessDelete(cs.baseRoute, id, w)
	})
}

func SetTableHandlers(route string, mux *goji.Mux, api *openapi.ApiDescription) Action {

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

	console.BluePrintln("Registering POST " + route)
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

	return CsvSetup{mux: mux, baseRoute: route, paramRoute: paramRoute}
}
