package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/channel"
	"go-broker/internal/models"
	"go-broker/internal/storage"
	"go-broker/internal/subscriber"
	"go-broker/internal/tcp"
	"path"
	"sync"
)

type Manager struct {
	socketServer *tcp.Server
	router       *Router
	msgMap       map[string]*channel.Channel
	chanMap      map[string]*channel.Channel
	lock         sync.Mutex
	conf         Config
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

	// create manager
	mgr := &Manager{
		socketServer: socketServer,
		router:       router,
		msgMap:       make(map[string]*channel.Channel),
		chanMap:      make(map[string]*channel.Channel),
	}

	// process incoming message
	go mgr.process(msgChan)

	return mgr, nil
}

func (m *Manager) process(ch chan *tcp.Context) {

	for {
		ctx := <-ch
		switch msg := ctx.Message.(type) {
		case *models.Message:
			m.processMessage(ctx.Client, msg)
		case *models.Ack:
			m.processAck(ctx.Client, msg)
		case *models.Nack:
			m.processNack(ctx.Client, msg)
		case *models.Register:
			m.processSubscribe(ctx.Client, msg)
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

func (m *Manager) processMessage(client *tcp.Client, msg *models.Message) {
	// check if channel exits
	ch, ok := m.chanMap[msg.Route]

	// create a new channel
	if !ok {

		fPath := path.Join(m.conf.FilePath, msg.Route)

		ch = channel.NewChannel(msg.Route, channel.ChannelOptions{
			FilePath:    fPath,
			StorageType: m.conf.StorageType,
		})

		m.chanMap[msg.Route] = ch
	}

	// add to mapping
	m.msgMap[msg.Id] = msg.Route

	// enqueue message
	ch.Enqueue(msg)

}

func (m *Manager) processAck(client *tcp.Client, ack *models.Ack) {

	// check if channel exits
	ch, ok := m.chanMap[ack.Id]

	// create a new channel
	if !ok {
	}

	// enqueue message
	ch.Enqueue(msg)
}

func (m *Manager) processNack(msgId string) {
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

func (m *Manager) processSubscribe(msgId string) {
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
