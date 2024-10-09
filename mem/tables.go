package mem

import (
	"fmt"
	"fold/console"
)

type Index = map[string][]Data
type Indexes map[string]Index
type Table struct {
	name           string
	indexes        Indexes
	rows           [][]Data
	cols           []*ColumnDefinition
	primaryIndex   string
	foreignIndexes []*ColumnDefinition
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
	return t.indexes[t.primaryIndex][id]
}

func (t Table) MapRow(row []Data) map[string]string {
	res := make(map[string]string)
	for index, value := range row {
		res[t.cols[index].name] = value.Str()
	}
	return res
}

func (t Table) MapJoinRow(row []Data, store *Store, tablePathMap *map[string]bool) map[string]any {
	res := make(map[string]any)
	for index, value := range row {
		res[t.cols[index].name] = value.Val()
	}
	console.YellowPrintln("Map join on table: " + t.name)
	pathMap := *tablePathMap
	for _, column := range t.foreignIndexes {
		console.YellowPrintln(fmt.Sprintf("Checking foreign index: %s->%s ", column.foreignTable, column.foreignColumn))
		_, ok := pathMap[column.foreignTable]
		if ok {
			continue
		}
		pathMap[column.foreignTable] = true

		val := row[column.number]
		joinTable := store.GetTable(column.foreignTable)
		joinRow := joinTable.GetRowByIndex(column.foreignColumn, val.Str())
		res[column.foreignTable] = joinTable.MapJoinRow(joinRow, store, tablePathMap)
	}
	return res
}

func (t Table) Get(id string, store *Store) map[string]any {
	pathMap := make(map[string]bool)
	return t.MapJoinRow(t.GetRow(id), store, &pathMap)
}

func (t Table) All() []map[string]string {
	res := make([]map[string]string, len(t.rows))
	for index, row := range t.rows {
		res[index] = t.MapRow(row)
	}
	return res
}

func InitTable(indexes Indexes, cols []*ColumnDefinition, nColumns int, nRows int, primaryIndex string) *Table {
	a := make([][]Data, nRows)
	for i := range a {
		a[i] = make([]Data, nColumns)
	}
	return &Table{indexes: indexes, rows: a, cols: cols, primaryIndex: primaryIndex, foreignIndexes: make([]*ColumnDefinition, 0)}
}

type ColumnDefinition struct {
	name          string
	isIndex       bool
	isUnique      bool
	foreignTable  string
	foreignColumn string
	foreignUnique bool
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

func ColumnsPrintln(columns []*ColumnDefinition) {
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
