package security

import (
	"errors"
	"fmt"
	"fold/util"
	"net/http"
	"regexp"
)

type Rule struct {
	Id        int    `json:"id" validate:"required"`
	Authority string `json:"security_authorities_id" validate:"required"`
	Audience  string `json:"security_audience_id" validate:"required"`
	Filter    string `json:"filter"`
	Action    string `json:"action"`
}

func (r Rule) PathMatches(s string) bool {
	var path string
	if r.Audience == "*" {
		path = "^/.*"
	} else if r.Audience == "/" {
		// just root page or address without path
		path = "^/?$"
	} else {
		path = fmt.Sprintf("^/%s.*", r.Audience)
	}
	re := regexp.MustCompile(path)
	return re.MatchString(s)
}

func (r Rule) Authorize(p *Principle, req *http.Request) (bool, error) {
	roleFound := false
	for _, role := range p.Roles {
		if role == r.Authority {
			roleFound = true
		}
	}
	if !roleFound {
		return false, errors.New("not authorized for this action")
	}
	if r.Filter == "" {
		return true, nil
	}
	pathParam, queryParams := util.GetUrlParams(r.Filter, req)

	allParams := make([]string, 0)
	if pathParam != "" {
		allParams = append(allParams, pathParam)
	}
	if queryParams != nil {
		allParams = append(allParams, queryParams...)
	}
	for _, param := range allParams {
		if param == "" {
			continue
		}
		if param != p.Id {
			return false, errors.New("you don't have access to the requested resource")
		}
	}
	return true, nil
}

func (r Rule) Matches(req *http.Request) bool {
	pathMatches := r.PathMatches(req.URL.Path)
	if !pathMatches {
		return false
	}
	if r.Filter != "" {
		param, queryParams := util.GetUrlParams(r.Filter, req)
		if param == "" && queryParams == nil {
			return false
		}
	}

	return r.ActionMatches(req.Method)
}

func (r Rule) Authenticate(req *http.Request) (*Principle, error) {
	if r.Authority == "guest" {
		return &Guest, nil
	}
	principle, err := Authenticate(req)
	return principle, err
}

func (r Rule) ActionMatches(method string) bool {
	if r.Action == "" || r.Action == "*" {
		return true
	}

	if r.Action == method {
		return true
	}

	if r.Action == "read" {
		return method == http.MethodHead ||
			method == http.MethodOptions ||
			method == http.MethodConnect ||
			method == http.MethodTrace ||
			method == http.MethodGet
	}

	if r.Action == "write" || r.Action == "update" {
		return method == http.MethodPut || method == http.MethodPatch || method == http.MethodPost || method == http.MethodDelete
	}

	return false
}

func (r Rule) IsPublic() bool {
	return r.Authority == "guest"
}

func (r Rule) AppliesTo(path string, method string) bool {
	return r.ActionMatches(method) && r.PathMatches(path)
}
