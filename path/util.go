package path

import (
	"fmt"
	"fold/util"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
)

func Structure(fileName string, info fs.FileInfo, clean RootCleaner) (string, string, string) {
	var route = filepath.ToSlash("/" + clean(filepath.Dir(fileName)))
	var name = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

	if name != "index" {
		if strings.HasSuffix(route, "/") {
			route = fmt.Sprintf("%s%s", route, name)
		} else {
			route = fmt.Sprintf("%s/%s", route, name)
		}
	}
	route = util.TableToPath(route)
	return route, fileName, filepath.Ext(fileName)
}

func SubStructure(path string, info fs.FileInfo, clean RootCleaner) (string, string) {
	var name = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
	var route = filepath.ToSlash("/" + clean(filepath.Dir(path)))

	return route, name
}

func PrepareFileName(rootDataDir string, fileName string, clean RootCleaner) (string, string, string, string, string) {
	route := clean(filepath.Dir(fileName))
	filename := fmt.Sprintf("%s/%s", rootDataDir, clean(fileName))
	ext := filepath.Ext(fileName)
	m := ""
	if ext != "" {
		m = mime.TypeByExtension(ext)
	}
	return route, filename, filepath.Base(fileName), ext, m
}

func PrepareForRegex(str string) string {
	return strings.Replace(str, "\\", "\\\\", -1)
}
