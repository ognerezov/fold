package main

import (
	"flag"
	"fmt"
	"fold/arguments"
	"fold/configurator"
	"fold/console"
	"fold/threads"
	"strings"
)

func main() {
	var dataPath string
	flag.StringVar(&dataPath, "dir", "", "Working directory")
	flag.StringVar(&dataPath, "d", "", "Working directory (shorthand)")

	flag.BoolVar(&arguments.AppArguments.Cache, "cache", true, "Cache files requests")
	flag.BoolVar(&arguments.AppArguments.RegistrationAllowed, "reg", true, "User registration allowed")

	var filePath string
	flag.StringVar(&filePath, "file", "", "Serve single file")
	flag.StringVar(&filePath, "f", "", "Working single (shorthand)")

	flag.IntVar(&arguments.AppArguments.Port, "port", 8000, "Port to listen on")
	flag.IntVar(&arguments.AppArguments.Port, "p", 8000, "Port to listen on (shorthand)")

	flag.StringVar(&arguments.AppArguments.ApiPath, "api", "", "Server base path")
	flag.StringVar(&arguments.AppArguments.ApiPath, "a", "", "Server base path (shorthand)")

	flag.StringVar(&arguments.AppArguments.Name, "name", "", "Application name")
	flag.StringVar(&arguments.AppArguments.Name, "n", "", "Application name (shorthand)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	flag.Parse()
	fmt.Println(arguments.AppArguments)
	if help {
		flag.PrintDefaults()
		return
	}

	if dataPath != "" && filePath != "" {
		panic("Both directory and single data file paths specified. Only one of them allowed.")
	}

	if dataPath == "" {
		dataPath = "./"
	}

	if arguments.AppArguments.ApiPath != "" && !strings.HasPrefix(arguments.AppArguments.ApiPath, "/") {
		arguments.AppArguments.ApiPath = "/" + arguments.AppArguments.ApiPath
	}

	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting Async service")
	console.GreenPrintln("___________________________")
	threads.Start(configurator.AppProviders)

	console.GreenPrintln("___________________________")
	console.GreenPrintln("Starting server")
	console.GreenPrintln("___________________________")
	if filePath != "" {
		configurator.CreateFileApplication("", filePath)
		return
	}
	configurator.CreateApplication("", dataPath)
}
