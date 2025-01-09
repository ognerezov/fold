package oauth

import (
	"fmt"
	"fold/console"
	"fold/mem"
)

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
