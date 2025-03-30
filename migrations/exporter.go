package migrations

import (
	"fold/arguments"
	"fold/mem"
)

type FileData struct {
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
	Path     string `json:"path"`
	NoSql    *mem.NoSql
	Binary   []byte
	Table    *mem.Table
}

type Importer interface {
	CreateFolder(name string) error
	SaveFile(data FileData) error
}

type Exporter interface {
	Process(importer Importer) error
}

type ImporterFactory func(arguments *arguments.InitArguments) (Importer, error)
type ExporterFactory func(arguments *arguments.InitArguments) (Exporter, error)
