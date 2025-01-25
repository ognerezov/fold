package configurator

import (
	"fmt"
	"fold/console"
	"fold/interfaces"
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

	next := interfaces.NewPhase()
	controlEndpoints := make(Endpoints)
	err = path.ProcessPath(dataPath, func(p string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		route, filename, extension := path.Structure(dataPath, p, info, clean)
		fileHandler := FilePath(filename)
		switch extension {
		case ".csv":
			console.GreenPrintln("Registering table handlers for " + filename)
			table, err := fileHandler.FetchCsv()
			if err != nil {
				console.RedPrintln(fmt.Sprintf("Error loading csc %s %v", filename, err))
			} else {
				store.SetTable(route, table, fileHandler)
				next.Append(SetTableHandlers(route, mux, apiDescription))
			}
		case ".json":
			console.GreenPrintln("Registering json handlers for " + filename)
			noSql, e := fileHandler.FetchNoSql()
			if e != nil {
				console.RedPrintln(e.Error())
			} else {
				store.SetNoSql(route, noSql, fileHandler)
				SetJsonHandlers(route, mux, apiDescription)
			}
		case ".fold":
			console.GreenPrintln("Registering fold action handler " + filename)
			SetControlHandlers(route, filename, mux, apiDescription, controlEndpoints)
		case ".drive":
			console.GreenPrintln("Registering google drive folder handler " + filename)
			routePath, id := path.SubStructure(p, info, clean)
			SetDriveHandlers(routePath, id, mux, apiDescription, next)
		default:
			console.RedPrintln("Registering raw file handler for " + filename)
			SetRawHandlers(route, filename, mux, apiDescription)
		}

		return nil
	})
	next.Act()
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
			pubRules := make([]security.Rule, 0)
			for _, rule := range rules {
				if rule.IsPublic() {
					pubRules = append(pubRules, rule)
				}
			}
			for _, endpoint := range controlEndpoints {
				for _, rule := range pubRules {
					if rule.AppliesTo(endpoint.Path, endpoint.Method) {
						endpoint.Public = true
					}
				}
			}
		} else {
			mux.Use(security.Public.AuthorizeRequest)
			for _, endpoint := range controlEndpoints {
				endpoint.Public = true
			}
		}
	}
	// even empty file is required
	controlFilename := dataPath + controlRoute
	err = util.SaveJavascript(controlFilename, controlEndpoints)
	if err != nil {
		console.RedPrintln(err.Error())
	} else {
		SetRawHandlers(controlRoute, controlFilename, mux, apiDescription)
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
