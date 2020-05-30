package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/storage"
	"go-broker/internal/subscribe"
	"go-broker/internal/tcp"
)

type Manager struct {
	socketServer *tcp.Server
	storage      *storage.Storage
	router       *Router
}

func InitManager(basePath string) (*Manager, error) {

	// create socket server config
	socketServerConf := tcp.ServerConfig{
		ConnectionPort: 8085,
	}

	// create socket server
	socketServer := tcp.New(socketServerConf)

	// start socket server
	socketServer.Start()

	// register handlers
	socketServer.RegisterHandler("SUB", handleSubscribeMessage)
	socketServer.RegisterHandler("PUB", handlePublishMessage)
	socketServer.RegisterHandler("ACK", handlePublishMessage)

	// init router
	router := NewRouter()

	// init storage
	s, err := storage.Init(basePath)

	if err != nil {
		return nil, err
	}

	mgr := &Manager{
		socketServer: socketServer,
		storage:      s,
		router:       router,
	}

	go mgr.processPublishedMessage(publishMessageChan, subscriberChan)

	return mgr, nil
}

func (m *Manager) processPublishedMessage(publishedChan chan *publish.PublishedMessage, subChan chan *subscribe.PublishedMessage) {
	for {
		msg := <-publishedChan

		log.Infof("publishing messages with id %s to subscribers", msg.MsgId)

		subChan <- &subscribe.PublishedMessage{
			MsgId:   msg.MsgId,
			Payload: msg.Payload,
			Routes:  msg.Routes,
		}
	}
}
