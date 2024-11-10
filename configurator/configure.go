package configurator

import (
	"fmt"
	"fold/console"
	"fold/csv"
	"fold/mem"
	"fold/path"
	"fold/router"
	"fold/security"
	"fold/util"
	goji "goji.io"
	"io/fs"
	"path/filepath"
	"strings"
)

func ConfigureServer(dataPath string) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + dataPath)
	mux := goji.NewMux()
	mux.Use(router.LogRequest)

	store := *mem.TheStore
	clean := path.CreateRootCleaner(dataPath)
	path.Root = dataPath
	var err = path.ProcessPath(dataPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		var route = "/" + clean(filepath.Dir(path))
		var name = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

		if name != "index" {
			if strings.HasSuffix(route, "/") {
				route = fmt.Sprintf("%s%s", route, name)
			} else {
				route = fmt.Sprintf("%s/%s", route, name)
			}
		}
		route = util.TableToPath(route)
		var filename = fmt.Sprintf("%s/%s", dataPath, clean(path))
		switch filepath.Ext(path) {
		case ".csv":
			console.GreenPrintln("Registering table handlers for " + filename)
			records := csv.ReadCsvFile(filename)
			table := mem.TableFromRecords(records)
			table.File = filename
			store.SetTable(route, table)
			SetTableHandlers(route, mux)
		case ".json":
			console.GreenPrintln("Registering json handlers for " + filename)
			noSql, e := mem.LoadJson(filename)
			if e != nil {
				console.RedPrintln(e.Error())
			} else {
				store.SetNoSql(route, noSql)
				SetJsonHandlers(route, mux)
			}
		default:
			SetRawHandlers(route, filename, mux)
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
		console.YellowPrintln(fmt.Sprintf("%v", rules))
		if len(rules) > 0 {
			config := security.RulesSecurityConfig(rules)
			mux.Use(config.AuthorizeRequest)
		} else {
			mux.Use(security.Public.AuthorizeRequest)
		}
	}
	security.SetAuthHandlers(mux)
	return mux, nil
}
