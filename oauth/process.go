package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/util"
	"net/http"
	"net/url"
	"time"
)

func (tr *GoogleTokenResponse) ExchangeForToken(uri string, iss string) (*UserInfo, error) {
	userInfo, err := tr.GetUserInfo(uri)
	if err != nil {
		return nil, err
	}

	principle, err := userInfo.Principle()
	if err != nil {
		return nil, err
	}
	token, err := principle.TokenFor(iss)
	if err != nil {
		return nil, err
	}
	userInfo.Token = token

	return userInfo, nil
}

func (tr *GoogleTokenResponse) GetRow(token string, email string, provider string) map[string]string {
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

func (tr *GoogleTokenResponse) GetUserInfo(uri string) (*UserInfo, error) {
	data := url.Values{}
	data.Set("access_token", tr.AccessToken)

	request, err := http.NewRequest("GET", uri+"?"+data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := util.SendRequest(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		var errResp map[string]any
		err = decoder.Decode(&errResp)
		err = errors.New(resp.Status)

		return nil, err
	}
	defer util.HideBody(resp.Body)
	decoder := json.NewDecoder(resp.Body)
	var userInfo UserInfo
	err = decoder.Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
