package api

import (
	srv "github.com/clr1107/facade/server/server"
	"net/url"
	"strings"
	"sync"
)

// ---------- MultiMatcher ----------

type MultiMatcher struct {
	matchers []srv.Matcher
	mux *sync.RWMutex
}

func (multiMatcher *MultiMatcher) AppendMatcher(matcher srv.Matcher) {
	multiMatcher.mux.Lock()
	multiMatcher.matchers = append(multiMatcher.matchers, matcher)
	multiMatcher.mux.Unlock()
}

func (multiMatcher *MultiMatcher) PrependMatcher(matcher srv.Matcher) {
	multiMatcher.mux.Lock()
	multiMatcher.matchers = append([]srv.Matcher{matcher}, multiMatcher.matchers...)
	multiMatcher.mux.Unlock()
}

func (multiMatcher *MultiMatcher) Match(request []byte) srv.Pipe {
	multiMatcher.mux.RLock()
	defer multiMatcher.mux.RUnlock()

	for _, matcher := range multiMatcher.matchers {
		if x := matcher.Match(request); x != nil {
			return x
		}
	}

	return nil
}

// ---------- PredicatedMatcher ----------

type MatcherPredicate func(request []byte) bool

type PredicatedMatcher struct {
	predicate MatcherPredicate
	matcher srv.Matcher
}

func (pMatcher *PredicatedMatcher) Match(request []byte) srv.Pipe {
	if pMatcher.predicate(request) {
		return pMatcher.matcher.Match(request)
	}

	return nil
}

// ---------- API functions ----------

func NewRedirectMatcher(host string) (*srv.RedirectHostMatcher, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &srv.RedirectHostMatcher{
		Host: strings.TrimSuffix(u.String(), "/"),
	}, nil
}

func NewMultiMatcher(matchers ...srv.Matcher) *MultiMatcher {
	return &MultiMatcher{
		matchers: matchers,
		mux: new(sync.RWMutex),
	}
}

func NewPredicateMatcher(predicate MatcherPredicate, matcher srv.Matcher) *PredicatedMatcher {
	return &PredicatedMatcher{
		predicate: predicate,
		matcher: matcher,
	}
}