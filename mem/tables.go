package mem

import "fmt"

type Index = map[string][]Data
type Indexes map[string]Index
type Table struct {
	indexes        Indexes
	rows           [][]Data
	cols           []ColumnDefinition
	primaryIndex   string
	foreignIndexes []ColumnDefinition
}

func (t Table) Print() {
	ColumnsPrintln(t.cols)
	for index, row := range t.rows {
		fmt.Println(index, row)
	}
}

func (t Table) GetRowByIndex(col string, id string) []Data {
	return t.indexes[col][id]
}

func (t Table) GetRow(id string) []Data {
	row := t.indexes[t.primaryIndex][id]
	//store := *TheStore
	//for index, column := range t.foreignIndexes {
	//	val := row[column.number]
	//	join := store.GetTable(column.foreignTable).GetRowByIndex(column.foreignColumn, string(val))
	//}

	return row
}

func (t Table) MapRow(row []Data) map[string]string {
	res := make(map[string]string)
	for index, value := range row {
		res[t.cols[index].name] = value.Str()
	}
	return res
}

func (t Table) Get(id string) map[string]string {
	return t.MapRow(t.GetRow(id))
}

func (t Table) All() []map[string]string {
	res := make([]map[string]string, len(t.rows))
	for index, row := range t.rows {
		res[index] = t.MapRow(row)
	}
	return res
}

func InitTable(indexes Indexes, cols []ColumnDefinition, nColumns int, nRows int, primaryIndex string) *Table {
	a := make([][]Data, nRows)
	for i := range a {
		a[i] = make([]Data, nColumns)
	}
	return &Table{indexes: indexes, rows: a, cols: cols, primaryIndex: primaryIndex, foreignIndexes: make([]ColumnDefinition, 0)}
}

type ColumnDefinition struct {
	name          string
	isIndex       bool
	isUnique      bool
	foreignTable  string
	foreignColumn string
	dataType      string
	number        int
}

func (c ColumnDefinition) ToString() string {
	if c.isIndex {
		return fmt.Sprintf("[%s*]", c.name)
	}
	if c.isIndex {
		return fmt.Sprintf("%s*", c.name)
	}
	return c.name
}

func ColumnsPrintln(columns []ColumnDefinition) {
	fmt.Print("_ |")
	for _, column := range columns {
		fmt.Print(column.ToString() + " |")
	}
	fmt.Println()
}

func SimpleDefinition(name string, isIndex bool, index int) *ColumnDefinition {
	return &ColumnDefinition{name: name, isIndex: isIndex, isUnique: isIndex, number: index}
}

func (c ColumnDefinition) IsIndex() bool {
	return c.isIndex
}
