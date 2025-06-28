package util

import (
	"encoding/json"
	"fmt"
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

type ValueTransformer func(any) any

func ReplaceValue(dest *map[string]any, key string, f ValueTransformer) {
	for k, v := range *dest {
		if k == key {
			(*dest)[k] = f(v)
			continue
		}
		switch child := v.(type) {
		case map[string]any:
			ReplaceValue(&child, key, f)
		case []any:
			for _, element := range child {
				switch ch := element.(type) {
				case map[string]any:
					ReplaceValue(&ch, key, f)
				}
			}
		default:
			fmt.Println(k, child)
		}
	}
}

func Restructure[T any](data any, out *T) error {
	h, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(h, out)
}

func RestructureArrays(values [][]any, factor int) [][]any {
	if values == nil || len(values) == 0 {
		return nil
	}

	result := make([][]any, factor)

	if factor >= len(values) {
		j := 0
		for i := 0; i < factor; i++ {
			result[i] = values[j]
			j++
			if j == len(values) {
				j = 0
			}
		}
		return result
	}
	i := 0
	for _, v := range values {

		if result[i] == nil {
			result[i] = v
			continue
		}
		result[i] = append(result[i], v...)
		i++
		if i == factor {
			i = 0
		}
	}
	return result

}
