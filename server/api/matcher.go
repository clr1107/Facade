package api

import (
	"github.com/clr1107/facade/server/server"
	"net/url"
	"strings"
	"sync"
)

type Matcher interface {
	Match(request []byte) server.Pipe
}

// ---------- MultiMatcher ----------

type MultiMatcher struct {
	matchers []Matcher
	mux *sync.RWMutex
}

func NewMultiMatcher(matchers []Matcher) *MultiMatcher {
	return &MultiMatcher{
		matchers: matchers,
		mux: new(sync.RWMutex),
	}
}

func (multiMatcher *MultiMatcher) AppendMatcher(matcher Matcher) {
	multiMatcher.mux.Lock()
	multiMatcher.matchers = append(multiMatcher.matchers, matcher)
	multiMatcher.mux.Unlock()
}

func (multiMatcher *MultiMatcher) PrependMatcher(matcher Matcher) {
	multiMatcher.mux.Lock()
	multiMatcher.matchers = append([]Matcher{matcher}, multiMatcher.matchers...)
	multiMatcher.mux.Unlock()
}

func (multiMatcher *MultiMatcher) Match(request []byte) server.Pipe {
	multiMatcher.mux.RLock()
	defer multiMatcher.mux.RUnlock()

	for _, matcher := range multiMatcher.matchers {
		if x := matcher.Match(request); x != nil {
			return x
		}
	}

	return nil
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