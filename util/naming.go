package util

import (
	"fmt"
	"strings"
)

func NamingLookups(name string) (bool, []string, []string, error) {
	var literals = strings.Split(name, HeaderDelimiter)
	lenLiterals := len(literals)
	if lenLiterals < 2 {
		return false, nil, nil, nil
	}
	tables := make([]string, lenLiterals-1)
	cols := make([]string, lenLiterals-1)

	for split := 1; split < len(literals); split++ {
		head, tail := SplitArray(literals, split)
		tables[split-1] = strings.Join(head, HeaderDelimiter)
		cols[split-1] = strings.Join(tail, HeaderDelimiter)
		fmt.Printf("Found variant %s -> %s", tables[split-1], cols[split-1])
		fmt.Println()
	}

	return true, tables, cols, nil
}

func TableToPath(table string) string {
	return "/" + strings.Replace(table, "_", "/", -1)
}

func SplitArray[t any](arr []t, split int) ([]t, []t) {
	totalLen := len(arr)
	tailLen := totalLen - split

	head := make([]t, split)
	tail := make([]t, tailLen)

	for i := 0; i < split; i++ {
		head[i] = arr[i]
	}

	for i := split; i < totalLen; i++ {
		tail[i-split] = arr[i]
	}

	return head, tail
}
