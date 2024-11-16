package mem

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"fold/console"
	"fold/db"
	"fold/openapi"
	"fold/util"
	"github.com/google/uuid"
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
	primaryColumn  *ColumnDefinition
	File           string
}

type JoinPathMap *map[string]int

func (t *Table) DescribeColumns() []ColumnDefinition {
	result := make([]ColumnDefinition, len(t.cols))
	for i, col := range t.cols {
		result[i] = *col
	}
	return result
}

func (t *Table) Print() {
	ColumnsPrintln(t.cols)
	for index, row := range t.rows {
		fmt.Println(index, row)
	}
}

func (t *Table) InitPathTable() JoinPathMap {
	result := make(map[string]int)
	result[util.PathToTable(t.name)] = 0
	return &result
}

func (t *Table) GetRowByIndex(col string, id string) []Data {
	return t.indexes[col][id]
}

func (t *Table) GetRow(id string) []Data {
	return t.indexes[t.primaryIndex][id]
}

func (t *Table) MapRow(row []Data) map[string]any {
	res := make(map[string]any)
	for index, value := range row {
		res[t.cols[index].name] = value.Val()
	}
	return res
}

func (t *Table) MapJoinRow(row []Data, store *Store, tablePathMap JoinPathMap, level int) map[string]any {
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

func (t *Table) Get(id string, store *Store) map[string]any {
	return t.MapJoinRow(t.GetRow(id), store, t.InitPathTable(), 0)
}

func (t *Table) update(row *[]Data, data map[string]string) int {
	colUpdates := 0
	for _, column := range t.cols {
		val, ok := data[column.name]
		isPrimaryIndex := t.primaryIndex == column.name
		if isPrimaryIndex {
			console.YellowPrintln(fmt.Sprintf("Skip updating %s because it's primary key ", column.name))
			continue
		}
		if !ok {
			continue
		}
		if (*row)[column.number].Str() == val {
			continue
		}
		colUpdates++
		(*row)[column.number] = *FromString(val)
	}
	return colUpdates
}

func (t *Table) Update(id string, record map[string]string, store *Store) (map[string]any, error) {
	row := t.GetRow(id)
	if row == nil {
		return nil, errors.New("entity not found")
	}
	colsUpdated := t.update(&row, record)
	if colsUpdated > 0 {
		t.OnUpdate()
	}
	return t.MapJoinRow(t.GetRow(id), store, t.InitPathTable(), 0), nil
}

func (t *Table) DeleteById(id string) int {
	row := t.GetRow(id)

	if row == nil || len(t.rows) == 0 {
		return 0
	}

	newLength := len(t.rows) - 1
	newRows := make([][]Data, newLength)
	rawQuery := make(map[string][]string)
	rawQuery[t.primaryIndex] = []string{id}
	query := PrepareQuery(rawQuery, t)
	index := 0
	count := 0
	for _, r := range t.rows {
		if query.Matches(r) {
			count++
			continue
		}
		newRows[index] = r
		index++
	}

	t.rows = newRows
	t.OnUpdate()
	return count
}

func (t *Table) PlainGet(id string) map[string]any {
	return t.MapRow(t.GetRow(id))
}

func (t *Table) All() []map[string]any {
	res := make([]map[string]any, len(t.rows))
	for index, row := range t.rows {
		res[index] = t.MapRow(row)
	}
	return res
}

func (t *Table) Search(query Query) any {
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

func (t *Table) QueryRows(query Query) []*[]Data {
	rows := make([]*[]Data, 0)
	for _, row := range t.rows {
		if !query.Matches(row) {
			continue
		}
		rows = append(rows, &row)
	}
	return rows
}

func (t *Table) SearchRows(colName string, value string) [][]Data {
	rows := make([][]Data, 0)

	query := PrepareQuery(util.OneValueQuery(colName, value), t)
	for _, row := range t.rows {
		if !query.Matches(row) {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func (t *Table) Insert(data map[string]string) (string, error) {
	row := make([]Data, len(t.cols))
	var index string
	for _, column := range t.cols {
		val, ok := data[column.name]
		isPrimaryIndex := t.primaryIndex == column.name
		console.YellowPrintln(fmt.Sprintf("Setting column %s -> value %s ", column.name, val))
		if ok {
			row[column.number] = *FromString(val)
			if isPrimaryIndex {
				exist := t.GetRow(val)
				if exist != nil {
					return "", errors.New("duplicated primary key")
				}
				index = val
			}
			continue
		}
		if t.primaryIndex != column.name {
			continue
		}
		row[column.number] = *t.nextPrimaryIndex()
		index = row[column.number].Str()
	}
	t.rows = append(t.rows, row)
	t.indexes[t.primaryIndex][index] = row
	t.OnUpdate()
	return index, nil
}

func (t *Table) nextPrimaryIndex() *Data {
	primaryIndex := t.indexes[t.primaryIndex]
	if len(primaryIndex) == 0 {
		return FromString("0")
	}
	allNumbers := true
	var maxIndex int64
	maxIndex = 0
	for k := range primaryIndex {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			allNumbers = false
			break
		}
		if i > maxIndex {
			maxIndex = i
		}
	}
	if allNumbers {
		return FromString(strconv.FormatInt(maxIndex+1, 10))
	}
	return FromString(uuid.NewString())
}

func (t *Table) ToCsv() [][]string {
	w := len(t.cols)
	h := len(t.rows)

	res := make([][]string, h+1)
	header := make([]string, w)

	for i, col := range t.cols {
		header[i] = col.name
	}
	res[0] = header

	for j := range res {
		if j == 0 {
			continue
		}
		r := j - 1
		res[j] = make([]string, w)

		for i, data := range t.rows[r] {
			res[j][i] = data.Str()
		}

	}
	fmt.Println(res)
	return res
}

func TableToStructs[A any](t *Table, query Query, array *[]A) error {
	data := t.Search(query)
	h, err := json.Marshal(data)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(h))
	err = decoder.Decode(array)
	if err != nil {
		return err
	}
	return nil
}

func (t *Table) OnUpdate() {
	db.OnTableUpdate(t.File, t.ToCsv())
}

func (t *Table) BatchUpdate(batch Batch) {
	updatedCount := 0
	for _, row := range t.rows {
		err := batch.UpdateRow(&row)
		if err == nil {
			updatedCount++
		}
	}
	if updatedCount == 0 {
		console.YellowPrintln("Batch update 0 rows")
		return
	} else {
		console.GreenPrintln(fmt.Sprintf("Batch update %v rows", updatedCount))
	}
	t.OnUpdate()
}

func (t *Table) ColNumber(name string) int {
	for _, col := range t.cols {
		if col.name == name {
			return col.number
		}
	}
	return -1
}

func (t *Table) Schema() openapi.Schema {
	result := openapi.Schema{
		Type:       "object",
		Properties: make(map[string]openapi.Schema),
		Required:   []string{t.primaryIndex},
	}

	for _, col := range t.cols {
		result.Properties[col.name] = col.Schema()
	}

	return result
}

func InitTable(indexes Indexes, cols []*ColumnDefinition, nColumns int, nRows int, primaryIndex string) *Table {
	a := make([][]Data, nRows)
	for i := range a {
		a[i] = make([]Data, nColumns)
	}
	return &Table{
		indexes:        indexes,
		rows:           a,
		cols:           cols,
		primaryIndex:   primaryIndex,
		foreignIndexes: make([]*ColumnDefinition, 0)}
}
