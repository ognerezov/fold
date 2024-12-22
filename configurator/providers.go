package configurator

import (
	"fmt"
	"fold/api"
	"fold/console"
	"os"
)

const (
	ProvidersOutputPath = "/providers/"
)

var (
	AppProviders                                    = &Providers{}
	GoogleProvider ProviderExporter[api.GoogleJson] = &api.GoogleJson{}
)

type Providers struct {
	Google *api.GoogleJson
}

// TODO create folder
func (p *Providers) Export(path string) error {
	err := os.MkdirAll(path+ProvidersOutputPath, os.ModePerm)
	if p.Google != nil {
		err = p.Google.Export(path + ProvidersOutputPath + "google.json")
		if err != nil {
			console.RedPrintln(err.Error())
		}
		return err
	}
	return nil
}

type ProviderExporter[T any] interface {
	WithoutSecret() T
	Export(string) error
}

func InitProviders(dataPath string) error {
	files, err := ReadDir(dataPath + "/providers")
	if err != nil {
		console.RedPrintln(err.Error())
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		name := file.Name()
		if name == "google" {
			console.YellowPrintln(fmt.Sprintf("Loading provider %s", name))
			var provider *api.GoogleJson
			provider, err = InitGoogleProvider(fmt.Sprintf("%s/providers/%s", dataPath, name))
			if err != nil {
				console.RedPrintln(err.Error())
			} else {
				AppProviders.Google = provider
				fmt.Println(provider.Web)
				fmt.Println(provider.Installed)
			}

			continue
		}
		console.RedPrintln(fmt.Sprintf("Unknown provider %s", name))
	}
	return nil
}
