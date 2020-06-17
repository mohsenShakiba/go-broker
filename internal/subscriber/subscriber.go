package subscriber

import (
	"go-broker/internal/internal/rate_controller"
	"go-broker/internal/models"
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
		rController: rate_controller.New(config.parallelism),
	}
}

type Subscriber struct {
	client      *tcp.Client
	config      subscriberConfig
	rController rate_controller.RateController
}

func (s *Subscriber) OnAck(msgId string) {
	s.rController.ReleaseOne(msgId)
}

func (s *Subscriber) OnNack(msgId string) {
	s.rController.ReleaseOne(msgId)
}

func (s *Subscriber) OnMessage(msg *models.Message) {
	s.rController.WaitOne(msg.Id)
	s.client.Write(msg.ToTcpMessage())
}
