package security

import (
	"fold/mem"
	"net/http"
)

const (
	blindGuard = BlindGuard(true)
)

var (
	NoGuard Guard = blindGuard
)

type BlindGuard bool

func (b BlindGuard) Authorize(mem.Principle, *http.Request) (bool, error) {
	return bool(b), nil
}

func (b BlindGuard) Matches(*http.Request) bool {
	return true
}

func (b BlindGuard) Authenticate(*http.Request) (mem.Principle, error) {
	return mem.Guest, nil
}
