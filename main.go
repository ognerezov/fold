package main

import (
	"fmt"
	"fold/csv"
	"fold/path"
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

	err := path.WalkPath(dataPath)
	if err != nil {
		return
	}
}
