package configurator

import "fold/openapi"

type Endpoint struct {
	Path       string              `json:"path"`
	Method     string              `json:"method"`
	Parameters []openapi.Parameter `json:"parameters"`
	Label      string              `json:"label"`
	Public     bool                `json:"public"`
}

type Endpoints map[string]*Endpoint
