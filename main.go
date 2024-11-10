package main

import (
	"fmt"
	"fold/configurator"
	"fold/console"
	"fold/threads"
	"os"
)

func main() {
	var argsWithProg = os.Args
	var argsWithoutProg = os.Args[1:]

	fmt.Println(argsWithProg)

	var dataPath = argsWithoutProg[0]
	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting Async service")
	console.GreenPrintln("___________________________")
	threads.Start()

	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting server")
	console.GreenPrintln("___________________________")
	configurator.App = configurator.CreateApplication("", dataPath)
}
