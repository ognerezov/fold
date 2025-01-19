package configurator

import (
	"fold/csv"
	"fold/mem"
	"fold/util"
)

type FilePath string

func (f FilePath) Fetch() ([]byte, error) {
	return util.ReadFile(string(f))
}

func (f FilePath) FetchCsv() (*mem.Table, error) {
	filename := string(f)
	records, err := csv.ReadCsvFile(filename)
	if err != nil {
		return nil, err
	}
	table := mem.TableFromRecords(records)
	table.File = filename
	return table, nil
}

func (f FilePath) FetchNoSql() (*mem.NoSql, error) {
	filename := string(f)
	res, err := mem.LoadJson(filename)
	if err != nil {
		return nil, err
	}
	res.File = filename
	return res, nil
}
