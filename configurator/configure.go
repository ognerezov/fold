package configurator

import (
	"fold/console"
	"fold/csv"
	"fold/mem"
	"fold/openapi"
	"fold/path"
	"fold/router"
	"fold/security"
	"fold/util"
	goji "goji.io"
	"io/fs"
	"os"
)

func ConfigureServer(dataPath string, port int) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + dataPath)
	mux := goji.NewMux()
	mux.Use(router.LogRequest)

	store := *mem.TheStore
	clean := path.CreateRootCleaner(dataPath)
	apiDescription := openapi.InitApi(dataPath, port, "1")
	openapiRoute := openapi.Filename
	openapiFileName := dataPath + openapiRoute
	_ = os.Remove(openapiFileName)
	err := AppProviders.Export(dataPath)
	if err != nil {
		console.RedPrintln("export error: " + err.Error())
	}
	// TODO we should somehow register path parameter path latter
	// Probably by making array of tasks and then sorting them
	err = path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		route, filename, extension := path.Structure(dataPath, p, info, clean)
		switch extension {
		case ".csv":
			console.GreenPrintln("Registering table handlers for " + filename)
			records := csv.ReadCsvFile(filename)
			table := mem.TableFromRecords(records)
			table.File = filename
			store.SetTable(route, table)
			SetTableHandlers(route, mux, apiDescription)
		case ".json":
			console.GreenPrintln("Registering json handlers for " + filename)
			noSql, e := mem.LoadJson(filename)
			if e != nil {
				console.RedPrintln(e.Error())
			} else {
				store.SetNoSql(route, noSql)
				SetJsonHandlers(route, mux, apiDescription)
			}
		case ".fold":
			console.GreenPrintln("Registering fold action handler " + filename)
			SetControlHandlers(route, filename, mux, apiDescription)
		default:
			SetRawHandlers(route, filename, mux, apiDescription)
		}

		return nil
	})
	store.ReIndex()
	if err != nil {
		return nil, err
	}
	userTable, _ := store.GetTable(util.UserPath)
	if userTable != nil {
		security.EncodePasswords(userTable)
	}
	securityRulesTable, _ := store.GetTable("/security/rules")
	if securityRulesTable != nil {
		var rules []security.Rule
		e := mem.TableToStructs(securityRulesTable, mem.AllQuery(), &rules)
		if e != nil {
			console.RedPrintln(e.Error())
		}
		console.YellowPrintln("Applying security rules")
		if len(rules) > 0 {
			config := security.RulesSecurityConfig(rules)
			mux.Use(config.AuthorizeRequest)
		} else {
			mux.Use(security.Public.AuthorizeRequest)
		}
	}
	security.SetAuthHandlers(mux, App.Name())

	err = apiDescription.Save(openapiFileName)
	if err == nil {
		SetRawHandlers(openapi.Route, openapiFileName, mux, nil)
	} else {
		console.RedPrintln(err.Error())
	}
	return mux, nil
}
