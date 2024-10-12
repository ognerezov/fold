package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/mem"
	"net/http"
)

func ProcessSearch(route string, w http.ResponseWriter, r *http.Request) error {
	store := *mem.TheStore
	table, ok := store.GetTable(route)
	if !ok {
		err := errors.New("table not found")
		return WriteError(err, w)
	}

	query := mem.QueryForTable(r, table)

	data := table.Search(query)

	return WriteResponse(data, w)
}

func ProcessGet(route string, id string, w http.ResponseWriter) error {
	store := *mem.TheStore
	table, ok := store.GetTable(route)
	if !ok {
		err := errors.New("table not found")
		return WriteError(err, w)
	}

	data := table.Get(id, &store)

	return WriteResponse(data, w)
}

func WriteResponse(data any, w http.ResponseWriter) error {
	h, err := json.Marshal(data)
	if err != nil {
		return WriteError(err, w)
	}
	_, e := w.Write(h)
	return e
}

func WriteError(err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(500)
	_, e := fmt.Fprintln(w, err)
	return e
}
