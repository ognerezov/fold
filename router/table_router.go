package router

import (
	"errors"
	"fold/mem"
	"net/http"
)

func ProcessSearch(route string, w http.ResponseWriter, r *http.Request) {
	store := *mem.TheStore
	table, ok := store.GetTable(route)
	if !ok {
		err := errors.New("table not found")
		ServerError(err, w)
		return
	}

	query := mem.QueryForTable(r, table)

	data := table.Search(query)

	WriteResponse(data, w)
}

func ProcessGet(route string, id string, w http.ResponseWriter) {
	store := *mem.TheStore
	table, ok := store.GetTable(route)
	if !ok {
		err := errors.New("table not found")
		ServerError(err, w)
		return
	}

	data := table.Get(id, &store)

	WriteResponse(data, w)
}
