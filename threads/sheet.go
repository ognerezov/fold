package threads

import (
	"fmt"
	"fold/db"
	"google.golang.org/api/sheets/v4"
)

type SheetsBatchWriter struct {
	requests []*sheets.Request
	id       string
}

func UpdateSheetAsync(update *db.SheetUpdate) {

	writer := SheetsBatchWriter{
		requests: update.Requests,
		id:       update.File.Id,
	}
	Async(writer, nil)
}

func (w SheetsBatchWriter) Call(any) (Message[string], ErrorMessage) {
	sheetService := providers.Google.Sheets
	process := fmt.Sprintf("Batch update of %v spread sheets ", w.id)
	if sheetService == nil {
		e := fmt.Errorf("sheets service not found in configurator")
		return EmptyMessage(process), CommonError(process, e)
	}

	response, err := sheetService.Spreadsheets.BatchUpdate(w.id, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: w.requests,
	}).Do()
	fmt.Println(response)
	// not retrying on Error response not to hit the quote
	db.ClearSheetsUpdate(w.id)
	if err != nil {
		return EmptyMessage(process), CommonError(process, err)
	}

	return CommonMessage(process, "success"), CommonError(process, nil)
}
