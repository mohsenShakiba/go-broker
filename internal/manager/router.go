package manager

import "go-broker/internal/tcp"

// Router is in charge of routing incoming messages to appropriate subscribers
//based on the message and subscriber routes
type Router struct {
	Routes []*Route
}

func NewRouter() *Router {
	return &Router{
		Routes: make([]Route, 0),
	}
}

func (r Router) AddRoute(path string, client *tcp.Client) {
	route := &Route{
		Path:     path,
		segments: ProcessPath(path),
		cache:    make(map[string]bool),
		Client:   client,
	}
	r.Routes = append(r.Routes, route)
}

func (r *Router) Match(path string) []*tcp.Client {
	clients := make([]*tcp.Client, 0)

	for _, r := range r.Routes {
		if r.DoesMatch(path) {
			clients = append(clients, r.Client)
		}
	}

	return clients
}
