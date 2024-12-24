package security

import (
	"encoding/json"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func SetAuthHandlers(mux *goji.Mux) {
	console.BluePrintln("Registering security root POST /login (User defined root would be overridden)")
	mux.HandleFunc(pat.Post("/login"), func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var login api.LoginRequest
		err := decoder.Decode(&login)
		if err != nil {
			router.ReturnError(err, 400, w)
			return
		}
		err = login.Validate()
		if err != nil {
			router.ReturnError(err, 400, w)
			return
		}

		console.GreenPrintln(fmt.Sprintf("Processing login request for: %s", login.Username))
		principle, err := FromTable(login.Username, login.Password)

		if err != nil {
			router.ReturnError(err, 401, w)
			return
		}
		token, err := principle.BearerToken()

		if err != nil {
			router.ServerError(err, w)
			return
		}

		router.WriteResponse(api.LoginResponse{
			Token: token,
		}, w)
	})

	console.BluePrintln("Registering security root GET /me (User defined root would be overridden)")
	mux.HandleFunc(pat.Get("/me"), func(w http.ResponseWriter, r *http.Request) {
		principle := FromRequest(r)
		if principle == nil {
			router.ReturnError(fmt.Errorf("principle not found"), 401, w)
			return
		}
		router.WriteResponse(*principle, w)
	})
}
