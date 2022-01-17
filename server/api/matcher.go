package api

import (
	"github.com/clr1107/facade/server/server"
	"net/url"
	"strings"
)

type Matcher interface {
	Match(request []byte) server.Pipe
}

// ---------- RedirectHostMatcher ----------

type RedirectHostMatcher struct {
	Matcher
	Host string
}

func NewRedirectMatcher(host string) (*RedirectHostMatcher, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &RedirectHostMatcher{
		Host: strings.TrimSuffix(u.String(), "/"),
	}, nil
}

func (matcher RedirectHostMatcher) Match(request []byte) server.Pipe {
	return server.NewOutboundPipe(
		[]byte(matcher.Host + string(request)),
	)
}