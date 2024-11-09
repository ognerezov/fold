package path

import (
	"fold/console"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	pathLeadingSymbols = regexp.MustCompile(`^./|^/`)
	Root               string
)

func FullPath(path string) string {
	return filepath.Join(Root, path)
}

type DirMapper func(path string) string

func CreateRootCleaner(root string) DirMapper {
	var rootPrefix = pathLeadingSymbols.ReplaceAllString(root, "")

	return func(path string) string {
		if path == rootPrefix {
			return ""
		}
		return strings.Replace(path, rootPrefix+"/", "", -1)
	}
}

func ProcessPath(root string, f filepath.WalkFunc) error {
	console.GreenPrintln("Scanning root path " + root)
	var err = filepath.Walk(root, f)
	return err
}
