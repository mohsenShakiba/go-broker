package manager

// Router is in charge of routing incoming messages to appropriate subscribers
//based on the message and subscriber routes
type Router struct {
	Routes []*Route
}

func NewRouter() *Router {
	return &Router{
		Routes: make([]*Route, 0),
	}
}

func (r *Router) AddRoute(routes []string, client *Subscriber) {

	for _, route := range routes {
		route := &Route{
			Path:       route,
			segments:   ProcessPath(route),
			cache:      make(map[string]bool),
			subscriber: client,
		}
		r.Routes = append(r.Routes, route)
	}

}

func (r *Router) Match(path []string) map[string]*Subscriber {
	clients := make(map[string]*Subscriber, 0)

	for _, r := range r.Routes {
		for _, p := range path {
			if r.DoesMatch(p) {
				clients[r.subscriber.client.ClientId] = r.subscriber
			}
		}
	}

	return clients
}
