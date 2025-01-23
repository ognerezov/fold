package controls

import (
	"errors"
	"fold/console"
	"fold/openapi"
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

func (id EchoControl) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return make([]openapi.Parameter, 0), map[string]openapi.Response{
		"200": {
			Description: "Raw data",
			Content: map[string]openapi.Content{
				openapi.ApplicationJson: {
					Schema: openapi.Binary,
				},
			},
		},
	}
}

func GetEcho(id string, _ any) *Control {
	var ctr Control
	ctr = EchoControl(id)
	return &ctr
}
