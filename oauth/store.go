package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/util"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

type TokenSource struct {
	InternalToken string
	Table         *mem.Table
	GoogleJson    *GoogleJson
}

func (ts *TokenSource) Provider() string {
	if ts.GoogleJson != nil {
		return "google"
	}
	return ""
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	row := ts.Table.GetRow(ts.InternalToken)
	if row == nil {
		return nil, fmt.Errorf("user not found")
	}
	oauth := ts.Table.MapRow(row)

	token := mapToken(oauth)
	if util.IsTimePast(token.Expiry, 10) {
		tokenResponse, err := ts.RefreshToken(token.RefreshToken)
		fmt.Println(tokenResponse)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		email := oauth["user_id"].(string)
		row, err := ts.Table.Update(ts.InternalToken, tokenResponse.GetRow(ts.InternalToken, email, ts.Provider()))
		if err != nil {
			return nil, err
		}
		return mapToken(row), nil
	}
	return token, nil
}

func (ts *TokenSource) RefreshToken(refreshToken string) (*GoogleTokenResponse, error) {
	req, err := ts.GoogleJson.RefreshRequest(refreshToken)
	if err != nil {
		return nil, err
	}

	resp, err := util.SendRequest(req)
	if err != nil {
		return nil, err
	}
	defer util.HideBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		var errResp map[string]any
		err = decoder.Decode(&errResp)
		console.RedPrintln(fmt.Sprintf("RefreshToken error: %v, response code: %v", errResp, resp.Status))
		return nil, errors.New(fmt.Sprintf("server reponded with status %v", resp.Status))
	}
	decoder := json.NewDecoder(resp.Body)
	var tokenResp GoogleTokenResponse
	err = decoder.Decode(&tokenResp)
	if err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (ts *TokenSource) Build(ctx context.Context) (*oauth2.TokenSource, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	res := googleJson.OAuthConfig().TokenSource(ctx, token)
	return &res, nil
}

func mapToken(oauth map[string]any) *oauth2.Token {
	const layout = "2006-01-02 15:04:05.999999999 -0700 MST"
	updated := oauth["updated"].(string)
	tm, err := time.Parse(layout, strings.Split(updated, " m=")[0])
	if err != nil {
		console.RedPrintln(err.Error())
		return nil
	}
	expires := oauth["expires_in"].(int64)

	return &oauth2.Token{
		AccessToken:  oauth["access_token"].(string),
		TokenType:    oauth["token_type"].(string),
		RefreshToken: oauth["refresh_token"].(string),
		ExpiresIn:    expires,
		Expiry:       tm.Add(time.Duration(expires) * time.Second),
	}
}

func (ts *TokenSource) Interface() oauth2.TokenSource {
	return ts
}

func StoreOauth(row map[string]string) {
	table, ok := mem.TheStore.GetTable("/user/oauth")
	if ok {
		_, err := table.Insert(row)
		if err != nil {
			console.RedPrintln("error saving oauth to table")
		}
	} else {
		console.RedPrintln("oauth table not found, skip oath storage")
	}
}

func SourceToken(internalToken string) TokenSource {
	table, ok := mem.TheStore.GetTable("/user/oauth")
	if !ok {
		panic("oauth table not found")
	}

	return TokenSource{
		InternalToken: internalToken,
		Table:         table,
		GoogleJson:    googleJson,
	}
}

func TradeToken(internalToken string) *string {
	table, ok := mem.TheStore.GetTable("/user/oauth")
	if !ok {
		console.RedPrintln("oauth table not found")
		return nil
	}
	row := table.GetRow(internalToken)
	if row == nil {
		console.RedPrintln("oauth token not found")
		return nil
	}
	return table.GetStringValue(row, "access_token")
}

func RemovePreviousRecords(user string, provider string) {
	table, ok := mem.TheStore.GetTable("/user/oauth")
	rawQuery := map[string][]string{
		"provider": {provider},
		"user_id":  {user},
	}
	if ok {
		ids := table.RawQueryIds(rawQuery)
		for _, id := range ids {
			table.DeleteById(id)
		}
		console.YellowPrintln(fmt.Sprintf("%v previous oauth records removed", len(ids)))
	} else {
		console.RedPrintln("oauth table not found, skip oath storage")
	}
}
