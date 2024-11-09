package configurator

import (
	"context"
	"fmt"
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
	fmt.Println(len(ports))
	l := len(ports)
	for i := 0; i < l-1; i++ {
		go initPort(address, ports, i, services)
	}
	// last port goes to console
	initPort(address, ports, l-1, services)

	return &Application{services: services}
}

func initPort(address string, ports PortsConfig, index int, services []*Service) {
	port := ports[index]
	server, err := port.Serve(address)
	if err != nil {
		console.RedPrintln(err.Error())
		panic(err)
	}
	services[index] = server
}

func (app *Application) Stop() {
	for _, service := range app.services {
		err := service.server.Shutdown(context.Background())
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
}
