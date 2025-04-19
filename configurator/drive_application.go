package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/interfaces"
	"fold/migrations"
	"fold/oauth"
	"fold/openapi"
	"fold/security"
	"fold/util"
	goji "goji.io"
	"google.golang.org/api/drive/v3"
	"strings"
)

func CreateDriveApplication(address string, folder string, credentialsFile string) {
	gj, err := oauth.FromFile(credentialsFile)
	AppProviders.Google = gj
	if err != nil {
		panic(err)
	}
	files, err := gj.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	ports := PortsConfig{}
	var resourceFolder *drive.File
	for _, file := range files {
		fileName := file.Name

		if file.MimeType == "application/vnd.google-apps.folder" {
			console.YellowPrintln("File is a nested folder")
			if fileName == "www" {
				resourceFolder = file
				continue
			}
			port, found := util.PathToInt(fileName)

			if !found {
				console.RedPrintln("Unknown folder name " + fileName)
				continue
			}
			ports = append(ports, PortConfig{
				port:      port,
				path:      file.Id,
				configure: ConfigureDriveServer,
				driveFile: file,
			})
			continue
		}

		fileHandler := DriveFileHandler(*file)
		if file.Name == projectJson {
			bytes, err := fileHandler.Fetch()
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(bytes, &config)

			if err != nil {
				panic(err)
			}
			continue
		}
		console.RedPrintln(fmt.Sprintf("Found file %s is of type %s in a root folder is skipped", fileName, file.MimeType))
	}
	if resourceFolder != nil {
		console.GreenPrintln(fmt.Sprintf("Found resource folder: %s", resourceFolder.Name))
		// TODO configure resources
	}

	if len(ports) == 0 {
		panic(errors.New("no server folders found"))
	}

	ServePorts(address, ports)
}

func ConfigureDriveServer(folderId string, port int) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + folderId)
	mux, store, apiDescription := initialize(folderId, port)

	//TODO add oauth client to provider

	next := interfaces.NewPhase()
	controlEndpoints := make(Endpoints)
	migrationHandler := migrations.CreateDriveHandler(AppProviders.Google, folderId)
	migrationHandlers[folderId] = migrationHandler
	SetDriveHandlers(arguments.AppArguments.ApiPath, folderId, mux, apiDescription, next, migrationHandler, controlEndpoints)
	next.Act()
	setupSecurity(store, mux, controlEndpoints)
	err := SaveAndServe(projectRoute, config, migrationHandler, mux, apiDescription)
	if err != nil {
		console.RedPrintln(err.Error())
	}
	err = SaveAndServe(controlRoute, controlEndpoints, migrationHandler, mux, apiDescription)
	if err != nil {
		console.RedPrintln(err.Error())
	}

	err = SaveAndServe(interfaces.ProvidersOutputPath+"google.js", AppProviders.Google.WithoutSecret(), migrationHandler, mux, apiDescription)
	if err != nil {
		console.RedPrintln(err.Error())
	}

	security.SetAuthHandlers(arguments.AppArguments.ApiPath, mux, App.Name())

	err = SaveAndServe(openapi.Filename, apiDescription, migrationHandler, mux, apiDescription)
	if err != nil {
		console.RedPrintln(err.Error())
	}
	return mux, nil
}

func SaveAndServe(route string, data any, migrationHandler *migrations.DriverHandler, mux *goji.Mux, api *openapi.ApiDescription) error {
	binary, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	path := strings.Split(route, "/")
	fileName := path[len(path)-1]
	folder := strings.Join(path[:len(path)-1], "/")
	m, ext, _ := GetContentType(fileName)
	if ext == ".js" {
		str := "export default\n" + string(binary)
		binary = []byte(str)
	}
	file, err := migrationHandler.CreateOrUpdateFile(fileName, folder, binary, m)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fileHandler := DriveFileHandler(*file)
	SetRawDriveHandlers(route, fileHandler, m, binary, mux, api)
	return nil
}
