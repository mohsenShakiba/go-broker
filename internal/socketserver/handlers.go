package socketserver

import log "github.com/sirupsen/logrus"

func processRoutedMessageHandler(msg *tcpMessage, ch *chan ServerEvents) {

	id, ok := msg.readStr("id")

	if !ok {
		log.Warnf("the routed message doesn't have a valid id")
		return
	}

	route, ok := msg.readStr("route")

	if !ok {
		log.Warnf("the routed message doesn't have a valid route")
		return
	}

	payload, ok := msg.readByteArr("payload")

	if !ok {
		log.Warnf("the routed message doesn't have a valid payload")
		return
	}
}

func processAckMessageHandler(msg *tcpMessage, ch *chan ServerEvents) {

}

func processNackMessageHandler(msg *tcpMessage, ch *chan ServerEvents) {

}

func processSubscribeMessageHandler(msg *tcpMessage, ch *chan ServerEvents) {

}
