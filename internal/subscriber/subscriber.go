package subscriber

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/internal/rate_controller"
	"go-broker/internal/models"
	"go-broker/internal/tcp"
	"io"
)

type Config struct {
	Dop    int
	Routes []string
}

func NewSubscriber(c *tcp.Client, config Config) *Subscriber {
	return &Subscriber{
		client:      c,
		config:      config,
		rController: rate_controller.New(config.Dop),
	}
}

type Subscriber struct {
	client      *tcp.Client
	config      Config
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
	s.client.BeginWrite(func(w io.Writer) {
		err := msg.Write(w)
		if err != nil {
			log.Errorf("failed to write the message, error: %s", err)
		}
	})
}
