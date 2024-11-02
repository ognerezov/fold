package configurator

import (
	"fold/mem"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetJsonHandlers(route string, noSql *mem.NoSql, mux *goji.Mux) {
	mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
		n, ok := (*mem.TheStore).NoSql(route)
		if !ok {
			router.NotFound(w)
		}
		router.WriteResponse(n.Val(), w)
	})
}
