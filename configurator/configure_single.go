package configurator

import (
	"fold/console"
	"fold/interfaces"
	"fold/mem"
	"fold/openapi"
	"fold/path"
	"fold/router"
	"fold/security"
	goji "goji.io"
	"os"
	"path/filepath"
)

func ConfigureSingleFileServer(dataPath string, port int) (*goji.Mux, error) {
	console.YellowPrintln("Configure server for dir " + dataPath)
	mux := goji.NewMux()
	mux.Use(router.LogRequest)

	store := mem.TheStore
	dir := filepath.Dir(dataPath)
	clean := path.CreateRootCleaner(dir)
	apiDescription := openapi.InitApi(dataPath, port, "1")
	info, err := os.Lstat(dataPath)
	if err != nil {
		panic(err)
	}
	next := interfaces.NewPhase()
	controlEndpoints := make(Endpoints)
	ConfigureFile(info.Name(), info, dir, clean, next, mux, apiDescription, controlEndpoints)
	next.Act()
	store.ReIndex()

	mux.Use(security.Public.AuthorizeRequest)
	for _, endpoint := range controlEndpoints {
		endpoint.Public = true
	}
	return mux, nil
}
