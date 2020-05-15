package socketserver

//
//func parseMessage(tcpMessage []byte) (interface{}, error) {
//
//	// extract the tcpMessage type
//	msgType := tcpMessage[:3]
//
//	payload := tcpMessage[3:]
//
//	s := serializer.NewJsonSerializer()
//
//	// parse based on type
//	switch string(msgType) {
//	case authenticateMessageType:
//		msg := &authenticateMessage{}
//		err := s.Deserialize(payload, msg)
//		return msg, err
//	case routedMessageType:
//		msg := &routedMessage{}
//		err := s.Deserialize(payload, msg)
//		return msg, err
//	case subscribeMessageType:
//		msg := &subscribeMessage{}
//		err := s.Deserialize(payload, msg)
//		return msg, err
//	case ackMessageType:
//		msg := &ackMessage{}
//		err := s.Deserialize(payload, msg)
//		return msg, err
//	case nackMessageType:
//		msg := &nackMessage{}
//		err := s.Deserialize(payload, msg)
//		return msg, err
//
//	}
//
//	errMsg := fmt.Sprintf("the tcpMessage type %s cannot be parsed", string(msgType))
//	return nil, errors.New(errMsg)
//
//}
