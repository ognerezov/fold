package threads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fold/db"
	"google.golang.org/api/drive/v3"
)

type DriveWriter drive.File

func WriteDriveAsync(update db.DriveUpdate) {
	writer := DriveWriter(drive.File{
		Id:       update.File.Id,
		MimeType: update.File.MimeType,
		Name:     update.File.Name,
	})
	Async(writer, update.Value)
}

func (w DriveWriter) Call(val any) (Message[string], ErrorMessage) {
	driveService := providers.Google.Drive
	process := "SaveJson data to file " + w.Id
	if driveService == nil {
		e := fmt.Errorf("drive service not found in configurator")
		return EmptyMessage(process), CommonError(process, e)
	}
	b, err := json.MarshalIndent(val, "", "    ")
	reader := bytes.NewReader(b)
	if err != nil {
		return EmptyMessage(process), CommonError(process, err)
	}
	f := drive.File{MimeType: w.MimeType, Name: w.Name}
	_, err = driveService.Files.Update(w.Id, &f).Media(reader).Do()

	if err != nil {
		fmt.Println(err.Error())
		return EmptyMessage(process), CommonError(process, err)
	}
	fmt.Println("Successfully updated file " + w.Id)
	db.ClearDriveUpdate(w.Id)
	return CommonMessage(process, "success"), CommonError(process, nil)
}
