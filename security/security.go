package security

import (
	"net/http"
)

type Guard interface {
	Authorize(p Principle, r *http.Request) (bool, error)
	Matches(r *http.Request) bool
	Authenticate(r *http.Request) (Principle, error)
}
