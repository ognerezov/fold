package security

import (
	"context"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/util"
	"net/http"
)

type Principle struct {
	Roles    []string `json:"roles"`
	Id       string   `json:"id"`
	password string
}

var Guest = Principle{
	Id:       "guest",
	Roles:    []string{"pub"},
	password: guestPassword,
}

var Root = Principle{
	Id:       "root",
	Roles:    []string{"root", "admin", "user", "pub"},
	password: adminPassword,
}

func WithPrinciple(r *http.Request, p *Principle) *http.Request {
	ctx := r.Context()
	return r.Clone(context.WithValue(ctx, "principle", p))
}

func FromRequest(r *http.Request) *Principle {
	ctx := r.Context()
	return ctx.Value("principle").(*Principle)
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

func FromToken(tokenString string) (*Principle, error) {
	token, err := verifyToken(tokenString[len("Bearer "):])
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
		Id:    sub,
		Roles: roles,
	}, nil
}

func FromTable(id string, password string) (*Principle, error) {
	userData := mem.TheStore.PlainGet(util.UserPath, id)
	if userData == nil {
		return nil, fmt.Errorf("user not found")
	}
	e := verifyPassword(password, userData)
	if e != nil {
		return nil, e
	}

	rolesTable, _ := mem.TheStore.GetTable("/user/roles")
	roleRows := rolesTable.SearchRows("user_id", id)
	roles := make([]string, len(roleRows))
	for i, row := range roleRows {
		roles[i] = row[2].Str()
	}

	return &Principle{
		Id:    id,
		Roles: roles,
	}, nil
}

func verifyPassword(password string, data map[string]string) error {
	p := data["password"]
	if !CheckPassword(password, p) {
		return fmt.Errorf("password missmatch")
	}
	return nil
}
