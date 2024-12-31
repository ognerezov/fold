package controls

import "net/http"

type Control interface {
	Do(map[string]any, http.ResponseWriter, *http.Request)
}

type ControlFactory func(string, any) *Control
