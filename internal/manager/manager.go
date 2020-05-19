package manager

import (
	"go-broker/internal/publish"
	"go-broker/internal/subscribe"
	"go-broker/internal/tcp"
)

type Manager struct {
	socketServer      *tcp.Server
	publisherManager  *publish.PublisherManager
	subscriberManager *subscribe.SubscriberManager
}

func InitManager() {
	mgr := Manager{
		socketServer:      nil,
		publisherManager:  nil,
		subscriberManager: nil,
	}
}
