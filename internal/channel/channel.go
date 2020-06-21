package channel

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/internal/queue"
	"go-broker/internal/models"
	"go-broker/internal/storage"
	"go-broker/internal/subscriber"
	"sync"
	"time"
)

type ChannelOptions struct {
	FilePath    string
	StorageType string
}

type Channel struct {

	// name of channel
	name string

	// channel options
	opt ChannelOptions

	// list of available subscribers
	subscribers []*subscriber.Subscriber

	// message queue
	// queue is by default
	msgQueue queue.Queue

	// mapping to know which message was sent by which subscriber
	messageMap map[string]*subscriber.Subscriber

	// to know which subscriber was last used for sending message
	// this field is required for enabling round robin
	sIndex int

	storage storage.Storage

	lock sync.RWMutex
}

func NewChannel(route string, opt ChannelOptions) *Channel {
	return &Channel{
		name:        route,
		opt:         opt,
		subscribers: make([]*subscriber.Subscriber, 0),
		msgQueue:    queue.New(),
		messageMap:  make(map[string]*subscriber.Subscriber),
		sIndex:      0,
		storage:     nil,
		lock:        sync.RWMutex{},
	}
}

func (c *Channel) Init() error {
	c.storage = storage.NewStorage(c.opt.FilePath, c.opt.StorageType)
	go c.processMessages()
	return c.storage.Init()
}

func (c *Channel) processMessages() {
	for {
		// retrieve messages
		item := c.msgQueue.Dequeue()

		// convert to payload
		msg := item.(*models.Message)

		// check if any subscriber is available
		if len(c.subscribers) <= 0 {
			time.Sleep(time.Second)
			continue
		}

		// lock
		c.lock.RLock()
		c.lock.RUnlock()

		// rotate the sIndex
		// this is required for enabling load balancing between subscribers
		if c.sIndex >= len(c.subscribers) {
			c.sIndex = 0
		}

		// find the next available subscriber
		subscriber := c.subscribers[c.sIndex]

		// add subscriber to message map
		c.lock.Lock()
		c.messageMap[msg.Id] = subscriber
		c.lock.Unlock()

		// send message
		subscriber.OnMessage(msg)
	}
}

func (c *Channel) Enqueue(msg *models.Message) {

	msgb, err := msg.ToBinary()

	if err != nil {
		log.Errorf("failed to serialize the msg, error: %s", err)
	}

	c.msgQueue.Enqueue(msg)

	err = c.storage.Write(msg.Id, msgb)

	if err != nil {
		log.Errorf("failed to modify storage, error: %s", err)
	}

}

func (c *Channel) Ack(m *models.Ack) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if s, ok := c.messageMap[m.Id]; ok {
		s.OnAck(m.Id)
	}

	delete(c.messageMap, m.Id)

	err := c.storage.Delete(m.Id)

	if err != nil {
		log.Errorf("failed to modify storage, error: %s", err)
	}
}

func (c *Channel) Nack(m *models.Nack) {

	c.lock.Lock()
	defer c.lock.Unlock()

	if s, ok := c.messageMap[m.Id]; ok {
		s.OnNack(m.Id)
	}

	delete(c.messageMap, m.Id)

	msgb, err := c.storage.Read(m.Id)

	if err != nil {
		log.Errorf("failed to find the id specified")
		return
	}

	msg := models.Message{}
	err = msg.FromBinary(msgb)

	if err != nil {
		log.Errorf("failed to deserialize message")
		return
	}

	c.msgQueue.Enqueue(msg)
}

func (c *Channel) Register(sub *subscriber.Subscriber) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.subscribers = append(c.subscribers, sub)
}
