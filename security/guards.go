package security

import (
	"net/http"
	"regexp"
)

const (
	blindGuard = BlindGuard(true)
	rootGuard  = MasterGuard("root")
)

var (
	NoGuard  Guard = blindGuard
	pubGuard       = PathGuard{
		regexp: regexp.MustCompile("/pub/.*"),
		role:   "guest",
	}
	PubGuard  Guard = pubGuard
	RootGuard Guard = rootGuard
)

type BlindGuard bool

func (b BlindGuard) Authorize(*Principle, *http.Request) (bool, error) {
	return bool(b), nil
}

func (b BlindGuard) Matches(*http.Request) bool {
	return true
}

func (b BlindGuard) Authenticate(*http.Request) (*Principle, error) {
	return &Guest, nil
}

type MasterGuard string

func (m MasterGuard) Authorize(p *Principle, _ *http.Request) (bool, error) {
	roles := p.roles
	for _, role := range roles {
		if role == string(m) {
			return true, nil
		}
	}

	return false, nil
}

func (m MasterGuard) Matches(*http.Request) bool {
	return true
}

func (m MasterGuard) Authenticate(r *http.Request) (*Principle, error) {
	principle, err := Authenticate(r)
	return principle, err
}

type PathGuard struct {
	regexp *regexp.Regexp
	role   string
}

func (g PathGuard) PathMatches(s string) bool {
	return g.regexp.MatchString(s)
}

func (g PathGuard) Authorize(_ *Principle, r *http.Request) (bool, error) {
	if !g.PathMatches(r.URL.Path) {
		return false, nil
	}

	return true, nil
}

func (g PathGuard) Matches(r *http.Request) bool {
	return g.PathMatches(r.URL.Path)
}

func (g PathGuard) Authenticate(r *http.Request) (*Principle, error) {
	principle, err := Authenticate(r)
	return principle, err
}
