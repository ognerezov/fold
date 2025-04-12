package initiator

import (
	"fmt"
	"fold/arguments"
	"fold/configurator"
	"fold/console"
	"fold/mem"
	"fold/migrations"
	"fold/path"
	"io"
	"strings"
)

type ExporterConfig struct {
	Template string
	port     int
	config   configurator.AppConfig
}

func (ec ExporterConfig) Process(importer migrations.Importer) error {
	folders, ok := FoldersToCopy[ec.Template]
	if !ok {
		return fmt.Errorf("template folders not found %v", ec.Template)
	}
	//for _, folder := range folders {
	//	err := importer.CreateFolder(folder)
	//	if err != nil {
	//		return err
	//	}
	//}
	for _, folder := range folders {
		err := ProcessFolder(embeddedPath(folder), importer)
		if err != nil {
			return err
		}
	}
	return nil
}

func ConfigureExporter(args *arguments.InitArguments, config *configurator.AppConfig) (migrations.Exporter, error) {
	var exporter migrations.Exporter
	exporter = ExporterConfig{
		Template: args.Template,
		port:     args.Port,
		config:   *config,
	}
	return exporter, nil
}

func ProcessFolder(folder string, importer migrations.Importer) error {
	files, isIndex, err := readEmbeddedFolder(folder)
	if err != nil {
		return err
	}

	if !isIndex {
		err = importer.CreateFolder(strings.TrimPrefix(folder, embeddedRoot))
		if err != nil {
			return err
		}
	}
	for _, file := range files {
		inputFileName, outputFile := structureFileName(folder, file)
		if file.IsDir() {
			err = importer.CreateFolder(outputFile)
			if err != nil {
				return err
			}
			err = ProcessFolder(inputFileName, importer)
			if err != nil {
				return err
			}
			continue
		}

		err = ProcessFile(inputFileName, importer)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessFile(inputFile string, importer migrations.Importer) error {

	clean := path.CreateRootCleaner(embeddedRoot)

	route, _, fileName, ext, mime := path.PrepareFileName(embeddedRoot, inputFile, clean)
	for _, banned := range migrations.BanedExtensions {
		if ext == banned {
			console.MagentaPrintln("skip file " + fileName)
			return nil
		}
	}
	in, closeFile, err := readEmbeddedFile(inputFile)
	if err != nil {
		return err
	}
	defer closeFile()

	bytes, err := io.ReadAll(in)
	if err != nil {
		console.YellowPrintln(err.Error())
		return err
	}
	parts := strings.Split(fileName, "/")
	data := migrations.FileData{
		Filename: parts[len(parts)-1],
		MimeType: mime,
		Binary:   bytes,
		Path:     cleanIndex(route, true),
	}
	if ext == ".csv" {
		console.YellowPrintln("file is csv")
		table, err := mem.TableFromBytes(bytes)
		if err != nil {
			console.YellowPrintln(err.Error())
			return err
		}
		data.Table = table
	}
	if ext == ".json" {
		console.YellowPrintln("file is json")
		noSql, err := mem.FromBytes(bytes)
		if err != nil {
			console.YellowPrintln(err.Error())
			return err
		}
		data.NoSql = noSql
	}

	return importer.SaveFile(data)
}
