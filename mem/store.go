package mem

import (
	"errors"
	"fmt"
	"fold/console"
	"fold/util"
)

type Store struct {
	noSql         map[string]*NoSql
	tables        map[string]*Table
	cache         map[string][]byte
	cacheFetchers map[string]BytesFetcher
	noSqlFetchers map[string]NoSqlFetcher
	csvFetchers   map[string]CsvFetcher
}

func (s Store) Cache(key string, b []byte, refresh BytesFetcher) {
	s.cache[key] = b
	s.cacheFetchers[key] = refresh
}

// TODO Refresh action
func (s Store) RefreshCache(key string) error {
	fetcher, ok := s.cacheFetchers[key]
	if !ok {
		return errors.New("key not found")
	}
	bytes, err := fetcher.Fetch()
	if err != nil {
		return err
	}
	s.cache[key] = bytes
	return nil
}

func (s Store) GetCached(key string) ([]byte, bool) {
	res, ok := s.cache[key]
	return res, ok
}

func (s Store) SetNoSql(key string, value *NoSql, fetcher NoSqlFetcher) {
	s.noSql[key] = value
	s.noSqlFetchers[key] = fetcher
}

func (s Store) NoSql(key string) (*NoSql, bool) {
	n, ok := s.noSql[key]
	return n, ok
}

func (s Store) RefreshNoSql(key string) error {
	fetcher, ok := s.noSqlFetchers[key]
	if !ok {
		return errors.New("key not found")
	}
	res, err := fetcher.FetchNoSql()
	if err != nil {
		return err
	}
	s.noSql[key] = res
	return nil
}

func (s Store) SetTable(key string, value *Table, fetcher CsvFetcher) {
	value.name = key
	s.tables[key] = value
	s.csvFetchers[key] = fetcher
}

func (s Store) RefreshTable(key string) error {
	fetcher, ok := s.csvFetchers[key]
	if !ok {
		return errors.New("key not found")
	}
	res, err := fetcher.FetchCsv()
	if err != nil {
		return err
	}
	s.tables[key] = res
	return nil
}

func (s Store) Refresh(key string) error {
	_, ok := s.cacheFetchers[key]
	if ok {
		return s.RefreshCache(key)
	}
	_, ok = s.noSqlFetchers[key]
	if ok {
		return s.RefreshNoSql(key)
	}
	_, ok = s.csvFetchers[key]
	if ok {
		return s.RefreshTable(key)
	}
	return errors.New("no handler found")
}

func (s Store) GetTable(key string) (*Table, bool) {
	console.CyanPrintln("Get table: " + key)
	t, ok := s.tables[key]
	return t, ok
}

func (s Store) DeleteTable(key string) {
	delete(s.tables, key)
}

func (s Store) Get(table string, id string) map[string]any {
	return s.tables[table].Get(id, &s)
}

func (s Store) PlainGet(table string, id string) map[string]any {
	return s.tables[table].PlainGet(id)
}

func (s Store) All(table string) []map[string]any {
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
					col.foreignTable = tableName
					col.foreignColumn = foreignColName
					col.foreignUnique = foreignCol.isUnique
					table.foreignIndexes = append(table.foreignIndexes, col)
					console.CyanPrintln(fmt.Sprintf(
						"Created foreign idex %s -> %s on table: %s, column: %s ",
						tableName, foreignColName,
						name, col.name))

					formattedName := util.PathToTable(name)
					foreignTable.foreignIndexes = append(foreignTable.foreignIndexes,
						col.BackReference(formattedName, foreignCol.number))

					console.GreenPrintln(fmt.Sprintf(
						"Created foreign back reference %s -> %s  table: %s",
						foreignTable.name,
						col.name, formattedName))
				}
			}
		}
	}
}

var (
	TheStore *Store = &Store{
		noSql:         make(map[string]*NoSql),
		tables:        make(map[string]*Table),
		cache:         make(map[string][]byte),
		cacheFetchers: make(map[string]BytesFetcher),
		noSqlFetchers: make(map[string]NoSqlFetcher),
		csvFetchers:   make(map[string]CsvFetcher),
	}
)
