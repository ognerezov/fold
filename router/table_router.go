package router

import (
	"errors"
	"fold/api"
	"fold/mem"
	"net/http"
)

func getTable(route string) (*mem.Table, error) {
	table, ok := mem.TheStore.GetTable(route)
	if !ok {
		return nil, errors.New("table not found")
	}
	return table, nil
}

func ProcessSearch(route string, w http.ResponseWriter, r *http.Request) {
	table, err := getTable(route)
	if err != nil {
		ReturnError(err, 404, w)
		return
	}

	query := mem.QueryForTable(r, table)

	data := table.Search(query)

	WriteResponse(data, w)
}

func ProcessGet(route string, id string, w http.ResponseWriter) {
	table, err := getTable(route)
	if err != nil {
		ReturnError(err, 404, w)
		return
	}

	data := table.Get(id, mem.TheStore)

	WriteResponse(data, w)
}

func ProcessPatch(route string, id string, record map[string]string, w http.ResponseWriter) {
	table, err := getTable(route)
	if err != nil {
		ReturnError(err, 404, w)
		return
	}

	data, err := table.Update(id, record, mem.TheStore)

	if err != nil {
		ReturnError(err, 404, w)
		return
	}

	WriteResponse(data, w)
}

func ProcessPost(route string, record map[string]string, w http.ResponseWriter) {
	table, err := getTable(route)
	if err != nil {
		ReturnError(err, 404, w)
		return
	}
	id, err := table.Insert(record)
	if err != nil {
		ReturnError(err, 409, w)
		return
	}

	WriteResponse(api.IdResponse{
		Id: id,
	}, w)
}

func ProcessDelete(route string, id string, w http.ResponseWriter) {
	table, err := getTable(route)
	if err != nil {
		ReturnError(err, 404, w)
		return
	}

	count := table.DeleteById(id)

	WriteResponse(api.CountResponse{Count: count}, w)
}
