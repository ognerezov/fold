package initiator

import (
	"encoding/json"
	"fold/arguments"
	"fold/configurator"
	"fold/migrations"
	"strconv"
)

func InitDrive(args *arguments.InitArguments, config *configurator.AppConfig) {
	importer, err := migrations.DriveImporter(args)
	if err != nil {
		panic(err)
	}
	portFolder := strconv.Itoa(args.Port)
	err = importer.CreateFolder(portFolder)
	if err != nil {
		panic(err)
	}

	if config.AuthProviders == nil || len(config.AuthProviders) == 0 {
		config.AuthProviders = []string{"google"}
	} else {
		config.AuthProviders = append(config.AuthProviders, "google")
	}
	bytes, err := json.MarshalIndent(*config, "", "  ")
	if err != nil {
		panic(err)
	}
	data := migrations.FileData{
		Filename: "project.json",
		MimeType: "application/json",
		Binary:   bytes,
		Path:     "",
	}
	err = importer.SaveFile(data)
	if err != nil {
		panic(err)
	}

	p, err := importer.GetFolderId(portFolder, args.Destination)
	if err != nil {
		panic(err)
	}
	// Creating new exporter for port folder
	args.Destination = *p
	importer, err = migrations.DriveImporter(args)
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
