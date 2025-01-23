package controls

import (
	"fold/openapi"
	"net/http"
)

type Control interface {
	Do(map[string]any, http.ResponseWriter, *http.Request)
	Describe() ([]openapi.Parameter, map[string]openapi.Response)
}

type ControlFactory func(string, any) *Control
