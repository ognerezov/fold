package configurator

import (
	"fold/console"
	"fold/router"
	goji "goji.io"
	"goji.io/pat"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func GetContentType(filePath string) (string, string, bool) {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return "", "", false
	}
	m := mime.TypeByExtension(ext)
	return m, ext, true
}

func SetRawHandlers(route string, filePath string, mux *goji.Mux) {
	m, ext, hasExt := GetContentType(filePath)
	console.BluePrintln("Registering GET " + route + ext)
	mux.HandleFunc(pat.Get(route+ext), func(w http.ResponseWriter, r *http.Request) {
		f, err := os.OpenFile(filePath, os.O_RDONLY, 0)
		if err != nil {
			if os.IsNotExist(err) {
				router.NotFound(w)
				return
			}
			router.ServerError(err, w)
			return
		}
		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				console.RedPrintln(err.Error())
			}
		}(f)
		bytes, err := io.ReadAll(f)
		if err != nil {
			console.RedPrintln(err.Error())
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
		console.RedPrintln(err.Error())
		router.ServerError(err, w)
	})
}
