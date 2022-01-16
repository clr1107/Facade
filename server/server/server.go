package server

import (
	"github.com/clr1107/facade/pkg/cache"
	"github.com/op/go-logging"
	"os"
)

func SetupLogger(level logging.Level) *logging.Logger {
	log := logging.MustGetLogger("facade")
	format := logging.MustStringFormatter(
		`%{color}[%{time:15:04:05}] [%{level}] (%{longfunc}): %{color:reset}%{message}`,
	)

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	levelled := logging.SetBackend(backendFormatter)
	levelled.SetLevel(level, "")

	return log
}

// Server serves requests for content.
type Server interface {
	Start()
	Stop() error
	GetCache(request []byte) []byte
}

// FacadeServer contains all data and embeddings a Server would need.
type FacadeServer struct {
	Server
	Name    string
	Address string
	Port    int
	Cache cache.Cacher
	Errors chan error
}

func NewFacadeServer(name string, address string, port int) FacadeServer {
	return FacadeServer{
		Name:    name,
		Address: address,
		Port:    port,
		Cache: cache.NewCache(nil, nil), // todo config
		Errors: make(chan error),
	}
}

func (server *FacadeServer) GetCache(request []byte) []byte {
	if cached, ok := server.Cache.Get(string(request)); ok {
		if cast, ok := cached.([]byte); ok {
			return cast
		}
	}

	return nil
}


