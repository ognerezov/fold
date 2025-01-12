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
	"google.golang.org/api/drive/v3"
	"net/http"
	"strconv"
	"strings"
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
	for _, file := range r.Files {
		fileId := file.Id
		fileName := file.Name
		m, ext, _ := GetContentType(fileName)
		fileRoute := basePath + "/" + fileName
		bytes, err := FetchDriveFile(fileId)
		if err != nil {
			console.RedPrintln(err.Error())
			continue
		}
		if ext == ".json" {
			SetJsonDriveHandlers(file, fileRoute, m, bytes, mux, api)
			continue
		}
		SetRawDriveHandlers(fileRoute, fileId, m, bytes, mux, api)
	}
}

func SetJsonDriveHandlers(file *drive.File, fileRoute string, mime string, bytes []byte, mux *goji.Mux, api *openapi.ApiDescription) {
	noSql, err := mem.FromBytes(bytes)
	if err != nil {
		console.RedPrintln(err.Error())
		return
	}
	route := strings.TrimSuffix(fileRoute, ".json")
	if file.MimeType == "" {
		file.MimeType = mime
	}
	noSql.DriveFile = file
	store := mem.TheStore
	store.SetNoSql(route, noSql)
	SetJsonHandlers(route, mux, api)
}

func SetRawDriveHandlers(fileRoute string, fileId string, mime string, bytes []byte, mux *goji.Mux, api *openapi.ApiDescription) {
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
		w.Header().Set("Content-Type", mime)

		_, err := w.Write(bytes)
		if err == nil {
			return
		}
		console.RedPrintln(err.Error())
		router.ServerError(err, w)
	})
}
