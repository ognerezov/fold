package security

import (
	"fmt"
	"fold/console"
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
		path = "/.*"
	} else {
		path = fmt.Sprintf("/%s.*", r.Audience)
	}
	re := regexp.MustCompile(path)
	return re.MatchString(s)
}

func (r Rule) Authorize(p *Principle, _ *http.Request) (bool, error) {
	for _, role := range p.Roles {
		if role == r.Authority {
			return true, nil
		}
	}
	return false, nil
}

func (r Rule) Matches(req *http.Request) bool {
	console.MagentaPrintln("Authenticating " + req.URL.Path)
	pathMatches := r.PathMatches(req.URL.Path)
	if !pathMatches {
		return false
	}
	if r.Filter != "" {
		var param string
		util.PathParamValue(req, r.Filter, &param)
		fmt.Println(param)
		q, err := util.MapQuery(req)
		var queryParam []string
		if err != nil {
			queryParam = q[r.Filter]
		}
		if param == "" && queryParam == nil {
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
