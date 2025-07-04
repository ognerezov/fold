package configurator

import (
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/gcloud"
	"fold/interfaces"
	"fold/mem"
	"fold/migrations"
	"fold/openapi"
	"fold/path"
	"fold/router"
	"fold/security"
	"fold/util"
	goji "goji.io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	ProcessedPersonalRoutes = map[string]bool{}
	doNotServe              = []string{"openapi.json", "fold.js", "project.js", "google.js"}
)

const (
	RawRoutesFolder = "_raw_routes"
	RawSeparator    = "_"
)

type ServerConfigurator func(string, int) (*goji.Mux, error)

func ConfigureServer(dataPath string, port int) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + dataPath)
	mux, store, apiDescription := initialize(dataPath, port)
	clean := path.CreateRootCleaner(dataPath)
	openapiRoute := openapi.Filename
	openapiFileName := dataPath + openapiRoute
	_ = os.Remove(filepath.FromSlash(openapiFileName))
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
		ConfigureFile(p, info, dataPath, clean, next, mux, apiDescription, controlEndpoints)

		return nil
	})
	next.Act()
	store.ReIndex()
	if err != nil {
		return nil, err
	}

	protected := setupSecurity(store, mux, controlEndpoints)
	ok, _ := util.Exists(dataPath + srtRoute)
	if ok {
		// files required to run ui
		err = util.SaveJavascript(dataPath+projectFile, config)
		if err != nil {
			panic(err)
		} else {
			SetRawHandlers(projectRoute, dataPath+projectFile, mux, apiDescription)
		}

		// even empty file is required
		controlFilename := dataPath + controlFile
		err = util.SaveJavascript(controlFilename, controlEndpoints)
		if err != nil {
			console.RedPrintln(err.Error())
		} else {
			SetRawHandlers(arguments.AppArguments.ApiPath+controlRoute, controlFilename, mux, apiDescription)
		}
	}
	if protected {
		security.SetAuthHandlers(arguments.AppArguments.ApiPath, mux, App.Name())
	}
	err = apiDescription.Save(openapiFileName)
	if err == nil {
		SetRawHandlers(arguments.AppArguments.ApiPath+openapi.Route, openapiFileName, mux, nil)
	} else {
		console.RedPrintln(err.Error())
	}
	return mux, nil
}

func setupSecurity(store *mem.Store, mux *goji.Mux, controlEndpoints Endpoints) bool {
	userTable, ok := store.GetUserTable()
	if userTable != nil {
		security.EncodePasswords(userTable)
	}
	securityRulesTable, _ := store.GetTable(arguments.AppArguments.ApiPath + "/security/rules")
	if ok {
		config.AuthProviders = append(config.AuthProviders, "password")
	}
	if securityRulesTable != nil {
		var rules []security.Rule
		e := mem.TableToStructs(securityRulesTable, mem.AllQuery(), &rules)
		if e != nil {
			console.RedPrintln(e.Error())
		}
		console.YellowPrintln("Applying security rules")
		if len(rules) == 0 {
			for _, endpoint := range controlEndpoints {
				endpoint.Public = true
			}
			return false
		}
		conf := security.RulesSecurityConfig(rules)
		mux.Use(conf.AuthorizeRequest)
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
		return true
	}
	for _, endpoint := range controlEndpoints {
		endpoint.Public = true
	}
	return false

}

func initialize(dataPath string, port int) (*goji.Mux, *mem.Store, *openapi.ApiDescription) {
	ProcessedPersonalRoutes = map[string]bool{}

	(*TheInstructions)[createProject] = gcloud.ConfigureProjectCreator(dataPath)

	mux := goji.NewMux()
	mux.Use(router.LogRequest)
	mux.Use(AddHeaders)

	store := mem.TheStore
	apiDescription := openapi.InitApi(config.Name, port, "1")
	return mux, store, apiDescription
}

func ConfigureFile(p string, info fs.FileInfo, dataPath string, clean path.DirMapper, next *interfaces.Phase, mux *goji.Mux, apiDescription *openapi.ApiDescription, controlEndpoints Endpoints) {
	if strings.Contains(p, RawRoutesFolder) {
		ConfigureRawFolder(dataPath, p, mux, apiDescription)
		return
	}

	store := mem.TheStore
	route, filename, extension := path.Structure(p, info, clean)
	route = arguments.AppArguments.ApiPath + route
	fileHandler := mem.FilePath(filename)
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
		migrationHandlers[id] = migrations.CreateDriveHandler(AppProviders.Google, id)
		SetDriveHandlers(arguments.AppArguments.ApiPath+routePath, id, mux, apiDescription, next, migrationHandlers[id], controlEndpoints)
	default:
		console.RedPrintln("Registering raw file handler for " + filename)
		SetRawHandlers(route, filename, mux, apiDescription)
	}
}

func ConfigureRawFolder(dataPath, _filename string, mux *goji.Mux, api *openapi.ApiDescription) {
	clean := path.CreateRootCleaner(fmt.Sprintf("%v/%v", dataPath, RawRoutesFolder))
	_name := clean(_filename)
	parts := strings.Split(_name, RawSeparator)
	pathParts := strings.Split(parts[0], "/")
	method := strings.ToUpper(pathParts[len(pathParts)-1])
	name := strings.Join(parts[1:], RawSeparator)
	route := filepath.ToSlash("/" + strings.Join(pathParts[:len(pathParts)-1], "/") + "/" + strings.TrimSuffix(name, filepath.Ext(name)))
	fmt.Println(_filename)
	SetRawMethodHandlers(route, _filename, method, mux, api)
}

func ParsePersonalRoute(path *string) *string {
	if util.PersonalDataRegexp.MatchString(*path) {
		match := util.PersonalDataRegexp.FindStringSubmatch(*path)

		if match != nil && len(match) > 1 {
			userId := match[1]
			value := strings.Replace(*path, userId, ":user", -1)

			_, ok := ProcessedPersonalRoutes[value]

			if ok {
				return nil
			}
			ProcessedPersonalRoutes[value] = true
			return &value
		}
		return path
	}

	return path
}
