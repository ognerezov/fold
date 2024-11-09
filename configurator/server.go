package configurator

import (
	"net/http"
)

type Service struct {
	server *http.Server
	config *PortConfig
}
