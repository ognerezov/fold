package csv

import (
	"encoding/csv"
	"fold/util"
	"log"
	"os"
)

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer util.CloseFie(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func WriteCsvFile(filePath string, records [][]string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal("Unable to read output file "+filePath, err)
	}
	defer util.CloseFie(f)
	w := csv.NewWriter(f)
	return w.WriteAll(records)
}
