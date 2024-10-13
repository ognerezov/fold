package mem

import (
	"fmt"
	"fold/console"
	"fold/util"
	"strconv"
)

const (
	MaxJoinDepth = 100
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

type JoinPathMap *map[string]int

func (t Table) DescribeColumns() []ColumnDefinition {
	result := make([]ColumnDefinition, len(t.cols))
	for i, col := range t.cols {
		result[i] = *col
	}
	return result
}

func (t Table) Print() {
	ColumnsPrintln(t.cols)
	for index, row := range t.rows {
		fmt.Println(index, row)
	}
}

func (t Table) InitPathTable() JoinPathMap {
	result := make(map[string]int)
	result[util.PathToTable(t.name)] = 0
	return &result
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

func (t Table) MapJoinRow(row []Data, store *Store, tablePathMap JoinPathMap, level int) map[string]any {
	res := make(map[string]any)
	for index, value := range row {
		res[t.cols[index].name] = value.Val()
	}
	if level >= MaxJoinDepth {
		console.RedPrintln("Maximum join depth exceeded " + strconv.Itoa(level))
		return res
	}
	console.YellowPrintln("Map join on table: " + t.name)
	pathMap := *tablePathMap
	for _, column := range t.foreignIndexes {
		console.YellowPrintln(fmt.Sprintf("Checking foreign index: %s->%s ", column.foreignTable, column.foreignColumn))
		previousLevel, ok := pathMap[column.foreignTable]

		/*
			We are storing path through tables to avoid circular joins
			If level == previousLevel we are querying multiple rows on current join
			Don't need to stop after first one
		*/
		if ok && previousLevel != level {
			continue
		}
		pathMap[column.foreignTable] = level

		val := row[column.number]
		joinTable, _ := store.GetTable(util.TableToPath(column.foreignTable))

		var joinRows [][]Data
		if column.foreignUnique {
			joinRows = make([][]Data, 0)
			joinRow := joinTable.GetRowByIndex(column.foreignColumn, val.Str())
			joinRows = append(joinRows, joinRow)
		} else {
			joinRows = joinTable.SearchRows(column.foreignColumn, val.Str())
		}
		joins := make([]map[string]any, len(joinRows))
		for index, joinRow := range joinRows {
			joins[index] = joinTable.MapJoinRow(joinRow, store, tablePathMap, level+1)
		}
		res[column.foreignTable] = joins
	}
	return res
}

func (t Table) Get(id string, store *Store) map[string]any {
	return t.MapJoinRow(t.GetRow(id), store, t.InitPathTable(), 0)
}

func (t Table) All() []map[string]string {
	res := make([]map[string]string, len(t.rows))
	for index, row := range t.rows {
		res[index] = t.MapRow(row)
	}
	return res
}

func (t Table) Search(query Query) any {
	if query.all {
		return t.All()
	}
	rows := make([]map[string]any, 0)
	for _, row := range t.rows {
		if !query.Matches(row) {
			continue
		}
		rows = append(rows, t.MapJoinRow(row, TheStore, t.InitPathTable(), 0))
	}
	return rows
}

func (t Table) QueryRows(query Query) []*[]Data {
	rows := make([]*[]Data, 0)
	for _, row := range t.rows {
		if !query.Matches(row) {
			continue
		}
		rows = append(rows, &row)
	}
	return rows
}

func (t Table) SearchRows(colName string, value string) [][]Data {
	rows := make([][]Data, 0)

	query := PrepareQuery(util.OneValueQuery(colName, value), &t)
	for _, row := range t.rows {
		if !query.Matches(row) {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func InitTable(indexes Indexes, cols []*ColumnDefinition, nColumns int, nRows int, primaryIndex string) *Table {
	a := make([][]Data, nRows)
	for i := range a {
		a[i] = make([]Data, nColumns)
	}
	return &Table{indexes: indexes, rows: a, cols: cols, primaryIndex: primaryIndex, foreignIndexes: make([]*ColumnDefinition, 0)}
}
