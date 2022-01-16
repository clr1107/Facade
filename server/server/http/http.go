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
		FacadeServer: server.NewFacadeServer(name, address, port),
		logger: server.SetupLogger(logging.DEBUG), // todo config
	}
	s.httpServer = &fasthttp.Server{
		Name:    s.Name,
		Handler: s.handler,
		ErrorHandler: s.errorHandler,
	}

	s.logger.Debugf("created new HttpServer `%s` to listen on `%s:%d`", name, address, port)
	return s
}

func (server *HttpServer) handler(ctx *fasthttp.RequestCtx) {
	server.logger.Debugf("server <=== %s", ctx.RequestURI())

	if cached := server.GetCache(ctx.RequestURI()); cached != nil {
		server.logger.Debugf("cache  ===> %s", ctx.RequestURI())
		ctx.SetBody(cached)

		return
	}

	// todo ok go outbound now

	server.logger.Debugf("remote ===> %s", ctx.RequestURI())
	ctx.SetBody([]byte("This will be the outbound..."))
}

func (server *HttpServer) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	server.logger.Errorf("error dealing with request %s: %s", ctx.RequestURI(), err)
	ctx.Error("Error", 500)
}

func (server *HttpServer) Start() {
	server.logger.Infof("starting HttpServer `%s` on `%s:%d` ...", server.Name, server.Address, server.Port)

	go func() {
		server.Errors <- server.httpServer.ListenAndServe(server.Address + ":" + strconv.Itoa(server.Port))
	}()

	server.logger.Info("HttpServer started and listening")
}

func (server *HttpServer) Stop() error {
	server.logger.Info("stopping HttpServer `%s`", server.Name)
	close(server.Errors)
	return server.httpServer.Shutdown()
}
