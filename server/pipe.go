package server

import (
	"github.com/clr1107/facade/pkg/loadbalance"
	"github.com/valyala/fasthttp"
)

type Pipe interface {
	Pull() ([]byte, error)
}

// ---------- OutboundPipe ----------

type OutboundPipe struct {
	Pipe
	Address []byte
	Client *fasthttp.Client
}

func NewOutboundPipe(address []byte) *OutboundPipe {
	return &OutboundPipe{
		Address: address,
		Client: nil,
	}
}

func (pipe OutboundPipe) Pull() ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURIBytes(pipe.Address)
	if pipe.Client == nil {
		pipe.Client = new(fasthttp.Client)
	}

	if err := pipe.Client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

// ---------- ProxiedPipe ----------

type BalancedOutboundPipe struct {
	OutboundPipe
	Balancer loadbalance.Balancer
}

func NewBalancedOutboundPipe(pipe *OutboundPipe, balancer loadbalance.Balancer) *BalancedOutboundPipe {
	return &BalancedOutboundPipe{
		*pipe, balancer,
	}
}

func (pipe BalancedOutboundPipe) Pull() ([]byte, error) {
	medium := pipe.Balancer.Get()
	medium.Apply(pipe.Client)

	return pipe.OutboundPipe.Pull()
}
