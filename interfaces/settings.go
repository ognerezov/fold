package interfaces

import (
	"fold/console"
	"fold/oauth"
	"os"
)

const (
	ProvidersOutputPath = "/providers/"
)

type Providers struct {
	Google *oauth.GoogleJson
}

func (p *Providers) Export(path string) error {
	err := os.MkdirAll(path+ProvidersOutputPath, os.ModePerm)
	if p.Google != nil {
		err = p.Google.Export(path + ProvidersOutputPath + "google")
		if err != nil {
			console.RedPrintln(err.Error())
		}
		return err
	}
	return nil
}
