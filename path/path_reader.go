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

func WalkPath(root string) error {
	var rootPrefix = pathLeadingSymbols.ReplaceAllString(root, "")
	fmt.Println("Scanning root path " + root)
	var err = filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {

		fmt.Println("Scanning " + strings.Replace(path, rootPrefix+"/", "", -1))
		fmt.Println(info.IsDir())
		fmt.Println(info.Name())
		return nil
	})

	return err
}
