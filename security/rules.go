package security

import (
	"fmt"
	"fold/util"
	"goji.io/pat"
	"net/http"
	"regexp"
	"strings"
)

type Rule struct {
	Id        string `json:"id" validate:"required"`
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
		path = fmt.Sprintf("/%s.*", s)
	}
	re := regexp.MustCompile(path)
	return re.MatchString(s)
}

func (r Rule) Authorize(p *Principle, _ *http.Request) (bool, error) {
	for _, role := range p.Roles {
		fmt.Printf("check %s against %s", role, r.Authority)
		fmt.Println("")
		if role == r.Authority {
			return true, nil
		}
	}
	return false, nil
}

func (r Rule) Matches(req *http.Request) bool {
	pathMatches := r.PathMatches(req.URL.Path)
	if !pathMatches {
		return false
	}
	if r.Filter != "" {
		param := pat.Param(req, r.Filter)
		q, err := util.MapQuery(req)
		var queryParam []string
		if err != nil {
			queryParam = q[r.Filter]
		}
		if param == "" && queryParam == nil {
			return false
		}
	}
	if r.Action != "" {
		if !strings.EqualFold(r.Action, req.Method) {
			return false
		}
	}
	return true
}

func (r Rule) Authenticate(req *http.Request) (*Principle, error) {
	principle, err := Authenticate(req)
	return principle, err
}
