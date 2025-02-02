package configurator

import (
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/interfaces"
	"fold/mem"
	"fold/openapi"
	"fold/path"
	"fold/router"
	"fold/security"
	"fold/util"
	goji "goji.io"
	"os"
	"path/filepath"
)

func ConfigureSingleFileServer(dataPath string, port int) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + dataPath)
	mux := goji.NewMux()
	mux.Use(router.LogRequest)

	store := *mem.TheStore
	dir := filepath.Dir(dataPath)
	clean := path.CreateRootCleaner(dir)
	apiDescription := openapi.InitApi(dataPath, port, "1")
	openapiRoute := openapi.Filename
	openapiFileName := "./" + openapiRoute
	_ = os.Remove(openapiFileName)
	info, err := os.Lstat(dataPath)
	if err != nil {
		panic(err)
	}
	next := interfaces.NewPhase()
	controlEndpoints := make(Endpoints)
	ConfigureFile(info.Name(), info, dir, clean, next, mux, apiDescription, controlEndpoints)
	next.Act()
	store.ReIndex()

	userTable, _ := store.GetTable(util.UserPath)
	if userTable != nil {
		security.EncodePasswords(userTable)
	}
	mux.Use(security.Public.AuthorizeRequest)
	for _, endpoint := range controlEndpoints {
		endpoint.Public = true
	}
	fmt.Println(arguments.AppArguments)
	security.SetAuthHandlers(arguments.AppArguments.ApiPath, mux, App.Name())

	err = apiDescription.Save(openapiFileName)
	if err == nil {
		SetRawHandlers(arguments.AppArguments.ApiPath+openapi.Route, openapiFileName, mux, nil)
	} else {
		console.RedPrintln(err.Error())
	}
	return mux, nil
}
