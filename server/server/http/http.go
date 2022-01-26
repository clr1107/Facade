package http

import (
	"fmt"
	"github.com/clr1107/facade/pkg/cache"
	"github.com/clr1107/facade/server/server"
	"github.com/valyala/fasthttp"
	"strconv"
)

// ---------- HttpServer ----------

// HttpServer ; using fasthttp.
type HttpServer struct {
	server.FacadeServer
	httpServer *fasthttp.Server
}

func NewServer(name string, address string, port int, matcher server.Matcher) *HttpServer {
	s := &HttpServer{
		FacadeServer: server.NewFacadeServer(name, address, port, matcher),
	}
	s.httpServer = &fasthttp.Server{
		Name:    s.Name,
		Handler: s.handler,
		ErrorHandler: s.errorHandler,
	}

	s.Logger.Debugf("created new HttpServer `%s` to listen on `%s:%d`", name, address, port)
	return s
}

func (httpServer *HttpServer) handler(ctx *fasthttp.RequestCtx) {
	httpServer.Logger.Debugf("httpServer <=== %s", ctx.RequestURI())

	if cached := httpServer.GetFromCache(ctx.RequestURI()); cached != nil {
		httpServer.Logger.Debugf("cache      ===> %s", ctx.RequestURI())
		ctx.SetBody(cached)

		return
	}

	httpServer.Logger.Debugf("remote     ===> %s", ctx.RequestURI())

	pipe := httpServer.Matcher.Match(ctx.RequestURI())
	if pipe == nil {
		ctx.Error("could not match url given", 400)
		return
	}

	if resp, err := pipe.Pull(); err != nil {
		ctx.Error("an unknown error has occurred whilst pulling pipe", 500)
		httpServer.Errors <- err
	} else {
		if err := httpServer.Cache.Put(string(ctx.RequestURI()), cache.NewCacheUnit(resp, nil)); err != nil {
			httpServer.Errors <- err
		}

		ctx.SetBody(resp)
	}
}

func (httpServer *HttpServer) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	httpServer.Errors <- fmt.Errorf("error dealing with request %s: %s", ctx.RequestURI(), err)
	ctx.Error("internal server error", 500)
}

func (httpServer *HttpServer) Start() {
	httpServer.Logger.Infof("starting HttpServer `%s` on `%s:%d` ...", httpServer.Name, httpServer.Address, httpServer.Port)

	go func() {
		err := httpServer.httpServer.ListenAndServe(httpServer.Address + ":" + strconv.Itoa(httpServer.Port))
		if err != nil {
			httpServer.Logger.Errorf("ListenAndServe: %s", err)
		}
	}()

	httpServer.Logger.Info("HttpServer started and listening")
}

func (httpServer *HttpServer) Stop() error {
	httpServer.Logger.Infof("stopping HttpServer `%s` ...", httpServer.Name)

	close(httpServer.Errors)
	err := httpServer.httpServer.Shutdown()

	httpServer.Logger.Info("stopped HttpServer `%s`", httpServer.Name)
	return err
}
