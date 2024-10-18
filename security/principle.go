package security

import (
	"context"
	"fmt"
	"net/http"
)

type Principle struct {
	roles []string
	id    string
}

var Guest = Principle{
	id:    "guest",
	roles: []string{"pub"},
}

var Root = Principle{
	id:    "root",
	roles: []string{"root", "admin", "user", "pub"},
}

func WithPrinciple(r *http.Request, p *Principle) *http.Request {
	ctx := r.Context()
	return r.Clone(context.WithValue(ctx, "principle", p))
}

func GetPrinciple(r *http.Request) *Principle {
	ctx := r.Context()
	return ctx.Value("principle").(*Principle)
}

func (principle *Principle) BearerToken() (string, error) {
	j, er := createToken(principle)
	if er != nil {
		return "", er
	}
	return fmt.Sprintf("Bearer %s", j), nil
}

func FromToken(tokeString string) (*Principle, error) {
	token, err := verifyToken(tokeString)
	if err != nil || token == nil {
		return nil, err
	}

	claims := token.Claims
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	aud, err := claims.GetAudience()
	var roles []string
	if err != nil || aud == nil {
		roles = []string{}
	} else {
		roles = aud
	}
	return &Principle{
		id:    sub,
		roles: roles,
	}, nil
}
