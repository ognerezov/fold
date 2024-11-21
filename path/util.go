package path

import (
	"fmt"
	"fold/util"
	"io/fs"
	"path/filepath"
	"strings"
)

func Structure(dataPath string, path string, info fs.FileInfo, clean DirMapper) (string, string, string) {
	var route = "/" + clean(filepath.Dir(path))
	var name = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

	if name != "index" {
		if strings.HasSuffix(route, "/") {
			route = fmt.Sprintf("%s%s", route, name)
		} else {
			route = fmt.Sprintf("%s/%s", route, name)
		}
	}
	route = util.TableToPath(route)
	filename := fmt.Sprintf("%s/%s", dataPath, clean(path))
	return route, filename, filepath.Ext(path)
}
