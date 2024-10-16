package security

import (
	"fold/mem"
	"net/http"
)

type Guard interface {
	Authorize(p mem.Principle, r *http.Request) (bool, error)
	Matches(r *http.Request) bool
	Authenticate(r *http.Request) (mem.Principle, error)
}
