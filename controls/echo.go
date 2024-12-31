package controls

import (
	"errors"
	"fold/console"
	"fold/router"
	"net/http"
)

type EchoControl string

func (id EchoControl) Do(data map[string]any, w http.ResponseWriter, _ *http.Request) {
	if data == nil {
		console.RedPrintln("request is empty")
		router.BadRequest(errors.New("request is empty"), w)
		return
	}

	router.WriteResponse(data, w)
}

func (id EchoControl) ConfigureControl(_ any) error {
	return nil
}

func GetEcho(id string) *Control {
	var ctr Control
	ctr = EchoControl(id)
	return &ctr
}
