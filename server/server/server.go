package server

import (
	"github.com/op/go-logging"
	"os"
)

// Server serves requests for content.
type Server interface {
	Start() error
	Stop() error
}

// FacadeServer contains all data and embeddings a Server would need.
type FacadeServer struct {
	Server
	Name    string
	Address string
	Port    int
}

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


