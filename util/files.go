package util

import (
	"encoding/json"
	"fold/console"
	"io"
	"os"
)

func CloseFie(f *os.File) {
	err := f.Close()
	if err != nil {
		console.RedPrintln(err.Error())
	}
}

func ReadFile(filePath string) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}
	defer CloseFie(f)
	bytes, err := io.ReadAll(f)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}
	return bytes, nil
}

func DoesFileExist(filePath string) bool {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	defer CloseFie(f)
	return true
}

func FromJson[T any](filename string, out *T) error {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(f)

	decoder := json.NewDecoder(f)
	return decoder.Decode(out)
}
