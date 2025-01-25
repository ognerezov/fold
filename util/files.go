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

func Save(file string, bytes []byte) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer CloseFie(f)
	_, err = f.Write(bytes)
	return err
}

func SaveJavascript(path string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	str := "export default\n" + string(bytes)
	err = Save(path, []byte(str))
	if err != nil {
		return err
	}
	return nil
}
