package main

import (
	"flag"
	"fold/configurator"
	"fold/console"
	"fold/threads"
)

func main() {
	var dataPath string
	flag.StringVar(&dataPath, "dir", "./", "Working directory")
	flag.StringVar(&dataPath, "d", "./", "Working directory (shorthand)")

	flag.BoolVar(&configurator.AppArguments.Cache, "cache", true, "Cache files requests")
	flag.BoolVar(&configurator.AppArguments.RegistrationAllowed, "reg", true, "User registration allowed")

	flag.Parse()
	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting Async service")
	console.GreenPrintln("___________________________")
	threads.Start()

	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting server")
	console.GreenPrintln("___________________________")
	configurator.CreateApplication("", dataPath)
}
