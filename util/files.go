package util

import (
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
