package configurator

import (
	"context"
	"log"
	"net/http"
)

var (
	TheServer *http.Server
)

func Start() {
	log.Fatal(TheServer.ListenAndServe())
}

func Stop() {
	log.Fatal(TheServer.Shutdown(context.Background()))
}
