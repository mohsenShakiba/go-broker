package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/storage"
	"go-broker/internal/tcp"
	"sync"
)

type Manager struct {
	socketServer   *tcp.Server
	storage        storage.Storage
	router         *Router
	messageMapping map[string]*Subscriber
	lock           sync.Mutex
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

	// init router
	router := NewRouter()

	storageConfig := storage.StorageConfig{
		Path:           basePath,
		FileMaxSize:    1024,
		FileNamePrefix: "go",
	}

	// init storage
	s := storage.New(storageConfig)

	err := s.Init()

	if err != nil {
		return nil, err
	}

	mgr := &Manager{
		socketServer:   socketServer,
		storage:        s,
		router:         router,
		messageMapping: make(map[string]*Subscriber),
	}

	// register handlers
	socketServer.RegisterHandler("SUB", mgr.handleSubscribeMessage)
	socketServer.RegisterHandler("PUB", mgr.handlePublishMessage)
	socketServer.RegisterHandler("ACK", mgr.handleAck)

	return mgr, nil
}

func (m *Manager) processMessage(p *PayloadMessage) {

	log.Infof("processing message with id %s", p.Id)

	subscribers := m.router.Match(p.Routes)

	msgB, err := p.ToBinary()

	if err != nil {
		log.Errorf("failed to serialize message, error: %s", err)
	}

	err = m.storage.Write(getStringHash(p.Id), msgB)

	if err != nil {
		log.Errorf("failed to persist message, error: %s", err)
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	for _, s := range subscribers {
		s.OnMessage(p)
		m.messageMapping[p.Id] = s
	}
}

func (m *Manager) processAck(msgId string) {
	m.lock.Lock()
	s, ok := m.messageMapping[msgId]
	m.lock.Unlock()

	err := m.storage.Delete(getStringHash(msgId))

	if err != nil {
		log.Errorf("failed to persist ack, error: %s", err)
	}

	log.Infof("processing ack for msgId: %s", msgId)

	if !ok {
		return
	}

	s.OnAck(msgId)
}
