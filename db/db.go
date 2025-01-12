package db

import (
	"fold/console"
	"google.golang.org/api/drive/v3"
)

type PendingTables map[string][][]string
type PendingUpdates map[string]any

type PendingDriveUpdates map[string]DriveUpdate
type DriveUpdate struct {
	File  *drive.File
	Value any
}

var (
	pendingDB           = make(PendingTables)
	pendingUpdates      = make(PendingUpdates)
	pendingDriveUpdates = make(PendingDriveUpdates)
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

func OnDriveUpdate(file *drive.File, val any) {
	pendingDriveUpdates[file.Id] = DriveUpdate{File: file, Value: val}
}

func ClearFileUpdate(name string) {
	delete(pendingUpdates, name)
}

func ClearDriveUpdate(name string) {
	delete(pendingDriveUpdates, name)
}

func Db() *PendingTables {
	return &pendingDB
}

func Pending() *PendingUpdates {
	return &pendingUpdates
}

func DrivePending() *PendingDriveUpdates {
	return &pendingDriveUpdates
}
