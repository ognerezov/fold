package oauth

import (
	"fmt"
	"fold/console"
	"fold/mem"
	"golang.org/x/oauth2"
	"time"
)

type TokenSource struct {
	InternalToken string
	Table         *mem.Table
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	row := ts.Table.GetRow(ts.InternalToken)
	if row == nil {
		return nil, fmt.Errorf("user not found")
	}
	oauth := ts.Table.MapRow(row)
	const layout = "2025-02-22 22:26:44.727221 +0100 CET m=+2789.332512543"

	tm, _ := time.Parse(layout, oauth["updated"].(string))
	expires := oauth["expires_in"].(int64)
	tm.Add(time.Duration(expires) * time.Second)
	token := &oauth2.Token{
		AccessToken:  oauth["access_token"].(string),
		TokenType:    oauth["token_type"].(string),
		RefreshToken: oauth["refresh_token"].(string),
		ExpiresIn:    expires,
		Expiry:       tm,
	}
	return token, nil
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
