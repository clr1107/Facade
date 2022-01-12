package http

import (
	"github.com/clr1107/facade/server/requests"
	"github.com/valyala/fasthttp"
)

// ---------- HttpServer ----------

// HttpServer with FacadeServer embedded; using fasthttp.
type HttpServer struct {
	requests.FacadeServer
	httpServer *fasthttp.Server
}

func NewServer(name string, address string, port int) *HttpServer {
	server := &HttpServer{
		FacadeServer: requests.FacadeServer{
			Name:    name,
			Address: address,
			Port:    port,
		},
	}
	server.httpServer = &fasthttp.Server{
		Name:    server.Name,
		Handler: server.handle,
	}

	return server
}

func (server *HttpServer) handle(ctx *fasthttp.RequestCtx) {

}
