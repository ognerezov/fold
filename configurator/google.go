package configurator

import (
	"context"
	"errors"
	"fmt"
	"fold/console"
	"fold/oauth"
	"fold/path"
	"fold/util"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/fs"
)

func InitGoogleProvider(dataPath string) (*oauth.GoogleJson, error) {
	clean := path.CreateRootCleaner(dataPath)
	googleJson := &oauth.GoogleJson{}
	gotAny := false
	apiClientCreated := false
	err := path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, filename, extension := path.Structure(dataPath, p, info, clean)
		switch extension {
		case ".json":
			console.YellowPrintln(fmt.Sprintf("Reading Google json file from %s", filename))
			var json *oauth.GoogleJson
			err = util.FromJson(filename, &json)
			if err != nil {
				console.RedPrintln(err.Error())
				return err
			}

			if !apiClientCreated && json.PrivateKey != "" {
				console.MagentaPrintln("Found Google service account json")
				ctx := context.Background()
				var drv *drive.Service
				drv, err = drive.NewService(ctx, option.WithCredentialsFile(filename))
				googleJson.Drive = drv
				googleJson.Iss = config.Name
				if err != nil {
					console.RedPrintln(err.Error())
				} else {
					apiClientCreated = true
					console.MagentaPrintln("Google API Client created")
				}
				s, err := sheets.NewService(ctx, option.WithCredentialsFile(filename))
				if err != nil {
					console.RedPrintln(err.Error())
				} else {
					googleJson.Sheets = s
				}
				cloud, err := cloudresourcemanager.NewService(ctx, option.WithCredentialsFile(filename))
				if err != nil {
					console.RedPrintln(err.Error())
				} else {
					googleJson.CloudResourceManager = cloud
				}
			} else {
				googleJson.Attach(json)
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
		//fmt.Println(*(googleJson.Web))
		//fmt.Println(*(googleJson.Installed))
		return googleJson, nil
	}

	return nil, errors.New("no valid Google json file found")
}
