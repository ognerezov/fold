package mem

import (
	"fold/csv"
	"fold/util"
)

type FilePath string

func (f FilePath) Fetch() ([]byte, error) {
	return util.ReadFile(string(f))
}

func (f FilePath) FetchCsv() (*Table, error) {
	filename := string(f)
	records, err := csv.ReadCsvFile(filename)
	if err != nil {
		return nil, err
	}
	table := TableFromRecords(records)
	table.File = filename
	return table, nil
}

func (f FilePath) FetchNoSql() (*NoSql, error) {
	filename := string(f)
	res, err := LoadJson(filename)
	if err != nil {
		return nil, err
	}
	res.File = filename
	return res, nil
}
