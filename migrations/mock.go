package migrations

import (
	"fmt"
	"fold/arguments"
	"fold/console"
)

type MockImporter bool

func (m MockImporter) CreateFolder(name string) error {
	console.CyanPrintln(fmt.Sprintf("Creating folder %s", name))
	return nil
}
func (m MockImporter) SaveFile(data FileData) error {
	console.GreenPrintln(fmt.Sprintf("Saving file %s to %s", data.Filename, data.Path))
	return nil
}

func MockImport(*arguments.InitArguments) (Importer, error) {
	var importer Importer
	importer = MockImporter(true)
	return importer, nil
}
