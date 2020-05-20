package manager

import (
	"go-broker/internal/publish"
	"go-broker/internal/storage"
	"go-broker/internal/subscribe"
	"go-broker/internal/tcp"
)

type Manager struct {
	socketServer      *tcp.Server
	publisherManager  *publish.PublisherManager
	subscriberManager *subscribe.SubscriberManager
	storage           *storage.Storage
}

func InitManager() *Manager {

	publishMessageChan := make(chan *publish.PublishedMessage)
	subscriberChan := make(chan *subscribe.PublishedMessage)

	socketServer := tcp.Init(tcp.TcpConfig{
		Credentials:    []string{""},
		ConnectionPort: 8085,
	})

	publisherManager := publish.InitPublisherManager(publishMessageChan, socketServer)

	subscriberManager := subscribe.InitSubscriberManager(socketServer, subscriberChan)

	storage, err := storage.Init("")

	mgr := &Manager{
		socketServer:      socketServer,
		publisherManager:  publisherManager,
		subscriberManager: subscriberManager,
		storage:           storage,
	}

	return mgr
}
