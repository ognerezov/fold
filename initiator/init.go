package initiator

import (
	"embed"
	"fmt"
	"fold/console"
	"fold/util"
	"io"
	"io/fs"
	"os"
	"strings"
)

const (
	DefaultTemplate = "default"
	embeddedRoot    = "data"
)

var (
	Initiators = map[string]Initiator{
		DefaultTemplate: DefaultInit,
	}
	//go:embed data/*
	dataOs embed.FS
)

func Init(template string, path string, port int) {
	console.CyanPrintln(fmt.Sprintf("Initializing new project with template %s in folder %s", template, path))
	err := util.IsEmptyDir(path)
	if err != nil {
		panic(err)
	}

	initiator, ok := Initiators[template]
	if !ok {
		panic("template not found " + template)
	}
	prodPath := fmt.Sprintf("%s/%v", path, port)
	err = os.MkdirAll(prodPath, os.ModePerm)
	err = os.MkdirAll(fmt.Sprintf("%s/%s", prodPath, "www"), os.ModePerm)
	err = os.MkdirAll(fmt.Sprintf("%s/%s", prodPath, "providers"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = initiator(path, port)
	if err != nil {
		panic(err)
	}
}

type Initiator func(string, int) error

func DefaultInit(path string, port int) error {
	portPath := fmt.Sprintf("%v/%v", path, port)
	err := exportFolders(portPath, []string{"user", "security", "pub"})
	if err != nil {
		return err
	}
	return nil
}

func exportFolders(path string, folders []string) error {
	for _, folder := range folders {
		err := exportFolder(path, fmt.Sprintf("%s/%s", embeddedRoot, folder))
		if err != nil {
			return err
		}
	}
	return nil
}

func exportFolder(path string, folder string) error {
	console.YellowPrintln("exporting folder: " + folder)
	files, err := dataOs.ReadDir(folder)
	if err != nil {
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s", path, strings.TrimPrefix(folder, embeddedRoot)), os.ModePerm)
	if err != nil {
		return err
	}
	for _, file := range files {
		fileName := folder + "/" + file.Name()
		if file.IsDir() {
			err = os.MkdirAll(fmt.Sprintf("%s/%s", path, fileName), os.ModePerm)
			if err != nil {
				return err
			}
			err = exportFolder(path, fileName)
			if err != nil {
				return err
			}
		}
		outputFile := strings.TrimPrefix(fileName, embeddedRoot)
		err = exportFile(fileName, path+outputFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func exportFile(src, dst string) error {
	console.YellowPrintln("exporting file: " + src)
	in, err := dataOs.Open(src)
	if err != nil {
		return err
	}
	defer func(in fs.File) {
		err = in.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer util.CloseFile(out)
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}
