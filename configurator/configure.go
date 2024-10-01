package configurator

import (
	"fmt"
	"fold/path"
	goji "goji.io"
	"goji.io/pat"
	"io/fs"
	"net/http"
	"path/filepath"
)

func Configure(dataPath string) *goji.Mux {
	var mux = goji.NewMux()
	clean := path.CreateRootCleaner(dataPath)
	var err = path.ProcessPath(dataPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		var route = "/" + clean(filepath.Dir(path))
		fmt.Println("Registering GET " + route)
		mux.HandleFunc(pat.Get(route), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Find file, %s!", info.Name())
		})
		return nil
	})
	if err != nil {
		return nil
	}

	return mux
}
