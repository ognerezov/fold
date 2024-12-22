package threads

import (
	"fold/csv"
	"fold/db"
)

func WriteCsvAsync(filePath string, records [][]string) {
	writer := CsvWriter(filePath)
	Async(writer, records)
}

type CsvWriter string

func (w CsvWriter) Call(args [][]string) (Message[string], ErrorMessage) {
	e := csv.WriteCsvFile(string(w), args)
	process := "SaveJson csv to file " + string(w)
	if e != nil {
		return EmptyMessage(process), CommonError(process, e)
	}
	db.ClearTableUpdate(string(w))
	return CommonMessage(process, "success"), CommonError(process, e)
}
