package configurator

import (
	"fmt"
	"fold/arguments"
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

const (
	HTML = ".html"
)

var (
	FrontendFiles = map[string]bool{
		HTML:   true,
		".js":  true,
		".css": true,
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
	feRoute := route + ext
	if ext == HTML || !hasExt {
		feRoute = route
	}
	if hasExt && FrontendFiles[ext] && arguments.AppArguments.Cache {
		console.BluePrintln("Registering GET " + feRoute)
		bytes, err := util.ReadFile(filePath)
		if err != nil {
			return
		}
		mem.TheStore.Cache(feRoute, bytes, mem.FilePath(filePath))
		mux.HandleFunc(pat.Get(feRoute), func(w http.ResponseWriter, r *http.Request) {
			console.BluePrintln("Searching cache for " + feRoute)
			var ok bool
			bytes, ok = mem.TheStore.GetCached(feRoute)
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
		if api != nil {
			api.DescribeRawGet(feRoute, "Get cached frontend file", m)
		}
		return
	}
	console.BluePrintln("Registering GET " + feRoute)
	mux.HandleFunc(pat.Get(feRoute), func(w http.ResponseWriter, r *http.Request) {
		console.BluePrintln("Searching filesystem for " + feRoute)
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
		if api != nil {
			api.DescribeRawGet(feRoute, "Get raw file from hdd", m)
		}
		console.RedPrintln(err.Error())
		router.ServerError(err, w)
	})
}
