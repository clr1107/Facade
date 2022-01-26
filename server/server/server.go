package server

import (
	"github.com/clr1107/facade/pkg/cache"
	"github.com/op/go-logging"
	"os"
)

func setupLogger(level logging.Level) *logging.Logger {
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
	GetFromCache(request []byte) []byte
}

// FacadeServer contains all data and embeddings a Server would need.
type FacadeServer struct {
	Server
	Name    string
	Address string
	Port    int
	Cache   cache.Cacher
	Matcher Matcher
	Errors  chan error
	Logger *logging.Logger
}

func NewFacadeServer(name string, address string, port int, matcher Matcher) FacadeServer {
	return FacadeServer{
		Name:    name,
		Address: address,
		Port:    port,
		Cache: cache.NewCache(nil, nil), // todo config
		Matcher: matcher,
		Errors: make(chan error, 1),
		Logger: setupLogger(logging.DEBUG), // todo config
	}
}

func (server *FacadeServer) GetFromCache(request []byte) []byte {
	if cached, ok := server.Cache.Get(string(request)); ok {
		if cast, ok := cached.([]byte); ok {
			return cast
		}
	}

	return nil
}


