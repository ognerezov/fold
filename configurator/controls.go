package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/controls"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"maps"
	"net/http"
)

const (
	echo    = "echo"
	restart = "restart"
)

type InstructionMap map[string]*controls.Control

var (
	TheInstructions = &InstructionMap{
		echo: controls.GetEcho(echo),
	}
)

type ControllerData struct {
	Id         string         `json:"id"`
	Parameters map[string]any `json:"parameters"`
	Method     string         `json:"method" default:"GET"`
}

func (c *ControllerData) Controller() (*Controller, bool) {
	ctr, ok := (*TheInstructions)[c.Id]
	if !ok {
		return nil, false
	}
	fmt.Println(c.Id)
	fmt.Println(*ctr)
	params := c.Parameters

	if params == nil {
		params = make(map[string]any)
	}

	return &Controller{Id: c.Id, Parameters: params, Method: c.Method, Control: ctr}, true
}

type Controller struct {
	Id         string
	Parameters map[string]any
	Method     string
	Control    *controls.Control
}

// TODO set openapi
func SetControlHandlers(route string, filePath string, mux *goji.Mux, api *openapi.ApiDescription) {
	var config ControllerData

	err := util.FromJson(filePath, &config)
	if err != nil {
		console.RedPrintln("Error registering route " + route + " : " + err.Error())
		return
	}
	console.BluePrintln("Registering  " + config.Method + " " + route)

	mux.HandleFunc(pat.NewWithMethods(route, config.Method), func(w http.ResponseWriter, r *http.Request) {
		console.GreenPrintln(fmt.Sprintf("Incoming request to %s: %s", config.Method, route))
		controller, ok := config.Controller()
		if !ok {
			console.RedPrintln("Error registering route " + route + " . Instructions not found")
			router.ServerError(errors.New("controller not found"), w)
			return
		}
		q, _ := util.MapQuery(r)
		data := make(map[string]any)
		maps.Copy(data, controller.Parameters)

		if controller.Method != "GET" && r.Body != nil {
			fmt.Println("Reading body")
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
		res, e := ctr.Do(data)
		if e != nil {
			router.ServerError(e, w)
			return
		}
		router.WriteResponse(res, w)
	})
}
