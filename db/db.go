package db

import "fold/console"

type PendingTables map[string][][]string
type PendingUpdates map[string]any

var (
	pendingDB      = make(PendingTables)
	pendingUpdates = make(PendingUpdates)
)

func OnTableUpdate(name string, table [][]string) {
	pendingDB[name] = table
}

func ClearTableUpdate(name string) {
	console.YellowPrintln("Clearing table updates for " + name)
	delete(pendingDB, name)
}

func OnFileUpdate(name string, val any) {
	pendingUpdates[name] = val
}

func ClearFileUpdate(name string) {
	delete(pendingUpdates, name)
}

func Db() *PendingTables {
	return &pendingDB
}

func Pending() *PendingUpdates {
	return &pendingUpdates
}
