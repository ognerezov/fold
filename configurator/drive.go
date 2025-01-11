package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"io"
	"net/http"
)

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
		fileRoute := basePath + "/" + fileName
		console.CyanPrintln("Registering GET " + fileRoute)
		mux.HandleFunc(pat.Get(fileRoute), func(w http.ResponseWriter, r *http.Request) {
			console.BluePrintln("Incoming request GET " + fileRoute)
			resp, e := driveService.Files.Get(fileId).Download()
			if e != nil {
				console.RedPrintln(e.Error())

				router.ServerError(err, w)
				return
			}
			if resp.StatusCode != http.StatusOK {
				decoder := json.NewDecoder(resp.Body)
				var errResp map[string]any
				err = decoder.Decode(&errResp)
				err = errors.New(resp.Status)
				router.ReturnError(err, resp.StatusCode, w)
				return
			}
			body := resp.Body
			defer util.HideBody(body)
			_, err = io.Copy(w, resp.Body)

			if err != nil {
				router.ServerError(err, w)
			}
		})
	}
}
