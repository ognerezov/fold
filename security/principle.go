package security

import (
	"context"
	"net/http"
)

type Principle struct {
	roles []string
	id    string
}

var Guest = Principle{
	roles: []string{"guest"},
}

func WithPrinciple(r *http.Request, p *Principle) *http.Request {
	ctx := r.Context()
	return r.Clone(context.WithValue(ctx, "principle", p))
}

func GetPrinciple(r *http.Request) *Principle {
	ctx := r.Context()
	return ctx.Value("principle").(*Principle)
}
