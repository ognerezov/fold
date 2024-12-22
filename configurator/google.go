package configurator

import (
	"errors"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/path"
	"fold/util"
	"io/fs"
)

func InitGoogleProvider(dataPath string) (*api.GoogleJson, error) {
	clean := path.CreateRootCleaner(dataPath)
	googleJson := &api.GoogleJson{}
	gotAny := false
	err := path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, filename, extension := path.Structure(dataPath, p, info, clean)
		switch extension {
		case ".json":
			console.YellowPrintln(fmt.Sprintf("Reading Google json file from %s", filename))
			var json *api.GoogleJson
			err = util.FromJson(filename, &json)
			if err != nil {
				console.RedPrintln(err.Error())
				return err
			}
			googleJson.Attach(json)
			gotAny = true
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		console.RedPrintln(err.Error())
	}
	if gotAny {
		fmt.Println(*(googleJson.Web))
		fmt.Println(*(googleJson.Installed))
		return googleJson, nil
	}

	return nil, errors.New("no valid Google json file found")
}
