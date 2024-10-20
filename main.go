package main

import (
	"fmt"
	"fold/configurator"
	"fold/console"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting server")

	var argsWithProg = os.Args
	var argsWithoutProg = os.Args[1:]

	fmt.Println(argsWithProg)

	var dataPath = argsWithoutProg[0]

	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting server")
	console.GreenPrintln("___________________________")
	mux, err := configurator.Configure(dataPath)
	if err != nil {
		console.RedPrintln("Can't start server")
		console.RedPrintln(err.Error())
		return
	}
	addr := "localhost:8000"
	configurator.TheServer = &http.Server{Addr: addr, Handler: mux}
	configurator.Start()
}
