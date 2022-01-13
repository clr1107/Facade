package server

// Server serves requests for content.
type Server interface {
}

// FacadeServer contains all data and embeddings a Server would need.
type FacadeServer struct {
	Server
	Name    string
	Address string
	Port    int
}
