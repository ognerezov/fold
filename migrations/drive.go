package migrations

import (
	"bytes"
	"errors"
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/oauth"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"strings"
)

type DriverHandler struct {
	googleJson *oauth.GoogleJson
	root       string
	folders    map[string]*string
}

func CreateDriveHandler(gj *oauth.GoogleJson, root string) *DriverHandler {
	return &DriverHandler{googleJson: gj, root: root, folders: make(map[string]*string)}
}

func DriveImporter(arguments *arguments.InitArguments) (Importer, error) {
	credFile := arguments.CredentialsFile
	if credFile == "" {
		return nil, errors.New("no credentials file specified")
	}
	driveFolder := arguments.Destination
	if driveFolder == "" {
		return nil, errors.New("no drive folder specified")
	}
	gj, err := oauth.FromFile(credFile)
	if err != nil {
		return nil, err
	}
	var importer Importer
	importer = CreateDriveHandler(gj, driveFolder)

	return importer, nil
}

func (dh DriverHandler) GetFolderId(folder string, parent string) (*string, error) {
	if folder == "" {
		return &dh.root, nil
	}

	id, ok := dh.folders[folder]
	if ok {
		return id, nil
	}

	newId, err := dh.googleJson.CreateFolder(folder, []string{parent})
	if err != nil {
		return nil, err
	}
	dh.folders[folder] = newId
	return newId, nil
}

func (dh DriverHandler) CreateFolder(name string) error {
	path := strings.Split(name, "/")
	if len(path) == 0 {
		return nil
	}

	currentRoot := path[0]
	currentRootId, err := dh.GetFolderId(currentRoot, dh.root)
	if err != nil || currentRootId == nil {
		return err
	}

	for index, subfolder := range path {
		if index == 0 {
			continue
		}
		currentRootId, err = dh.GetFolderId(subfolder, *currentRootId)
		if err != nil || currentRootId == nil {
			return err
		}
	}

	return nil
}

func (dh DriverHandler) GetRegisteredFolder(dataPath string) (*string, error) {
	path := strings.Split(dataPath, "/")
	var folderId *string
	folder := path[len(path)-1]
	if folder == "" {
		folderId = &dh.root
	} else {
		f, ok := dh.folders[folder]
		if !ok {
			return nil, errors.New("Could not find folder " + folder)
		}
		folderId = f
	}
	return folderId, nil
}

func (dh DriverHandler) SaveFile(data FileData) error {
	folderId, err := dh.GetRegisteredFolder(data.Path)
	if err != nil {
		console.RedPrintln(err.Error())
		return err
	}

	if data.Table != nil {
		spreadsheet := &sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: data.Filename,
			},
		}
		createdSpreadsheet, err := dh.googleJson.Sheets.Spreadsheets.Create(spreadsheet).Do()
		if err != nil {
			return err
		}
		table := data.Table
		table.Spreadsheet = createdSpreadsheet

		_, err = dh.googleJson.Sheets.Spreadsheets.BatchUpdate(
			createdSpreadsheet.SpreadsheetId,
			&sheets.BatchUpdateSpreadsheetRequest{
				Requests: table.CreateSheetRequest(),
			}).Do()
		if err != nil {
			return err
		}
		file, err := dh.googleJson.Drive.Files.Get(createdSpreadsheet.SpreadsheetId).Fields("parents").Do()
		if err != nil {
			return err
		}
		_, err = dh.googleJson.Drive.Files.Update(createdSpreadsheet.SpreadsheetId, &drive.File{}).
			AddParents(*folderId).
			RemoveParents(file.Parents[0]).
			Do()
		return nil
	}
	driveFile := &drive.File{
		Name:     data.Filename,
		MimeType: data.MimeType,
		Parents:  []string{*folderId},
	}

	_, err = dh.googleJson.Drive.Files.Create(driveFile).Media(bytes.NewReader(data.Binary)).Do()
	return err
}

func (dh DriverHandler) RegisterFolder(id *string, name string) {
	fmt.Println(fmt.Sprintf("Registering folder %s refer %s", name, *id))
	dh.folders[name] = id
}

func (dh DriverHandler) CreateOrUpdateFile(fileName string, folder string, binary []byte, mime string) (*drive.File, error) {
	folderId, err := dh.GetRegisteredFolder(folder)
	if err != nil {
		return nil, err
	}
	return dh.googleJson.SaveFile(fileName, *folderId, binary, mime)
}
