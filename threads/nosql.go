package threads

import (
	"fmt"
	"fold/db"
)

func WriteNosqlAsync(filePath string, val any) {
	writer := AnyWriter(filePath)
	Async(writer, val)
}

type AnyWriter string

func (w AnyWriter) Call(val any) (Message[string], ErrorMessage) {
	e := db.SaveJson(string(w), val)
	process := "SaveJson data to file " + string(w)
	if e != nil {
		fmt.Println(e.Error())
		return EmptyMessage(process), CommonError(process, e)
	}
	db.ClearFileUpdate(string(w))
	return CommonMessage(process, "success"), CommonError(process, e)
}
