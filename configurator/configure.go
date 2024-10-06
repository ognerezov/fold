package configurator

import (
	"fmt"
	"fold/csv"
	"fold/mem"
	"fold/path"
	goji "goji.io"
	"io/fs"
	"path/filepath"
	"strings"
)

func Configure(dataPath string) (*goji.Mux, error) {
	mux := goji.NewMux()
	store := *mem.TheStore
	clean := path.CreateRootCleaner(dataPath)
	var err = path.ProcessPath(dataPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		var route = "/" + clean(filepath.Dir(path))
		var name = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

		if name != "index" {
			if strings.HasSuffix(route, "/") {
				route = fmt.Sprintf("%s%s", route, name)
			} else {
				route = fmt.Sprintf("%s/%s", route, name)
			}
		}
		var filename = fmt.Sprintf("%s/%s", dataPath, clean(path))
		switch filepath.Ext(path) {
		case ".csv":
			records := csv.ReadCsvFile(filename)
			table := *mem.TableFromRecords(records)
			store.SetTable(route, table)
			SetCSVHandlers(route, mux)
		}

		return nil
	})
	store.ReIndex()
	if err != nil {
		return nil, err
	}

	return mux, nil
}
