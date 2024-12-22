package db

import (
	"encoding/json"
	"fold/util"
)

func SaveJson(file string, value any) error {
	bytes, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return err
	}
	return util.Save(file, bytes)
}
