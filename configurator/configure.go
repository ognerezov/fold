package configurator

import (
	"fmt"
	"fold/path"
	goji "goji.io"
	"io/fs"
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
		var filename = fmt.Sprintf("%s/%s", dataPath, clean(path))
		switch filepath.Ext(path) {
		case ".csv":
			SetCSVHandlers(route, filename, mux)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return mux
}
