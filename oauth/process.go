package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/security"
	"fold/util"
	"net/http"
	"net/url"
	"time"
)

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Hd            string `json:"hd"`
}

func (u UserInfo) Principle() (*security.Principle, error) {
	if u.Email == "" {
		return nil, errors.New("user email is empty")
	}
	return &security.Principle{
		Id:    u.Email,
		Roles: []string{"user"},
	}, nil
}

func (tr *CodeTokenResponse) ExchangeForToken(provider string, uri string, iss string) (string, error) {
	userInfo, err := tr.GetUserInfo(uri)
	if err != nil {
		return "", err
	}

	principle, err := userInfo.Principle()
	if err != nil {
		return "", err
	}
	token, err := principle.TokenFor(iss)
	if err != nil {
		return "", err
	}
	table, ok := mem.TheStore.GetTable("/user/oauth")
	if ok {
		_, err = table.Insert(tr.GetRow(token, userInfo.Email, provider))
		if err != nil {
			console.RedPrintln("error saving oauth to table")
		}
	} else {
		console.RedPrintln("oauth table not found, skip oath storage")
	}

	return token, nil
}

func (tr *CodeTokenResponse) GetRow(token string, email string, provider string) map[string]string {
	return map[string]string{
		"id":            token,
		"user_id":       email,
		"provider":      provider,
		"access_token":  tr.AccessToken,
		"expires_in":    fmt.Sprintf("%v", tr.ExpiresIn),
		"token_type":    tr.TokenType,
		"scope":         tr.Scope,
		"refresh_token": tr.RefreshToken,
		"id_token":      tr.IdToken,
		"updated":       time.Now().String(),
	}
}

func (tr *CodeTokenResponse) GetUserInfo(uri string) (*UserInfo, error) {
	data := url.Values{}
	data.Set("access_token", tr.AccessToken)
	fmt.Println(uri + "?" + data.Encode())
	request, err := http.NewRequest("GET", uri+"?"+data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := util.SendRequest(request)
	if err != nil {
		return nil, err
	}
	defer util.HideBody(resp.Body)
	decoder := json.NewDecoder(resp.Body)
	var userInfo UserInfo
	err = decoder.Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println(userInfo)
	return &userInfo, nil
}
