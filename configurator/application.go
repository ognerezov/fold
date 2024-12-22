package configurator

import (
	"context"
	"errors"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/controls"
	"fold/path"
	"fold/util"
)

var (
	App *Application
)

type Application struct {
	services map[int]*Service
	address  string
	ports    PortsConfig
}

func CreateApplication(address string, dataPath string) {
	path.Root = dataPath
	err := InitProviders(dataPath)
	if err != nil {
		console.RedPrintln("InitProviders error: " + err.Error())
	}
	resources := ConfigureResources(dataPath)
	for _, resource := range resources {
		fmt.Println(resource)
		resource.Start()
	}

	ports := ConfigurePorts(dataPath)

	services := make(map[int]*Service, len(ports))
	l := len(ports)
	for i := 0; i < l-1; i++ {
		go initPort(address, ports[i], services)
	}
	// last port goes to console
	console.GreenPrintln("___________________________")
	console.GreenPrintln("Server configured.")
	console.GreenPrintln("___________________________")
	App = &Application{services: services, address: address, ports: ports}
	(*TheInstructions)[restart] = App.RestartControl()
	console.MagentaPrintln(fmt.Sprintf("%v", *App))
	initPort(address, ports[l-1], services)
}

func initPort(address string, port PortConfig, services map[int]*Service) {
	err := port.Serve(address, services)
	if err != nil {
		console.RedPrintln(err.Error())
		console.MagentaPrintln(fmt.Sprintf("server at port %v was stopped", port))
	}
}

func (app *Application) RestartControl() *controls.Control {
	var restartControl controls.Control
	restartControl = app
	return &restartControl
}

func (app *Application) Stop() {
	for _, service := range app.services {
		err := service.server.Shutdown(context.Background())
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}
}

func (app *Application) GetService(port int) (*Service, error) {
	service, ok := app.services[port]
	if ok {
		return service, nil
	}
	return nil, errors.New("port not found")
}

func (app *Application) SetService(port int, service *Service) {
	app.services[port] = service
}

func (app *Application) Do(data map[string]any) (any, error) {
	var err error
	var parameters api.RestartData
	err = util.Restructure(data, &parameters)
	if err != nil {
		return api.GetErrorResponse(err), err
	}
	port := parameters.Port
	console.YellowPrintln(fmt.Sprintf("Restarting port %v", port))
	fmt.Println(app)
	service, ok := app.services[port]
	err = errors.New("port not found")
	if !ok {
		return api.GetErrorResponse(err), err
	}
	err = service.Restart(app.services)
	if err != nil {
		return api.GetErrorResponse(err), err
	}
	return api.Ok(), nil
}
