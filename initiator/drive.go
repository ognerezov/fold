package initiator

import (
	"fold/arguments"
	"fold/configurator"
	"fold/migrations"
)

func InitDrive(args *arguments.InitArguments, config *configurator.AppConfig) {
	importer, err := migrations.DriveImporter(args)
	if err != nil {
		panic(err)
	}

	exporter, err := ConfigureExporter(args, config)
	if err != nil {
		panic(err)
	}

	err = exporter.Process(importer)
	if err != nil {
		panic(err)
	}
}
