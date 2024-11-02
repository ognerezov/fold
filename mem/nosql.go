package mem

import (
	"encoding/json"
	"errors"
	"fold/console"
	"io"
	"os"
)

type NoSql struct {
	File       string
	document   map[string]any
	collection []map[string]any
	is         string
}

func (n *NoSql) Val() any {
	if n.is == Array {
		return n.collection
	}
	return n.document
}

func LoadJson(file string) (*NoSql, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln("Error closing json " + file)
		}
	}(f)
	raw, _ := io.ReadAll(f)
	var data any
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	switch x := data.(type) {
	case []map[string]any:
		return &NoSql{collection: x, is: Array, File: file}, nil
	case map[string]any:

		return &NoSql{document: x, is: Struct, File: file}, nil
	default:
		return nil, errors.New("invalid type")
	}
}
