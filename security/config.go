package security

import (
	"errors"
	"fold/console"
	"fold/router"
	"net/http"
)

var (
	Public      = &Config{guards: make(Guards, 0), lastGuard: &NoGuard}
	MasterGuest = &Config{guards: []*Guard{&PubGuard, &LoginGuard, &AuthenticateGuard}, lastGuard: &RootGuard}
)

type Guards []*Guard

type Config struct {
	guards    Guards
	lastGuard *Guard
}

func (c Config) Authenticate(r *http.Request) (*Principle, int, error) {
	code := 404
	errorMessage := "resource not found"
	for _, g := range c.guards {
		guard := *g
		if !guard.Matches(r) {
			continue
		}
		code = 401
		errorMessage = "unauthorized"
		principle, err := guard.Authenticate(r)
		if err != nil {
			continue
		}
		ok, err := guard.Authorize(principle, r)
		if !ok || err != nil {
			code = 403
			errorMessage = "access denied"
			continue
		}

		return principle, 200, nil
	}
	if c.lastGuard == nil {
		return nil, code, errors.New(errorMessage)
	}
	guard := *c.lastGuard
	if !guard.Matches(r) {
		return nil, code, errors.New(errorMessage)
	}
	code = 401
	principle, err := guard.Authenticate(r)
	console.RedPrintln("Request is not authenticated")
	if principle == nil || err != nil {
		return nil, code, err
	}
	ok, err := guard.Authorize(principle, r)
	console.RedPrintln("Request is not authorized")
	if !ok || err != nil {
		return nil, 403, err
	}
	return principle, 200, nil
}

func (c Config) Matches(r *http.Request) bool {
	for _, g := range c.guards {
		guard := *g
		if guard.Matches(r) {
			return true
		}
	}
	if c.lastGuard == nil {
		return false
	}

	return (*c.lastGuard).Matches(r)
}

func (c Config) AuthorizeRequest(f http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			currentPrinciple := FromRequest(r)
			if currentPrinciple != nil {
				console.YellowPrintln("There is a principle in the context skipping authorization")
				f.ServeHTTP(w, r)
				return
			}
			principle, code, err := c.Authenticate(r)
			if err != nil {
				console.RedPrint(err.Error())
				router.ReturnError(err, code, w)
				return
			}
			f.ServeHTTP(w, WithPrinciple(r, principle))
		})
}

func RulesSecurityConfig(rules []Rule) *Config {
	guards := make([]*Guard, len(rules)+1)
	guards[0] = &LoginGuard
	for i, rule := range rules {
		var guard Guard
		guard = rule
		guards[i+1] = &guard
	}
	return &Config{guards: guards, lastGuard: &RootGuard}
}
