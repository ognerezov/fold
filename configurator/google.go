package configurator

import (
	"errors"
	"fold/console"
	"fold/oauth"
	"fold/path"
	"io/fs"
	"path/filepath"
)

func InitGoogleProvider(dataPath string) (*oauth.GoogleJson, error) {
	clean := path.CreateRootCleaner(dataPath)
	googleJson := &oauth.GoogleJson{}
	googleJson.Iss = config.Name
	gotAny := false
	err := path.ProcessPath(dataPath, func(_p string, info fs.FileInfo, err error) error {
		p := filepath.ToSlash(_p)
		if info.IsDir() {
			return nil
		}

		_, filename, extension := path.Structure(dataPath, p, info, clean)
		switch extension {
		case ".json":
			err = googleJson.AttachFile(filename)
			if err == nil {
				gotAny = true
			}
		default:
			return nil
		}
		return nil
	})
	if err != nil {
		console.RedPrintln(err.Error())
	}
	oauth.SetGoogleJson(googleJson)
	if gotAny {
		return googleJson, nil
	}

	return nil, errors.New("no valid Google json file found")
}
