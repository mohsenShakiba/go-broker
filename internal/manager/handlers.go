package manager

import (
	log "github.com/sirupsen/logrus"
	"go-broker/internal/models"
	subscriber2 "go-broker/internal/subscriber"
	"go-broker/internal/tcp"
	"strings"
)

func (m *Manager) handleSubscribeMessage(ctx *tcp.Context) {
	// read the message
	parallelism, ok := ctx.Message.ReadInt("dop")

	if !ok {
		ctx.SendErr("INVALID_PARAM")
		return
	}

	routesStr, ok := ctx.Message.ReadStr("routes")

	if !ok {
		ctx.SendErr("INVALID_PARAM")
		return
	}

	routes := strings.Split(routesStr, ",")

	subConfig := subscriber2.subscriberConfig{
		parallelism: parallelism,
		routes:      routes,
	}

	subscriber := subscriber2.NewSubscriber(ctx.Client, subConfig)

	m.router.AddRoute(routes, subscriber)

	go subscriber.start()

	// send ack
	ctx.SendAck()

}

func (m *Manager) handlePublishMessage(ctx *tcp.Context) {
	// read the message
	payloadContent, ok := ctx.Message.ReadByteArr("payload")

	if !ok {
		ctx.SendErr("INVALID_PARAM")
		return
	}

	routesStr, ok := ctx.Message.ReadStr("routes")

	if !ok {
		ctx.SendErr("INVALID_PARAM")
		return
	}

	routes := strings.Split(routesStr, ",")

	p := &models.Message{
		Id:      ctx.Message.MsgId,
		Routes:  routes,
		Payload: payloadContent,
	}

	m.processMessage(p)
}

func (m *Manager) handleAck(ctx *tcp.Context) {
	msgId := ctx.Message.MsgId

	log.Infof("ack was received for msgId: %s", msgId)

	m.processAck(msgId)
}

func (m *Manager) handleNack(ctx *tcp.Context) {

}
