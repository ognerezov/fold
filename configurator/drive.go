package configurator

import (
	"errors"
	"fmt"
	"fold/console"
	"fold/csv"
	"fold/interfaces"
	"fold/mem"
	"fold/migrations"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	SpreadSheetMime = "application/vnd.google-apps.spreadsheet"
)

var (
	migrationHandlers = make(map[string]*migrations.DriverHandler)
)

type DriveFileHandler drive.File

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

func GetSpreadSheet(fileId string) (*sheets.Spreadsheet, error) {
	sheetsService := AppProviders.Google.Sheets
	if sheetsService == nil {
		return nil, errors.New("sheet service not found in configurator")
	}

	spreadSheet, err := sheetsService.Spreadsheets.Get(fileId).Do()
	if err != nil {
		return nil, err
	}
	return spreadSheet, nil
}

func FetchSpreadSheet(fileId string) ([][]string, *sheets.Spreadsheet, error) {
	driveService := AppProviders.Google.Drive
	if driveService == nil {
		return nil, nil, errors.New("drive service not found in configurator")
	}

	resp, err := driveService.Files.Export(fileId, "text/csv").Download()

	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("response status code is " + strconv.Itoa(resp.StatusCode))
	}
	bytes, err := util.HandleBody(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	records, err := csv.BytesToCsv(bytes)
	if err != nil {
		return nil, nil, err
	}

	spreadSheet, err := GetSpreadSheet(fileId)
	if err != nil {
		return nil, nil, err
	}
	return records, spreadSheet, nil
}

func (fid DriveFileHandler) Fetch() ([]byte, error) {
	return FetchDriveFile(fid.Id)
}

func (fid DriveFileHandler) FetchNoSql() (*mem.NoSql, error) {
	bytes, err := fid.Fetch()
	if err != nil {
		return nil, err
	}
	noSql, err := mem.FromBytes(bytes)
	noSql.DriveFile = fid.P()
	if err != nil {
		return nil, err
	}

	return noSql, nil
}

func (fid DriveFileHandler) FetchCsv() (*mem.Table, error) {
	records, spreadSheet, err := FetchSpreadSheet(fid.Id)
	if err != nil {
		return nil, err
	}
	table := mem.TableFromRecords(records)
	table.Spreadsheet = spreadSheet
	table.DriveFile = fid.P()
	return table, nil
}

func (fid DriveFileHandler) P() *drive.File {
	var res drive.File
	res = drive.File(fid)
	return &res
}

func SetDriveHandlers(basePath string, id string, mux *goji.Mux, api *openapi.ApiDescription, next *interfaces.Phase, mHandler *migrations.DriverHandler, endpoints Endpoints) {
	files, err := AppProviders.Google.ReadDir(id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, file := range files {
		fileId := file.Id
		fileName := file.Name
		// Will mount this files later
		if slices.Contains(doNotServe, fileName) {
			continue
		}
		fileHandler := DriveFileHandler(*file)
		m, ext, hasExt := GetContentType(fileName)
		fileRoute := basePath
		if strings.TrimSuffix(fileName, ext) != "index" {
			if strings.HasSuffix(basePath, "/") {
				fileRoute = fmt.Sprintf("%s%s", basePath, fileName)
			} else {
				fileRoute = fmt.Sprintf("%s/%s", basePath, fileName)
			}
		}
		fmt.Println(basePath)
		fmt.Println(fileRoute)
		fileRoute = util.TableToPath(fileRoute)
		route := fileRoute
		if hasExt {
			route = strings.TrimSuffix(fileRoute, ext)
		}
		if file.MimeType == SpreadSheetMime {
			console.YellowPrintln("File is a Spreadsheet")
			table, err := fileHandler.FetchCsv()

			if err != nil {
				console.RedPrintln(err.Error())
			} else {
				store := mem.TheStore
				store.SetTable(route, table, fileHandler)
				next.Append(SetTableHandlers(route, mux, api))
			}
			continue
		}

		if ext == ".json" {
			SetJsonDriveHandlers(route, fileHandler, mux, api)
			continue
		}
		if ext == ".fold" {
			console.GreenPrintln("Registering fold action handler " + fileName)
			SetControlHandlers(route, fileName, mux, api, endpoints)
		}
		if file.MimeType == "application/vnd.google-apps.folder" {
			console.YellowPrintln("File is a nested folder " + fileName)
			migrationHandlers[id] = migrations.CreateDriveHandler(AppProviders.Google, id)
			folderId := file.Id
			mHandler.RegisterFolder(&folderId, fileName)
			SetDriveHandlers(fileRoute, file.Id, mux, api, next, mHandler, endpoints)
			continue
		}
		bytes, err := FetchDriveFile(fileId)
		if err != nil {
			console.RedPrintln(err.Error())
			continue
		}
		if ext == ".html" {
			SetRawDriveHandlers(route, fileHandler, m, bytes, mux, api)
			continue
		}
		SetRawDriveHandlers(fileRoute, fileHandler, m, bytes, mux, api)
	}
}

func SetJsonDriveHandlers(fileRoute string, fileIdFetcher DriveFileHandler, mux *goji.Mux, api *openapi.ApiDescription) {
	noSql, err := fileIdFetcher.FetchNoSql()
	if err != nil {
		console.RedPrintln(err.Error())
		return
	}
	route := strings.TrimSuffix(fileRoute, ".json")
	store := mem.TheStore
	store.SetNoSql(route, noSql, fileIdFetcher)
	SetJsonHandlers(route, mux, api)
}

func SetRawDriveHandlers(fileRoute string, fileIdFetcher DriveFileHandler, mime string, bytes []byte, mux *goji.Mux, api *openapi.ApiDescription) {
	console.CyanPrintln("Registering GET " + fileRoute)
	// Drive data is always cached
	mem.TheStore.Cache(fileRoute, bytes, fileIdFetcher)
	api.Path(fileRoute).Get(openapi.Method{
		Summary: "GET " + fileRoute + " raw file fetcher",
		Responses: map[string]openapi.Response{
			"200": {
				Description: "Raw data",
				Content: map[string]openapi.Content{
					mime: {
						Schema: openapi.Binary,
					},
				},
			},
		},
	})
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
