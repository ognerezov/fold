package util

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"io"
	"os"
)

var (
	AllowedHiddenFiles = []string{".DS_Store"}
)

func CloseFile(f *os.File) {
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
	defer CloseFile(f)
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
	defer CloseFile(f)
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
	defer CloseFile(f)
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

func IsEmptyDir(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer CloseFile(f)

	files, err := f.ReadDir(1)
	if err != nil {
		return err
	}
	if len(files) != 0 {
		for _, fileName := range AllowedHiddenFiles {
			if fileName == files[0].Name() {
				return nil
			}
		}
		return fmt.Errorf("%s is not empty", name)
	}
	return nil
}
