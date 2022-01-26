package server

type Matcher interface {
	Match(request []byte) Pipe
}

// ---------- RedirectHostMatcher ----------

type RedirectHostMatcher struct {
	Matcher
	Host string
}

func (matcher RedirectHostMatcher) Match(request []byte) Pipe {
	return NewOutboundPipe(
		[]byte(matcher.Host + string(request)),
	)
}