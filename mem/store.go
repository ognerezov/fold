package mem

import (
	"fmt"
	"fold/console"
	"fold/util"
)

type Store struct {
	kv     map[string]any
	tables map[string]*Table
}

func (s Store) SValue(key string, value any) {
	s.kv[key] = value
}

func (s Store) Value(key string) any {
	return s.kv[key]
}

func (s Store) Delete(key string) {
	delete(s.kv, key)
}

func (s Store) SetTable(key string, value *Table) {
	value.name = key
	s.tables[key] = value
}

func (s Store) GetTable(key string) *Table {
	return s.tables[key]
}

func (s Store) DeleteTable(key string) {
	delete(s.tables, key)
}

func (s Store) Get(table string, id string) map[string]any {
	return s.tables[table].Get(id, &s)
}

func (s Store) All(table string) []map[string]string {
	return s.tables[table].All()
}

func (s Store) ReIndex() {
	for name, table := range s.tables {
		for _, col := range table.cols {
			found, tableNames, colNames, err := util.NamingLookups(col.name)
			if !found || err != nil {
				continue
			}

			for nameIndex, tableName := range tableNames {
				tablePath := util.TableToPath(tableName)
				foreignTable, ok := s.tables[tablePath]
				if !ok {
					continue
				}
				foreignColName := colNames[nameIndex]
				for _, foreignCol := range foreignTable.cols {
					if foreignCol.name != foreignColName {
						continue
					}
					col.foreignTable = tablePath
					col.foreignColumn = foreignColName
					col.foreignUnique = foreignCol.isUnique
					table.foreignIndexes = append(table.foreignIndexes, col)
					console.CyanPrintln(fmt.Sprintf(
						"Created foreign idex %s -> %s on table: %s, column: %s ",
						tablePath, foreignColName,
						name, col.name))
				}
			}
		}
	}
	console.MagentaPrintln(fmt.Sprintf("%v", s.tables))
}

var (
	TheStore *Store = &Store{make(map[string]any), make(map[string]*Table)}
)
