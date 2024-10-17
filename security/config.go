package security

import (
	"errors"
	"fold/console"
	"fold/router"
	"net/http"
)

var (
	Public = &Config{guards: make(Guards, 0), lastGuard: &NoGuard}
)

type Guards []*Guard

type Config struct {
	guards    Guards
	lastGuard *Guard
}

func (c Config) Authenticate(r *http.Request) (*Principle, error) {
	for _, g := range c.guards {
		guard := *g
		if !guard.Matches(r) {
			continue
		}
		principle, err := guard.Authenticate(r)
		if err != nil {
			return nil, err
		}
		ok, err := guard.Authorize(principle, r)
		if !ok || err != nil {
			return nil, err
		}

		return &principle, nil
	}
	if c.lastGuard == nil {
		return nil, errors.New("no guard found")
	}
	guard := *c.lastGuard
	principle, err := guard.Authenticate(r)
	ok, err := guard.Authorize(principle, r)
	if !ok || err != nil {
		return nil, err
	}

	return &principle, nil
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
			principle, err := c.Authenticate(r)
			if err != nil {
				e := router.WriteError(err, w)
				console.RedPrintln(e.Error())
				return
			}
			f.ServeHTTP(w, WithPrinciple(r, principle))
		})
}
