package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/models"
	"go-broker/internal/storage"
	"go-broker/internal/subscriber"
	"go-broker/internal/tcp"
	"sync"
)

type Manager struct {
	socketServer   *tcp.Server
	storage        storage.Storage
	router         *Router
	messageMapping map[string]*subscriber.Subscriber
	lock           sync.Mutex
	conf           Config
}

func InitManager(conf Config) (*Manager, error) {

	// create socket server config
	socketServerConf := tcp.ServerConfig{
		ConnectionPort: conf.Port,
	}

	// create msg chan
	msgChan := make(chan *tcp.Context)

	// create socket server
	socketServer := tcp.New(socketServerConf, msgChan)

	// start socket server
	socketServer.Start()

	// init router
	router := NewRouter()

	// init storage
	store := storage.NewStorage(conf.FilePath, conf.StorageType)
	err := store.Init()

	if err != nil {
		return nil, err
	}

	// create manager
	mgr := &Manager{
		socketServer:   socketServer,
		storage:        store,
		router:         router,
		messageMapping: make(map[string]*subscriber.Subscriber),
	}

	// process incoming message
	go mgr.processMessage(msgChan)

	return mgr, nil
}

func (m *Manager) processMessage(ch chan *tcp.Context) {

	for {
		msg := <-ch
		switch msg.(type) {

		}
	}

	log.Infof("processing message with id %s", p.Id)

	subscribers := m.router.Match(p.Routes)

	msgB, err := p.ToBinary()

	if err != nil {
		log.Errorf("failed to serialize message, error: %s", err)
	}

	err = m.storage.Write(models.getStringHash(p.Id), msgB)

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

	err := m.storage.Delete(models.getStringHash(msgId))

	if err != nil {
		log.Errorf("failed to persist ack, error: %s", err)
	}

	log.Infof("processing ack for msgId: %s", msgId)

	if !ok {
		return
	}

	s.OnAck(msgId)
}
