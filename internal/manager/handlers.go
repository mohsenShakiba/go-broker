package manager

import (
	"go-broker/internal/tcp"
	"strings"
)

func handleSubscribeMessage(ctx *tcp.Context) {
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

	// send ack
	ctx.SendAck()

}

func handlePublishMessage(ctx *tcp.Context) {

}

func handlePublishMessage(ctx *tcp.Context) {

}
