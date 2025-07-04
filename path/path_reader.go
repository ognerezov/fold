package path

import (
	"fmt"
	"fold/console"
	"path/filepath"
	"regexp"
)

var (
	pathLeadingSymbols = regexp.MustCompile(`^./`)
	Root               string
)

type RootCleaner func(path string) string

func CreateRootCleaner(root string) RootCleaner {
	rootPrefix := pathLeadingSymbols.ReplaceAllString(root, "")
	rootRegex := regexp.MustCompile(PrepareForRegex(fmt.Sprintf(`^%s%s`, rootPrefix, filepath.FromSlash("/"))))
	return func(path string) string {
		if path == rootPrefix {
			return ""
		}
		return rootRegex.ReplaceAllString(path, "")
	}
}

func ProcessPath(root string, f filepath.WalkFunc) error {
	console.GreenPrintln("Scanning root path " + root)
	var err = filepath.Walk(root, f)
	return err
}
