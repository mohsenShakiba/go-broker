package manager

import (
	"fmt"
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

	subConfig := subscriberConfig{
		parallelism: parallelism,
		routes:      routes,
	}

	subscriber := NewSubscriber(ctx.Client, subConfig)

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

	p := &PayloadMessage{
		Id:      ctx.Message.MsgId,
		Routes:  routes,
		Payload: payloadContent,
	}

	m.processMessage(p)
}

func (m *Manager) handleAck(ctx *tcp.Context) {

}

func (m *Manager) handleNack(ctx *tcp.Context) {

}
