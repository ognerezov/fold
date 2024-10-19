package security

import (
	"context"
	"fmt"
	"fold/console"
	"fold/mem"
	"net/http"
)

type Principle struct {
	roles    []string
	id       string
	password string
}

var Guest = Principle{
	id:       "guest",
	roles:    []string{"pub"},
	password: guestPassword,
}

var Root = Principle{
	id:       "root",
	roles:    []string{"root", "admin", "user", "pub"},
	password: adminPassword,
}

func WithPrinciple(r *http.Request, p *Principle) *http.Request {
	ctx := r.Context()
	return r.Clone(context.WithValue(ctx, "principle", p))
}

func Authenticate(r *http.Request) (*Principle, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		console.RedPrint("Token not found")
		return nil, fmt.Errorf("no token")
	}

	principle, err := FromToken(tokenString)
	return principle, err
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

func FromTable(id string, password string) (*Principle, error) {
	userData := mem.TheStore.PlainGet("/user", id)
	if userData == nil {
		return nil, fmt.Errorf("user not found")
	}
	e := verifyPassword(password, userData)
	if e != nil {
		return nil, e
	}
	return &Principle{
		id:    id,
		roles: []string{"user"},
	}, nil
}

func verifyPassword(password string, data map[string]string) error {
	p := data["password"]
	if p != password {
		return fmt.Errorf("password missmatch")
	}
	return nil
}
