package configurator

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"fold/controls"
	"fold/openapi"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"maps"
	"net/http"
	"strings"
)

const (
	echo       = "echo"
	restart    = "restart"
	googleAuth = "google_auth"
	adaptor    = "adaptor"
	reload     = "reload"
)

type InstructionMap map[string]controls.ControlFactory

var (
	TheInstructions = &InstructionMap{
		echo:    controls.GetEcho,
		adaptor: ConfigureAdaptor,
		reload:  controls.ConfigureReLoader,
	}
)

type ControllerData struct {
	Id         string         `json:"id"`
	Parameters map[string]any `json:"parameters"`
	Method     string         `json:"method" default:"GET"`
	PathParams []string       `json:"path_params"`
	Config     any            `json:"config"`
}

func (c *ControllerData) Controller() (*Controller, bool) {
	ctr, ok := (*TheInstructions)[c.Id]
	if !ok {
		return nil, false
	}
	params := c.Parameters

	if params == nil {
		params = make(map[string]any)
	}

	paramLiterals := util.ParamLiterals(c.PathParams, "/:")

	return &Controller{Id: c.Id,
		Parameters:    params,
		Method:        c.Method,
		Control:       ctr(c.Id, c.Config),
		Config:        c.Config,
		ParamLiterals: paramLiterals}, true
}

type Controller struct {
	Id            string
	Parameters    map[string]any
	Method        string
	Control       *controls.Control
	Config        any
	ParamLiterals string
}

func SetControlHandlers(route string, filePath string, mux *goji.Mux, api *openapi.ApiDescription) {
	var config ControllerData

	err := util.FromJson(filePath, &config)
	if err != nil {
		console.RedPrintln("Error registering route " + route + " : " + err.Error())
		return
	}

	controller, ok := config.Controller()
	if !ok {
		console.RedPrintln("Error registering route " + route + " . Instructions not found")
		return
	}
	paramRoute := route + controller.ParamLiterals
	console.BluePrintln("Registering  " + config.Method + " " + paramRoute)
	parameters, responses := (*controller.Control).Describe()
	api.Path(paramRoute).Method(strings.ToLower(config.Method), openapi.Method{
		Summary:    "Fold action " + config.Id,
		Responses:  responses,
		Parameters: parameters,
	})
	mux.HandleFunc(pat.NewWithMethods(paramRoute, config.Method), func(w http.ResponseWriter, r *http.Request) {
		console.GreenPrintln(fmt.Sprintf("Incoming request to %s: %s", config.Method, route))

		q, _ := util.MapQuery(r)
		data := make(map[string]any)
		maps.Copy(data, controller.Parameters)
		if controller.Method != "GET" && r.Body != nil {
			decoder := json.NewDecoder(r.Body)
			var record map[string]any
			err = decoder.Decode(&record)
			if err != nil {
				console.RedPrintln("Error decoding request body " + err.Error())
			} else {
				maps.Copy(data, record)
			}
		}
		for k, v := range q {
			data[k] = v
		}
		ctr := *controller.Control
		ctr.Do(data, w, r)
	})
}
