package main

import (
	"fmt"
	"fold/configurator"
	"fold/console"
	"fold/csv"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting server")

	var argsWithProg = os.Args
	var argsWithoutProg = os.Args[1:]

	fmt.Println(argsWithProg)

	var dataPath = argsWithoutProg[0]

	var progLanguages = csv.ReadCsvFile(dataPath + "/languages.csv")

	for index, progLanguage := range progLanguages {
		fmt.Println(index, progLanguage)
	}
	//var table = *mem.TableFromRecords(progLanguages)
	//table.Print()
	//fmt.Println(table.GetRow("0"))
	//err := path.WalkPath(dataPath)
	//if err != nil {
	//	return
	//}
	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting server")
	console.GreenPrintln("___________________________")
	mux, err := configurator.Configure(dataPath)
	if err != nil {
		console.RedPrintln("Can't start server")
		console.RedPrintln(err.Error())
		return
	}
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
