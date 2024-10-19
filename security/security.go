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
	console.BluePrintln("Registering POST /login")
	mux.HandleFunc(pat.Post("/login"), func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var login api.LoginRequest
		err := decoder.Decode(&login)
		if err != nil {
			_ = router.WriteError(err, w)
			return
		}
		err = login.Validate()
		if err != nil {
			_ = router.WriteError(err, w)
			return
		}

		console.GreenPrintln(fmt.Sprintf("Processing login request for: %s", login.Username))
		principle, err := FromTable(login.Username, login.Password)

		if err != nil {
			_ = router.WriteError(err, w)
			return
		}
		token, err := principle.BearerToken()

		if err != nil {
			_ = router.WriteError(err, w)
			return
		}

		_ = router.WriteResponse(api.LoginResponse{
			Token: token,
		}, w)
	})
}
