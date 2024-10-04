package configurator

import (
	"encoding/json"
	"fmt"
	"fold/mem"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetCSVHandlers(route string, mux *goji.Mux) {
	fmt.Println("Registering GET " + route)
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming Request GET: " + route)
		store := *mem.TheStore
		data := store.All(route)
		h, err := json.Marshal(data)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			return
		}
		w.Write(h)
	})
	paramLiteral := "/:id"
	if route == "/" {
		paramLiteral = ":id"
	}
	mux.HandleFunc(pat.Get(fmt.Sprintf("%s%s", route, paramLiteral)), func(w http.ResponseWriter, r *http.Request) {
		id := pat.Param(r, "id")
		fmt.Println(fmt.Sprintf("Incoming Request GET: %s/:%s", route, id))
		store := *mem.TheStore
		data := store.Get(route, id)
		h, err := json.Marshal(data)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(500)
			fmt.Fprintln(w, err)
			return
		}
		w.Write(h)
	})
}
