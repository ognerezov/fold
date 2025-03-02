package configurator

import (
	"context"
	"errors"
	"fmt"
	"fold/api"
	"fold/console"
	"fold/controls"
	"fold/mem"
	"fold/openapi"
	"fold/path"
	"fold/router"
	"fold/util"
	"net/http"
	"path/filepath"
)

var (
	App    *Application
	config AppConfig
)

const (
	projectRoute = "/src/project.js"
)

type Application struct {
	services map[int]*Service
	address  string
	ports    PortsConfig
	config   AppConfig
	fe       Endpoints
}

type AppConfig struct {
	Name          string   `json:"name" validate:"required"`
	Description   string   `json:"description"`
	Version       string   `json:"version"`
	AllowOrigin   string   `json:"allow_origin"`
	AuthProviders []string `json:"auth_providers"`
}

func loadConfig(dataPath string) AppConfig {
	var config AppConfig
	err := util.FromJson(dataPath+"/project.json", &config)
	if err != nil {
		panic(err)
	}
	return config
}

func CreateApplication(address string, dataPath string) {
	path.Root = dataPath
	config = loadConfig(dataPath)
	var err error
	config.AuthProviders, err = InitProviders(dataPath)
	_, ok := mem.TheStore.GetTable(util.UserPath)
	if ok {
		config.AuthProviders = append(config.AuthProviders, "password")
	}
	if err != nil {
		console.RedPrintln("InitProviders error: " + err.Error())
	}
	resources := ConfigureResources(dataPath)
	for _, resource := range resources {
		fmt.Println(resource)
		resource.Start()
	}

	ports := ConfigurePorts(dataPath, ConfigureServer)

	services := make(map[int]*Service, len(ports))
	l := len(ports)
	for i := 0; i < l-1; i++ {
		go initPort(address, ports[i], services)
	}
	// last port goes to console
	console.GreenPrintln("__________________________________")
	console.GreenPrintln("Server configured with default name")
	console.GreenPrintln("__________________________________")
	App = &Application{services: services, address: address, ports: ports, config: config}
	(*TheInstructions)[restart] = App.RestartControl
	console.MagentaPrintln(fmt.Sprintf("%v", *App))
	initPort(address, ports[l-1], services)
}

func CreateFileApplication(address string, dataPath string) {
	path.Root = dataPath
	ext := filepath.Ext(dataPath)
	if ext == ".drive" {
		panic("Serving drive file not supported. Json credentials have to be present in /providers/google folder ")
	}

	ports := SingleServer(dataPath, ConfigureSingleFileServer)

	services := make(map[int]*Service, len(ports))

	// last port goes to console
	console.GreenPrintln("__________________________________")
	console.GreenPrintln("Server configured with default name")
	console.GreenPrintln("__________________________________")
	App = &Application{services: services, address: address, ports: ports, config: AppConfig{Name: "fold"}}
	(*TheInstructions)[restart] = App.RestartControl
	console.MagentaPrintln(fmt.Sprintf("%v", *App))
	initPort(address, ports[0], services)
}

func initPort(address string, port PortConfig, services map[int]*Service) {
	err := port.Serve(address, services)
	if err != nil {
		console.RedPrintln(err.Error())
		console.MagentaPrintln(fmt.Sprintf("server at port %v was stopped", port))
	}
}

func (app *Application) RestartControl(string, any) *controls.Control {
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

func (app *Application) Do(data map[string]any, w http.ResponseWriter, _ *http.Request) {
	var err error
	var parameters api.RestartData
	err = util.Restructure(data, &parameters)
	if err != nil {
		router.ServerError(err, w)
		return
	}
	port := parameters.Port
	console.YellowPrintln(fmt.Sprintf("Restarting port %v", port))
	service, ok := app.services[port]
	err = errors.New("port not found")
	if !ok {
		router.ServerError(err, w)
		return
	}
	err = service.Restart(app.services)
	if err != nil {
		router.ServerError(err, w)
		return
	}
	router.WriteResponse(api.Ok(), w)
}

func (app *Application) Describe() ([]openapi.Parameter, map[string]openapi.Response) {
	return []openapi.Parameter{
		{
			Name: "port",
			Schema: openapi.Schema{
				Type: "integer",
			},
		},
	}, openapi.StatusResponse
}

func (app *Application) Name() string {
	return app.config.Name
}
