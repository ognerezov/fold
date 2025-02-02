package configurator

import (
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/interfaces"
	"fold/oauth"
)

var (
	AppProviders                                      = &interfaces.Providers{}
	GoogleProvider ProviderExporter[oauth.GoogleJson] = &oauth.GoogleJson{}
)

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
			var provider *oauth.GoogleJson
			provider, err = InitGoogleProvider(fmt.Sprintf("%s/providers/%s", dataPath, name))
			if err != nil {
				console.RedPrintln(err.Error())
			} else {
				provider.RegistrationAllowed = arguments.AppArguments.RegistrationAllowed
				AppProviders.Google = provider
				(*TheInstructions)[googleAuth] = provider.AuthControl
			}

			continue
		}
		console.RedPrintln(fmt.Sprintf("Unknown provider %s", name))
	}
	return nil
}
