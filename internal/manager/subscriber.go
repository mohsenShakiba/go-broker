package manager

import (
	"go-broker/internal/manager/internal/queue"
	"go-broker/internal/manager/internal/rate_controller"
	"go-broker/internal/tcp"
)

type subscriberConfig struct {
	parallelism int
	routes      []string
}

func NewSubscriber(c *tcp.Client, config subscriberConfig) *Subscriber {
	return &Subscriber{
		client:      c,
		config:      config,
		queue:       queue.New(),
		rController: rate_controller.New(config.parallelism),
	}
}

type Subscriber struct {
	client      *tcp.Client
	config      subscriberConfig
	queue       queue.Queue
	rController rate_controller.RateController
}

// start will start sending messages from the queue
func (s *Subscriber) start() {

	for {

		// retrieve messages
		item := s.queue.Dequeue()

		// convert to payload
		msg := item.(PayloadMessage)

		// pass it through the rate controller
		s.rController.WaitOne(msg.Id)

	}
}

func (s *Subscriber) OnAck(msgId string) {
	s.rController.ReleaseOne(msgId)
}

func (s *Subscriber) OnNack(msgId string) {
	s.rController.ReleaseOne(msgId)
}

func (s *Subscriber) OnMessage(message *PayloadMessage) {
	s.queue.Enqueue(message)
}
