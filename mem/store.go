package mem

type Store struct {
	kv     map[string]any
	tables map[string]Table
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

func (s Store) SetTable(key string, value Table) {
	s.tables[key] = value
}

func (s Store) GetTable(key string) Table {
	return s.tables[key]
}

func (s Store) DeleteTable(key string) {
	delete(s.tables, key)
}

func (s Store) Get(table string, id string) map[string]any {
	return s.tables[table].Get(id)
}

func (s Store) All(table string) []map[string]any {
	return s.tables[table].All()
}

var (
	TheStore *Store = &Store{make(map[string]any), make(map[string]Table)}
)
