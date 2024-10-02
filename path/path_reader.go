package path

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	pathLeadingSymbols = regexp.MustCompile(`^./|^/`)
)

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

func WalkPath(root string) error {
	clean := CreateRootCleaner(root)

	return ProcessPath(root, func(path string, info os.FileInfo, _ error) error {

		fmt.Println("Scanning " + clean(filepath.Dir(path)))
		fmt.Println(info.IsDir())
		fmt.Println("Dir: " + filepath.Dir(path))
		fmt.Println("Ext: " + filepath.Ext(path))
		fmt.Println(info.Name())
		fmt.Printf("%s/%s", root, clean(path))
		fmt.Println("")
		fmt.Println("---------")
		return nil
	})
}

func ProcessPath(root string, f filepath.WalkFunc) error {
	fmt.Println("Scanning root path " + root)
	var err = filepath.Walk(root, f)
	return err
}
