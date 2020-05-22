package manager

import (
	log "github.com/sirupsen/logrus"
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

func InitManager(basePath string) (*Manager, error) {

	publishMessageChan := make(chan *publish.PublishedMessage)
	subscriberChan := make(chan *subscribe.PublishedMessage)

	socketServer := tcp.Init(tcp.TcpConfig{
		Credentials:    []string{""},
		ConnectionPort: 8085,
	})

	publisherManager := publish.InitPublisherManager(publishMessageChan, socketServer)

	subscriberManager := subscribe.InitSubscriberManager(socketServer, subscriberChan)

	storage, err := storage.Init(basePath)

	if err != nil {
		return nil, err
	}

	mgr := &Manager{
		socketServer:      socketServer,
		publisherManager:  publisherManager,
		subscriberManager: subscriberManager,
		storage:           storage,
	}

	go mgr.processPublishedMessage(publishMessageChan, subscriberChan)

	return mgr, nil
}

func (m *Manager) processPublishedMessage(publishedChan chan *publish.PublishedMessage, subChan chan *subscribe.PublishedMessage) {
	for {
		msg := <-publishedChan

		log.Infof("publishing message with id %s to subscribers", msg.MsgId)

		subChan <- &subscribe.PublishedMessage{
			MsgId:   msg.MsgId,
			Payload: msg.Payload,
			Routes:  msg.Routes,
		}
	}
}
