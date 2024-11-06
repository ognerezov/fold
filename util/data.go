package util

import (
	"maps"
	"strconv"
)

func Flatten(prefix string, src *map[string]any, dest *map[string]any) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range *src {
		switch child := v.(type) {
		case map[string]any:
			Flatten(prefix+k, &child, dest)
		case []any:
			for i := 0; i < len(child); i++ {
				(*dest)[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			(*dest)[prefix+k] = v
		}
	}
}

func MergeMaps(dest map[string]any, src map[string]any) {
	for k, v := range dest {
		switch child := v.(type) {
		case map[string]any:
			vv, ok := src[k]
			if ok {
				switch cc := vv.(type) {
				case map[string]any:
					MergeMaps(cc, child)
				}
			}
		}
	}

	maps.Copy(dest, src)
}
