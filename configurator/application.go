package configurator

import (
	"context"
	"fold/console"
)

var (
	App *Application
)

type Application struct {
	services []*Service
}

func CreateApplication(address string, dataPath string) *Application {
	ports := ConfigurePorts(dataPath)
	services := make([]*Service, len(ports))
	for i, port := range ports {
		server, err := port.Serve(address)
		if err != nil {
			panic(err)
		}
		services[i] = server
	}

	return &Application{services: services}
}

func (app *Application) Stop() {
	for _, service := range app.services {
		err := service.server.Shutdown(context.Background())
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
}
