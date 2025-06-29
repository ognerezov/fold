package main

import (
	"flag"
	"fmt"
	"fold/arguments"
	"fold/configurator"
	"fold/console"
	"fold/initiator"
	"fold/migrations"
	"fold/recorder"
	"fold/threads"
	"fold/util"
	"strings"
)

func main() {
	var dataPath string
	flag.StringVar(&dataPath, "dir", "", "Working directory")
	flag.StringVar(&dataPath, "d", "", "Working directory (shorthand)")

	flag.BoolVar(&arguments.AppArguments.Cache, "cache", true, "Cache files requests")
	flag.BoolVar(&arguments.AppArguments.RegistrationAllowed, "reg", true, "User registration allowed")

	var filePath string
	flag.StringVar(&filePath, "file", "", "Serve or init with single file")
	flag.StringVar(&filePath, "f", "", "Serve or init with single file (shorthand)")

	flag.IntVar(&arguments.AppArguments.Port, "port", 3333, "Port to listen on")
	flag.IntVar(&arguments.AppArguments.Port, "p", 3333, "Port to listen on (shorthand)")

	flag.StringVar(&arguments.AppArguments.ApiPath, "api", "", "Server base path")
	flag.StringVar(&arguments.AppArguments.ApiPath, "a", "", "Server base path (shorthand)")

	flag.StringVar(&arguments.AppArguments.Name, "name", "fold", "Application name")
	flag.StringVar(&arguments.AppArguments.Name, "n", "fold", "Application name (shorthand)")

	var description string
	flag.StringVar(&description, "description", "", "Project description")
	var version string
	flag.StringVar(&version, "version", "1.0.0", "Project version")
	var allowOrigin string
	flag.StringVar(&allowOrigin, "origin", "", "Allow origin")
	flag.StringVar(&allowOrigin, "o", "", "Allow origin (shorthand)")

	var help bool
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	var init bool
	flag.BoolVar(&init, "init", false, "Init new project folder")
	flag.BoolVar(&init, "i", false, "Init new project folder (shorthand)")

	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "Migrate data")
	flag.BoolVar(&migrate, "m", false, "Migrate data (shorthand)")

	flag.StringVar(&arguments.InitArgs.Template, "template", initiator.DefaultTemplate, "Project template")
	flag.StringVar(&arguments.InitArgs.Template, "t", initiator.DefaultTemplate, "Project template (shorthand)")

	flag.StringVar(&arguments.InitArgs.SourceType, "source", migrations.Dir, "Data migration source type")
	flag.StringVar(&arguments.InitArgs.SourceType, "s", migrations.Dir, "Data migration source type (shorthand)")

	flag.StringVar(&arguments.InitArgs.DestinationType, "destination", migrations.Drive, "Data migration destination type")
	flag.StringVar(&arguments.InitArgs.DestinationType, "ds", migrations.Drive, "Data migration destination type (shorthand)")

	flag.StringVar(&arguments.InitArgs.CredentialsFile, "credentials", "", "Credentials file")
	flag.StringVar(&arguments.InitArgs.CredentialsFile, "c", "", "Credentials file (shorthand)")

	var drive string
	flag.StringVar(&drive, migrations.Drive, "", "Google Drive folder id")
	flag.StringVar(&drive, "dr", "", "Google Drive folder id (shorthand)")

	var recordFile string
	flag.StringVar(&recordFile, "record", "", "Record api json file")
	flag.StringVar(&recordFile, "r", "", "Record api json file")

	flag.Parse()
	fmt.Println(arguments.AppArguments)
	if help {
		flag.PrintDefaults()
		return
	}

	if dataPath != "" && filePath != "" && !init {
		panic("Both directory and single data file paths specified. Only one of them allowed.")
	}

	if init && migrate || (migrate && recordFile != "") {
		panic("Only one step at a time. Init or migrate. Not both at the same time.")
	}

	if dataPath == "" {
		dataPath = "./"
	}

	arguments.InitArgs.Output = dataPath
	arguments.InitArgs.Port = arguments.AppArguments.Port

	if init {
		if allowOrigin == "" {
			allowOrigin = fmt.Sprintf("http://localhost:%v", arguments.AppArguments.Port)
		}
		appConfig := configurator.AppConfig{
			Name:          arguments.AppArguments.Name,
			Description:   description,
			Version:       version,
			AllowOrigin:   allowOrigin,
			JWTSecret:     util.RandomString(12),
			GuestPassword: util.RandomString(8),
		}

		if drive != "" && arguments.InitArgs.CredentialsFile == "" {
			panic("Credentials file required for this init type")
		}

		if drive != "" {
			arguments.InitArgs.Destination = drive
			initiator.InitDrive(arguments.InitArgs, &appConfig)
			return
		}

		initiator.Init(arguments.InitArgs.Template, dataPath, arguments.InitArgs.Port, appConfig)
		if recordFile == "" {
			return
		}
	}

	if recordFile != "" {
		var recordApiDescription recorder.RecordApiDescription
		err := util.AnyFromJson(recordFile, &recordApiDescription)
		if err != nil {
			panic(err)
		}
		if recordApiDescription.Port == 0 {
			recordApiDescription.Port = arguments.AppArguments.Port
		}

		err = recorder.Record(dataPath, &recordApiDescription)
		if err != nil {
			panic(err)
		}
		return
	}

	if migrate {
		console.RedPrintln("Init data migration...")
		if arguments.InitArgs.SourceType == migrations.Dir {
			arguments.InitArgs.Source = dataPath
		} else if arguments.InitArgs.SourceType == migrations.Drive {
			arguments.InitArgs.Source = drive
		}
		if arguments.InitArgs.DestinationType == migrations.Drive {
			arguments.InitArgs.Destination = drive
		} else if arguments.InitArgs.DestinationType == migrations.Dir {
			arguments.InitArgs.Destination = dataPath
		}
		if arguments.InitArgs.Source == "" || arguments.InitArgs.Destination == "" {
			panic("Both arguments source and destination are required")
		}

		migrations.Migrate(arguments.InitArgs)
		return
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

	if drive != "" {
		if arguments.InitArgs.CredentialsFile == "" {
			panic("credentials file not specified")
		}
		configurator.CreateDriveApplication("", drive, arguments.InitArgs.CredentialsFile)
	}

	if filePath != "" {
		configurator.CreateFileApplication("", filePath)
		return
	}
	configurator.CreateApplication("", dataPath)
}
