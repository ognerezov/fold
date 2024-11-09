package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"net/http"
)

func WriteResponse(data any, w http.ResponseWriter) {
	h, err := json.Marshal(data)
	if err != nil {
		ServerError(err, w)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, e := w.Write(h)
	if e != nil {
		ServerError(e, w)
	}
}

func ServerError(err error, w http.ResponseWriter) {
	ReturnError(err, 500, w)
}

func NotFound(w http.ResponseWriter) {
	ReturnError(errors.New("not found"), 404, w)
}

func ReturnError(err error, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, e := fmt.Fprintln(w, err)

	if e != nil {
		console.RedPrintln("Error writing error response " + e.Error())
	}
}
