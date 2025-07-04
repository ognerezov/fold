package configurator

import (
	"fmt"
	"fold/arguments"
	"fold/console"
	"fold/util"
	"google.golang.org/api/drive/v3"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type PortConfig struct {
	port      int
	path      string
	configure ServerConfigurator
	driveFile *drive.File
}

type PortsConfig []PortConfig

func SingleServer(dataPath string, configure ServerConfigurator) PortsConfig {
	return PortsConfig{
		PortConfig{
			port:      arguments.AppArguments.Port,
			path:      dataPath,
			configure: configure,
		}}
}

func ReadDir(dataPath string) ([]os.DirEntry, error) {
	f, err := os.OpenFile(filepath.FromSlash(dataPath), os.O_RDONLY, 0)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}
	defer util.CloseFile(f)

	files, err := f.ReadDir(0)
	return files, err
}

func ConfigurePorts(dataPath string, configure ServerConfigurator) PortsConfig {
	files, err := ReadDir(dataPath)

	if err != nil {
		return SingleServer(dataPath, configure)
	}

	res := make([]PortConfig, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		port, found := util.PathToInt(file.Name())

		if !found {
			continue
		}
		console.YellowPrintln("Checking root path " + file.Name())
		res = append(res, PortConfig{
			port:      port,
			path:      filepath.FromSlash(util.JoinedPath(dataPath, file)),
			configure: configure,
		})
	}

	if len(res) == 0 {
		return SingleServer(dataPath, configure)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].port < res[j].port
	})

	return res
}

func (p *PortConfig) Serve(address string, services map[int]*Service) error {
	mux, err := p.configure(p.path, p.port)
	if err != nil {
		console.RedPrintln("Can't start server")
		console.RedPrintln(err.Error())
		return err
	}
	addr := fmt.Sprintf("%s:%v", address, p.port)
	server := http.Server{Addr: addr, Handler: mux}
	services[p.port] = &Service{
		server:  &server,
		address: address,
		config:  p,
	}
	err = server.ListenAndServe()
	if err != nil {
		delete(services, p.port)
		return err
	}
	return nil
}
