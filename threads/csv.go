package threads

import "fold/csv"

func WriteCsvAsync(filePath string, records [][]string) {
	writer := Writer(filePath)
	Async(writer, records)
}

type Writer string

func (w Writer) Call(args [][]string) (Message[string], ErrorMessage) {
	e := csv.WriteCsvFile(string(w), args)
	process := "Save csv to file " + string(w)
	if e != nil {
		return EmptyMessage(process), CommonError(process, e)
	}
	return CommonMessage(process, "success"), CommonError(process, e)
}
