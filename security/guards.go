package security

import (
	"net/http"
)

const (
	blindGuard = BlindGuard(true)
)

var (
	NoGuard Guard = blindGuard
)

type BlindGuard bool

func (b BlindGuard) Authorize(Principle, *http.Request) (bool, error) {
	return bool(b), nil
}

func (b BlindGuard) Matches(*http.Request) bool {
	return true
}

func (b BlindGuard) Authenticate(*http.Request) (Principle, error) {
	return Guest, nil
}
