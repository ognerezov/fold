package configurator

import (
	"context"
	"fmt"
	"fold/console"
	"net/http"
)

type Service struct {
	server  *http.Server
	config  *PortConfig
	address string
}

func (s *Service) Restart(services map[int]*Service) error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		console.RedPrintln(fmt.Sprintf("Failed to shutdown server on port %v", s.config.port))
		console.RedPrintln(err.Error())
		return err
	}
	console.YellowPrintln("Service stopped")
	go func() {
		err = s.config.Serve(s.address, services)
		if err != nil {
			console.RedPrintln(err.Error())
		}
	}()
	return nil
}
