package main

import (
	"fmt"
	"fold/configurator"
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

	var progLanguages = csv.ReadCsvFile(dataPath + "/programming_languages.csv")

	for index, progLanguage := range progLanguages {
		fmt.Println(index, progLanguage)
	}

	//err := path.WalkPath(dataPath)
	//if err != nil {
	//	return
	//}

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:8000", configurator.Configure(dataPath)))
}
