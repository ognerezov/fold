package mem

import "fmt"

type Index = map[string][]any
type Indexes map[string]Index
type Table struct {
	indexes      Indexes
	rows         [][]any
	cols         []ColumnDefinition
	primaryIndex string
}

func (t Table) Print() {
	ColumnsPrintln(t.cols)
	for index, row := range t.rows {
		fmt.Println(index, row)
	}
}

func (t Table) GetRow(id string) []any {
	return t.indexes[t.primaryIndex][id]
}

func (t Table) MapRow(row []any) map[string]any {
	res := make(map[string]any)
	for index, value := range row {
		res[t.cols[index].name] = value
	}
	return res
}

func (t Table) Get(id string) map[string]any {
	return t.MapRow(t.GetRow(id))
}

func (t Table) All() []map[string]any {
	res := make([]map[string]any, len(t.rows))
	for index, row := range t.rows {
		res[index] = t.MapRow(row)
	}
	return res
}

func InitTable(indexes Indexes, cols []ColumnDefinition, nColumns int, nRows int, primaryIndex string) *Table {
	a := make([][]any, nRows)
	for i := range a {
		a[i] = make([]any, nColumns)
	}
	return &Table{indexes: indexes, rows: a, cols: cols, primaryIndex: primaryIndex}
}

type ColumnDefinition struct {
	name          string
	isIndex       bool
	isUnique      bool
	foreignTable  string
	foreignColumn string
	dataType      string
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

func SimpleDefinition(name string, isIndex bool) *ColumnDefinition {
	return &ColumnDefinition{name: name, isIndex: isIndex, isUnique: isIndex}
}

func (c ColumnDefinition) IsIndex() bool {
	return c.isIndex
}
