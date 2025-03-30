package migrations

import (
	"fmt"
	"fold/arguments"
)

const (
	Dir   = "dir"
	Drive = "drive"
)

var (
	Exporters = map[string]ExporterFactory{
		Dir: FsExporter,
	}
	Importers = map[string]ImporterFactory{
		Drive: DriveImporter,
	}
	BanedExtensions = []string{".DS_Store"}
)

func Migrate(arguments *arguments.InitArguments) {

	exporter, importer, err := ConfigureMigration(arguments)
	if err != nil {
		panic(err)
	}
	err = exporter.Process(importer)
	if err != nil {
		panic(err)
	}
}

func ConfigureMigration(args *arguments.InitArguments) (Exporter, Importer, error) {
	eFactory, ok := Exporters[args.SourceType]
	if !ok {
		return nil, nil, fmt.Errorf("unsupported export type: %s", args.SourceType)
	}

	iFactory, ok := Importers[args.DestinationType]
	if !ok {
		return nil, nil, fmt.Errorf("unsupported importer type: %s", args.DestinationType)
	}

	exporter, err := eFactory(args)
	if err != nil {
		return nil, nil, err
	}

	importer, err := iFactory(args)
	if err != nil {
		return nil, nil, err
	}

	return exporter, importer, nil
}
