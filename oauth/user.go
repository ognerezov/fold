package oauth

import (
	"errors"
	"fold/arguments"
	"fold/security"
)

type Mapping interface {
	Map() map[string]string
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Hd            string `json:"hd"`
	Token         string `json:"token"`
}

func (u UserInfo) Principle() (*security.Principle, error) {
	if u.Email == "" {
		return nil, errors.New("user email is empty")
	}
	roles := security.GetUserRoles(u.Email, arguments.AppArguments.ApiPath)
	return &security.Principle{
		Id:    u.Email,
		Roles: roles,
	}, nil
}

func (u UserInfo) Map() map[string]string {
	return map[string]string{
		"id":   u.Email,
		"name": u.Name,
	}
}
