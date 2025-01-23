package controls

import (
	"fold/api"
	"fold/console"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	"fold/util"
	"net/http"
)

type ReLoaderConfig struct {
	Path string `json:"path"`
}

type ReLoader string

func (r ReLoader) Do(_ map[string]any, w http.ResponseWriter, _ *http.Request) {
	err := mem.TheStore.Refresh(string(r))
	if err != nil {
		console.RedPrintln(err.Error())
		router.ServerError(err, w)
	}

	router.WriteResponse(api.Ok(), w)
}

func (r ReLoader) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return nil, openapi.StatusResponse
}

func ConfigureReLoader(_ string, config any) *Control {
	var a ReLoaderConfig
	err := util.Restructure(config, &a)
	if err != nil {
		panic(err)
	}
	var res Control
	res = ReLoader(a.Path)
	return &res
}
