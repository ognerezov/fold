package configurator

import (
	"errors"
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
	"strconv"
)

type DriveFileId string

func FetchDriveFile(fileId string) ([]byte, error) {
	driveService := AppProviders.Google.Drive
	if driveService == nil {
		return nil, errors.New("drive service not found in configurator")
	}

	resp, err := driveService.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code is " + strconv.Itoa(resp.StatusCode))
	}
	bytes, err := util.HandleBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (fid DriveFileId) Fetch() ([]byte, error) {
	return FetchDriveFile(string(fid))
}

func SetDriveHandlers(basePath string, id string, mux *goji.Mux, api *openapi.ApiDescription) {
	fmt.Println(id)
	fmt.Println(basePath)
	driveService := AppProviders.Google.Drive
	if driveService == nil {
		console.RedPrintln("drive service not found in configurator")
		return
	}
	q := fmt.Sprintf("'%s' in parents", id)
	fmt.Println(q)
	r, err := driveService.Files.List().Q(q).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		PageSize(10).
		Fields("nextPageToken, files(id, name)").
		Do()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(r)
	for _, i := range r.Files {
		fileId := i.Id
		fileName := i.Name
		m, _, _ := GetContentType(fileName)
		fileRoute := basePath + "/" + fileName
		bytes, err := FetchDriveFile(fileId)
		if err != nil {
			console.RedPrintln(err.Error())
			continue
		}
		console.CyanPrintln("Registering GET " + fileRoute)
		// Drive data is always cached
		mem.TheStore.Cache(fileRoute, bytes, DriveFileId(fileId))
		mux.HandleFunc(pat.Get(fileRoute), func(w http.ResponseWriter, r *http.Request) {
			console.BluePrintln("Incoming request GET " + fileRoute)
			var ok bool
			bytes, ok = mem.TheStore.GetCached(fileRoute)
			if !ok {
				router.NotFound(w)
				return
			}
			w.Header().Set("Content-Type", m)

			_, err = w.Write(bytes)
			if err == nil {
				return
			}
			console.RedPrintln(err.Error())
			router.ServerError(err, w)
		})
	}
}
