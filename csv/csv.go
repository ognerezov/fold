package csv

import (
	"encoding/csv"
	"fmt"
	"fold/console"
	"fold/threads"
	"log"
	"os"
)

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func WriteCsvFile(filePath string, records [][]string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		log.Fatal("Unable to read output file "+filePath, err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(f)
	w := csv.NewWriter(f)
	return w.WriteAll(records)
}

func WriteCsvAsync(filePath string, records [][]string) {
	fmt.Println(filePath)
	fmt.Println(records)
	writer := Writer(filePath)
	threads.Async(writer, records)
}

type Writer string

func (w Writer) Call(args [][]string) (threads.Message[string], threads.ErrorMessage) {
	e := WriteCsvFile(string(w), args)
	process := "Save csv to file " + string(w)
	if e != nil {
		return threads.EmptyMessage(process), threads.CommonError(process, e)
	}
	return threads.CommonMessage(process, "success"), threads.CommonError(process, e)
}
