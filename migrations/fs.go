package migrations

import (
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/mem"
	"fold/path"
	"io/fs"
)

type DataPath string

func (d DataPath) Process(importer Importer) error {
	dataPath := string(d)
	clean := path.CreateRootCleaner(dataPath)
	return path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		fmt.Println(p)
		if p == dataPath {
			return nil
		}
		if info.IsDir() {
			return importer.CreateFolder(clean(p))
		}
		route, fullName, fileName, ext, mime := path.PrepareFileName(dataPath, p, clean)
		for _, banned := range BanedExtensions {
			if ext == banned {
				console.MagentaPrintln("skip file " + fileName)
				return nil
			}
		}
		fileHandler := mem.FilePath(fullName)
		bytes, err := fileHandler.Fetch()
		if err != nil {
			console.YellowPrintln(err.Error())
			return err
		}

		data := FileData{
			Filename: fileName,
			MimeType: mime,
			Binary:   bytes,
			Path:     route,
		}
		if ext == ".csv" {
			console.YellowPrintln("file is csv")
			table, err := fileHandler.FetchCsv()
			if err != nil {
				console.YellowPrintln(err.Error())
				return err
			}
			data.Table = table
		}
		if ext == ".json" {
			console.YellowPrintln("file is json")
			noSql, err := fileHandler.FetchNoSql()
			if err != nil {
				console.YellowPrintln(err.Error())
				return err
			}
			data.NoSql = noSql
		}

		return importer.SaveFile(data)
	})
}

func FsExporter(arguments *arguments.InitArguments) (Exporter, error) {
	p := arguments.Source
	if p == "" {
		return nil, fmt.Errorf("path can't be empty")
	}
	var exporter Exporter
	exporter = DataPath(p)

	return exporter, nil
}
