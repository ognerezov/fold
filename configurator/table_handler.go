package configurator

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetTableHandlers(route string, mux *goji.Mux, api *openapi.ApiDescription) {
	console.BluePrintln("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		router.ProcessSearch(route, w, r)
	})
	paramLiteral := ":id"
	paramBaseRoute := route + "/"
	if route == "/" {
		paramBaseRoute = "/"
	}
	paramRoute := fmt.Sprintf("%s%s", paramBaseRoute, paramLiteral)

	mux.HandleFunc(pat.Get(paramRoute), func(w http.ResponseWriter, r *http.Request) {
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

	api.Describe(route, openapi.Path{
		"get": openapi.Method{
			Summary: "Get table " + route,
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Raw file",
					Content: map[string]openapi.Content{
						util.ApplicationJson: {
							Schema: openapi.Schema{
								Type: "array",
							},
						},
					},
				},
			},
		},
		"post": openapi.Method{
			Summary: "Post new entity of " + route,
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Raw file",
					Content: map[string]openapi.Content{
						util.ApplicationJson: {
							Schema: openapi.Schema{
								Type: "object",
								Properties: map[string]openapi.Schema{
									"id": {
										Type: "string",
									},
								},
								Required: []string{"id"},
							},
						},
					},
				},
			},
		},
	})
	paramApiRoute := fmt.Sprintf("%s{id}", paramBaseRoute)
	api.Describe(paramApiRoute, openapi.Path{
		"get": openapi.Method{
			Summary: "Get entity by id from " + route,
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Raw file",
					Content: map[string]openapi.Content{
						util.ApplicationJson: {
							Schema: openapi.Schema{
								Type: "array",
							},
						},
					},
				},
			},
		},
		"post": openapi.Method{
			Summary: "Post new entity of " + route,
			Responses: map[string]openapi.Response{
				"200": {
					Description: "Raw file",
					Content: map[string]openapi.Content{
						util.ApplicationJson: {
							Schema: openapi.Schema{
								Type: "object",
								Properties: map[string]openapi.Schema{
									"id": {
										Type: "string",
									},
								},
								Required: []string{"id"},
							},
						},
					},
				},
			},
		},
	})
}
