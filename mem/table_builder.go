package mem

import "fold/csv"

func ReadHeader(header []string) ([]*ColumnDefinition, Indexes, string) {
	res := make([]*ColumnDefinition, len(header))
	var indexes = Indexes{}
	var primaryIndex = ""
	for index, element := range header {
		res[index] = GetHeaderDefinition(element, index)
		if res[index].IsIndex() {
			indexes[element] = Index{}
			if primaryIndex == "" {
				primaryIndex = element
			}
		}
	}

	return res, indexes, primaryIndex
}

func TableFromRecords(records [][]string) *Table {
	columns, indexes, primaryIndex := ReadHeader(records[0])
	nCols := len(records[0])
	nRows := len(records) - 1
	table := InitTable(indexes, columns, nCols, nRows, primaryIndex)
	for rIndex, record := range records[1:] {
		table.rows[rIndex] = make([]Data, nCols)
		for index, element := range record {
			column := columns[index]
			data := FromString(element)
			column.SetDataType(data.is)
			table.rows[rIndex][index] = *data
			if column.isIndex {
				indexes[column.name][element] = table.rows[rIndex]
			}
		}
	}

	return table
}

func IsId(column string) bool {
	return column == "id"
}

func GetHeaderDefinition(name string, index int) *ColumnDefinition {
	var isId = IsId(name)
	return SimpleDefinition(name, isId, index)
}

func TableFromBytes(bytes []byte) (*Table, error) {
	records, err := csv.BytesToCsv(bytes)
	if err != nil {
		return nil, err
	}

	return TableFromRecords(records), nil

}
