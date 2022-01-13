package http

import (
	"github.com/clr1107/facade/server"
	"github.com/valyala/fasthttp"
)

// ---------- HttpServer ----------

// HttpServer with FacadeServer embedded; using fasthttp.
type HttpServer struct {
	server.FacadeServer
	httpServer *fasthttp.Server
}

func NewServer(name string, address string, port int) *HttpServer {
	s := &HttpServer{
		FacadeServer: server.FacadeServer{
			Name:    name,
			Address: address,
			Port:    port,
		},
	}
	s.httpServer = &fasthttp.Server{
		Name:    s.Name,
		Handler: s.handle,
	}

	return s
}

func (server *HttpServer) handle(ctx *fasthttp.RequestCtx) {

}
