package configurator

import (
	"fmt"
	"fold/console"
	"fold/util"
	"net/http"
	"os"
	"sort"
)

const (
	DefaultPort = 8000
)

type PortConfig struct {
	port int
	path string
}

type PortsConfig []PortConfig

func SingleServer(dataPath string) PortsConfig {
	return PortsConfig{
		PortConfig{
			port: DefaultPort,
			path: dataPath,
		}}
}

func ReadDir(dataPath string) ([]os.DirEntry, error) {
	f, err := os.OpenFile(dataPath, os.O_RDONLY, 0)
	if err != nil {
		console.RedPrintln(err.Error())
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}(f)

	files, err := f.ReadDir(0)
	return files, err
}

func ConfigurePorts(dataPath string) PortsConfig {
	files, err := ReadDir(dataPath)

	if err != nil {
		return SingleServer(dataPath)
	}

	res := make([]PortConfig, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		port, found := util.PathToInt(file.Name())
		console.YellowPrintln("Checking root path " + file.Name())

		if !found {
			continue
		}
		res = append(res, PortConfig{port: port, path: util.JoinedPath(dataPath, file)})
	}

	if len(res) == 0 {
		return SingleServer(dataPath)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].port < res[j].port
	})

	return res
}

func (p *PortConfig) Serve(address string) (*Service, error) {
	mux, err := ConfigureServer(p.path, p.port)
	if err != nil {
		console.RedPrintln("Can't start server")
		console.RedPrintln(err.Error())
		return nil, err
	}
	addr := fmt.Sprintf("%s:%v", address, p.port)
	server := http.Server{Addr: addr, Handler: mux}
	err = server.ListenAndServe()
	if err != nil {
		return nil, err
	}

	return &Service{
		server: &server,
		config: p,
	}, nil
}
