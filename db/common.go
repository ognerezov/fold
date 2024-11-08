package db

import (
	"encoding/json"
	"fold/console"
	"os"
)

func Save(file string, value any) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln("Error closing file " + file)
		}
	}(f)
	bytes, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}
