package initiator

import (
	"embed"
	"fmt"
	"fold/console"
	"fold/util"
)

const (
	DefaultTemplate = "default"
)

var (
	Initiators = map[string]Initiator{
		DefaultTemplate: DefaultInit,
	}
)

func Init(fs *embed.FS, template string, path string) {
	console.CyanPrintln(fmt.Sprintf("Initializing new project with template %s in folder %s", template, path))
	err := util.IsEmptyDir(path)
	if err != nil {
		panic(err)
	}
	print(fs)
	initiator, ok := Initiators[template]
	if !ok {
		panic("template not found " + template)
	}
	err = initiator(path)
	if err != nil {
		panic(err)
	}
}

type Initiator func(string) error

func DefaultInit(path string) error {

	return nil
}
