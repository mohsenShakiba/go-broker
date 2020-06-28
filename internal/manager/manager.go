package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/channel"
	"go-broker/internal/models"
	"go-broker/internal/subscriber"
	"go-broker/internal/tcp"
	"sync"
)

type Manager struct {
	socketServer *tcp.Server
	msgMap       map[string]*channel.Channel
	chanMap      map[string]*channel.Channel
	lock         sync.RWMutex
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

	// create manager
	mgr := &Manager{
		conf:         conf,
		socketServer: socketServer,
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
		case *models.Ping:
			m.processPing(ctx.Client, msg)
		}
	}

}

func (m *Manager) processMessage(client *tcp.Client, msg *models.Message) {

	// check if channel exits
	ch, ok := m.chanMap[msg.Route]

	// create a new channel
	if !ok {

		ch = channel.NewChannel(msg.Route, m.conf.StorageConfig)

		err := ch.Init()

		if err != nil {
			log.Fatalf("failed to initialize channel, error: %s", err)
			return
		}

	}

	// add to mapping
	m.lock.Lock()
	defer m.lock.Unlock()
	m.msgMap[msg.Id] = ch
	m.chanMap[msg.Route] = ch

	// enqueue message
	ch.Enqueue(msg)

	// send message ack
	client.SendAck(msg.Id)

}

func (m *Manager) processAck(client *tcp.Client, ack *models.Ack) {

	// check if channel exits
	m.lock.RLock()
	ch, ok := m.msgMap[ack.Id]
	m.lock.RUnlock()

	// return error
	if !ok {
		client.SendError(ack.Id, "invalid route")
		return
	}

	// send ack
	ch.Ack(ack)
}

func (m *Manager) processNack(client *tcp.Client, nack *models.Nack) {

	// check if channel exits
	m.lock.RLock()
	ch, ok := m.chanMap[nack.Id]
	m.lock.RUnlock()

	// return error
	if !ok {
		client.SendError(nack.Id, "invalid route")
		return
	}

	// send ack
	ch.Nack(nack)
}

func (m *Manager) processSubscribe(client *tcp.Client, reg *models.Register) {

	// subscriber config
	conf := subscriber.Config{
		Dop:    reg.Dop,
		Routes: reg.Routes,
	}

	// create a new subscription
	sub := subscriber.NewSubscriber(client, conf)

	// create channels if not exists
	// add subscriber to channel
	for _, route := range reg.Routes {

		m.lock.RLock()
		ch, ok := m.chanMap[route]
		m.lock.RUnlock()

		if !ok {

			ch = channel.NewChannel(route, m.conf.StorageConfig)

			err := ch.Init()

			if err != nil {
				log.Fatalf("failed to initialize channel, error: %s", err)
				return
			}

			m.lock.Lock()
			m.chanMap[route] = ch
			m.lock.Unlock()
		}

		ch.Register(sub)
	}

	log.Infof("added new subscriber with id: %s", client.ClientId)

	// send ack
	client.SendAck(reg.Id)

}

func (m *Manager) processPing(client *tcp.Client, p *models.Ping) {
	client.SendAck(p.Id)
}
