package http

import (
	"github.com/clr1107/facade/server/server"
	"github.com/op/go-logging"
	"github.com/valyala/fasthttp"
	"strconv"
)

// ---------- HttpServer ----------

// HttpServer with FacadeServer embedded; using fasthttp.
type HttpServer struct {
	server.FacadeServer
	httpServer *fasthttp.Server
	logger *logging.Logger
}

func NewServer(name string, address string, port int) *HttpServer {
	s := &HttpServer{
		FacadeServer: server.FacadeServer{
			Name:    name,
			Address: address,
			Port:    port,
		},
		logger: server.SetupLogger(logging.DEBUG), // todo change
	}
	s.httpServer = &fasthttp.Server{
		Name:    s.Name,
		Handler: s.handle,
	}

	s.logger.Debugf("created new HttpServer `%s` to listen on `%s:%d`", name, address, port)
	return s
}

func (server *HttpServer) handle(ctx *fasthttp.RequestCtx) {
	// do smth
	server.logger.Infof("<=== %s", ctx.RequestURI())
}

func (server *HttpServer) Start() error {
	server.logger.Debugf("starting HttpServer `%s` on `%s:%d`", server.Name, server.Address, server.Port)
	return server.httpServer.ListenAndServe(server.Address + ":" + strconv.Itoa(server.Port))
}

func (server *HttpServer) Stop() error {
	server.logger.Debugf("stopping HttpServer `%s`", server.Name)
	return server.httpServer.Shutdown()
}
