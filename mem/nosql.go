package mem

import (
	"encoding/json"
	"fmt"
	"fold/console"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type NoSql struct {
	File       string
	document   map[string]any
	collection []*NoSql
	data       *Data
	is         string
}

func (n *NoSql) Val() any {
	if n.is == Array {
		return n.CollectionVal()
	}
	if n.is == Struct {
		return n.document
	}
	if n.data == nil {
		return nil
	}

	return n.data.Val()
}

func (n *NoSql) Str() string {
	if n.data == nil {
		return ""
	}
	return n.data.Str()
}

func ConvertIndex[T any](key string, arr []T) int {
	i, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return -1
	}
	if i < 0 || int(i) >= len(arr) {
		return -1
	}
	return int(i)
}

func (n *NoSql) Raw(key string) (any, bool) {
	if n.is == Array {
		i := ConvertIndex(key, n.collection)
		if i < 0 {
			return nil, false
		}
		return n.collection[i], true
	}
	if n.is == Struct {
		val, ok := n.document[key]
		return val, ok
	}
	return nil, false
}

func (n *NoSql) Get(key string) *NoSql {
	raw, ok := n.Raw(key)
	if !ok {
		return nil
	}
	return FromAny(raw)
}

func FromAnyArray(val []any) []*NoSql {
	res := make([]*NoSql, len(val))
	for i, v := range val {
		res[i] = FromAny(v)
	}
	return res
}

func FromAny(val any) *NoSql {
	fmt.Println(val)
	switch x := val.(type) {
	case *NoSql:
		return x
	case []any:
		fmt.Println("collection")
		return &NoSql{collection: FromAnyArray(x), is: Array}
	case map[string]any:
		fmt.Println("json")
		return &NoSql{document: x, is: Struct}
	default:
		fmt.Println("string")
		fmt.Println(reflect.TypeOf(x))
		str := fmt.Sprint(x)
		data := FromString(str)
		return &NoSql{data: data, is: data.is}
	}
}

func (n *NoSql) Set(key string, val any) bool {
	nn := FromAny(val)
	if n.is == Array {
		i := ConvertIndex(key, n.collection)
		if i < 0 {
			return false
		}
		n.collection[i] = nn
	}

	n.document[key] = nn.Val()
	return true
}

func LoadJson(file string) (*NoSql, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln("Error closing json " + file)
		}
	}(f)
	raw, _ := io.ReadAll(f)
	var data any
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	res := FromAny(data)
	res.File = file
	fmt.Println("------")
	fmt.Println(res.CollectionVal())
	return res, nil
}

func (n *NoSql) DeepGet(path string) (*NoSql, string, *NoSql, bool) {
	parts := strings.Split(path, ".")
	fmt.Println(n)
	fmt.Println(parts[0])
	var innerContainer = n.Get(parts[0])
	fmt.Println(innerContainer)
	if len(parts) < 2 {
		return n, parts[0], innerContainer, true
	}
	var value *NoSql
	for i := 1; i < len(parts)-1; i++ {
		value = innerContainer.Get(parts[i])
		if value == nil || (value.is != Array && value.is != Struct) {
			return innerContainer, parts[i], value, false
		}
		innerContainer = value
	}
	lastPart := parts[len(parts)-1]
	return innerContainer, lastPart, innerContainer.Get(lastPart), true
}

func (n *NoSql) RawSearch(query *map[string][]string) any {
	if n.is != Array || query == nil || len(*query) == 0 {
		console.YellowPrintln("taking raw val")
		return n.Val()
	}

	filtered, _, _, _ := n.CollectionSearch(query)
	res := make([]any, len(filtered))
	for i, item := range filtered {
		res[i] = item.Val()
	}
	return res
}

func Matches(n *NoSql, s string) bool {
	fmt.Println("compare " + s)
	fmt.Println(n)
	if n == nil {
		return s == "" || s == "null"
	}
	fmt.Println(n.Str())
	return n.Str() == s
}

func (n *NoSql) Matches(query *map[string][]string) (*NoSql, string, *NoSql, bool) {
	var container, value *NoSql
	field := ""
	ok := false
	for path, values := range *query {
		container, field, value, ok = n.DeepGet(path)
		if !ok {
			return container, field, value, false
		}
		anyValueMatches := false
		for _, v := range values {
			if Matches(value, v) {
				anyValueMatches = true
				break
			}
		}
		if !anyValueMatches {
			return container, field, value, false
		}
	}
	return container, field, value, true
}

func (n *NoSql) CollectionSearch(query *map[string][]string) ([]*NoSql, []*NoSql, string, []*NoSql) {
	all := n.collection
	res := make([]*NoSql, 0, len(all))
	containers := make([]*NoSql, 0, len(all))
	values := make([]*NoSql, 0, len(all))
	field := ""
	for _, item := range all {
		container, f, value, ok := item.Matches(query)
		if ok {
			res = append(res, item)
			containers = append(containers, container)
			values = append(values, value)
			field = f
		}
	}
	return res, containers, field, values
}

func (n *NoSql) CollectionVal() []any {
	length := 0
	if n.collection != nil {
		length = len(n.collection)
	}
	res := make([]any, length)
	for i := 0; i < length; i++ {
		res[i] = n.collection[i].Val()
	}

	return res
}

func (n *NoSql) Patch(query *map[string][]string, update any) any {
	if n.is != Array || query == nil || len(*query) == 0 {
		console.YellowPrintln("taking raw val")
		return n.Val()
	}

	filtered, containers, path, _ := n.CollectionSearch(query)
	res := make([]any, len(filtered))
	for i, item := range filtered {
		containers[i].Set(path, update)
		res[i] = item.Val()
	}
	return res
}
