package initiator

import (
	"fmt"
	"fold/console"
	"fold/interfaces"
	"io/fs"
	"strings"
)

func readEmbeddedFolder(folder string) ([]fs.DirEntry, bool, error) {
	console.YellowPrintln("exporting folder: " + folder)
	files, err := dataOs.ReadDir(folder)
	if err != nil {
		return nil, false, err
	}
	isIndex := folder == fmt.Sprintf("%s/%s", embeddedRoot, Index)
	return files, isIndex, nil
}

func embeddedPath(folder string) string {
	return fmt.Sprintf("%s/%s", embeddedRoot, folder)
}

func structureFileName(folder string, file fs.DirEntry) (string, string) {
	inputFileName := folder + "/" + file.Name()
	fileName := cleanIndex(inputFileName, file.IsDir())
	outputFile := strings.TrimPrefix(fileName, embeddedRoot)
	return inputFileName, outputFile
}

func readEmbeddedFile(src string) (fs.File, interfaces.F, error) {
	console.YellowPrintln("exporting file: " + src)
	in, err := dataOs.Open(src)
	if err != nil {
		return nil, nil, err
	}
	return in, func() {
		err = in.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}, nil
}
