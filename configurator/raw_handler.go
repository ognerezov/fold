package configurator

import (
	"fmt"
	"fold/console"
	"fold/mem"
	"fold/openapi"
	"fold/router"
	"fold/util"
	goji "goji.io"
	"goji.io/pat"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var (
	FrontendFiles = map[string]bool{
		".html": true,
	}
)

func GetContentType(filePath string) (string, string, bool) {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return "", "", false
	}
	m := mime.TypeByExtension(ext)
	return m, ext, true
}

func SetRawHandlers(route string, filePath string, mux *goji.Mux, api *openapi.ApiDescription) {
	m, ext, hasExt := GetContentType(filePath)

	if hasExt && FrontendFiles[ext] {
		console.BluePrintln("Registering GET " + route)
		bytes, err := util.ReadFile(filePath)
		if err != nil {
			return
		}
		mem.TheStore.Cache(route, bytes)
		mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
			console.BluePrintln("Searching cache for " + route)
			var ok bool
			bytes, ok = mem.TheStore.GetCached(route)
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
		api.DescribeRawGet(route, "Get cached frontend file", m)
		return
	}
	console.BluePrintln("Registering GET " + route + ext)
	mux.HandleFunc(pat.Get(route+ext), func(w http.ResponseWriter, r *http.Request) {
		console.BluePrintln("Searching filesystem for " + route)
		fmt.Println(filePath)
		bytes, err := util.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				router.NotFound(w)
				return
			}
			router.ServerError(err, w)
			return
		}
		if hasExt {
			w.Header().Set("Content-Type", m)
		}
		_, err = w.Write(bytes)
		if err == nil {
			return
		}
		api.DescribeRawGet(route, "Get raw file from hdd", m)
		console.RedPrintln(err.Error())
		router.ServerError(err, w)
	})
}
