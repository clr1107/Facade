package http

import (
	"fmt"
	"github.com/clr1107/facade/pkg/cache"
	srv "github.com/clr1107/facade/server"
	"github.com/clr1107/facade/server/server"
	"github.com/valyala/fasthttp"
	"strconv"
)

// ---------- HttpServer ----------

// HttpServer with FacadeServer embedded; using fasthttp.
type HttpServer struct {
	server.FacadeServer
	httpServer *fasthttp.Server
}

func NewServer(name string, address string, port int) *HttpServer {
	s := &HttpServer{
		FacadeServer: server.NewFacadeServer(name, address, port),
	}
	s.httpServer = &fasthttp.Server{
		Name:    s.Name,
		Handler: s.handler,
		ErrorHandler: s.errorHandler,
	}

	s.Logger.Debugf("created new HttpServer `%s` to listen on `%s:%d`", name, address, port)
	return s
}

func (server *HttpServer) handler(ctx *fasthttp.RequestCtx) {
	server.Logger.Debugf("server <=== %s", ctx.RequestURI())

	if cached := server.GetFromCache(ctx.RequestURI()); cached != nil {
		server.Logger.Debugf("cache  ===> %s", ctx.RequestURI())
		ctx.SetBody(cached)

		return
	}

	server.Logger.Debugf("remote ===> %s", ctx.RequestURI())
	pipe := srv.NewOutboundPipe([]byte("https://api.ipify.org?format=json")) // todo obviously for testing...
	// todo load balancing

	if resp, err := pipe.Pull(); err != nil {
		ctx.Error("an unknown error has occurred", 500)
		server.Errors <- err
	} else {
		if err := server.Cache.Put(string(ctx.RequestURI()), cache.NewCacheUnit(resp, nil)); err != nil {
			server.Errors <- err
		}

		ctx.SetBody(resp)
	}
}

func (server *HttpServer) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	server.Errors <- fmt.Errorf("error dealing with request %s: %s", ctx.RequestURI(), err)
	ctx.Error("Error", 500)
}

func (server *HttpServer) Start() {
	server.Logger.Infof("starting HttpServer `%s` on `%s:%d` ...", server.Name, server.Address, server.Port)

	go func() {
		err := server.httpServer.ListenAndServe(server.Address + ":" + strconv.Itoa(server.Port))
		if err != nil {
			server.Logger.Errorf("ListenAndServe: %s", err)
		}
	}()

	server.Logger.Info("HttpServer started and listening")
}

func (server *HttpServer) Stop() error {
	server.Logger.Infof("stopping HttpServer `%s` ...", server.Name)

	close(server.Errors)
	err := server.httpServer.Shutdown()

	server.Logger.Info("stopped HttpServer `%s`", server.Name)
	return err
}
