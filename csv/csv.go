package csv

import (
	"bytes"
	"encoding/csv"
	"fold/util"
	"log"
	"os"
)

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer util.CloseFie(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
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

func BytesToCsv(b []byte) ([][]string, error) {

	csvReader := csv.NewReader(bytes.NewReader(b))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
