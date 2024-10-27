package db

type PendingTables map[string][][]string

var (
	pendingDB = make(PendingTables)
)

func OnTableUpdate(name string, table [][]string) {
	pendingDB[name] = table
}

func Db() *PendingTables {
	return &pendingDB
}
