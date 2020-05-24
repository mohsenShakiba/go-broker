package tcp

//
//import (
//	"bytes"
//	"errors"
//	log "github.com/sirupsen/logrus"
//	"go-broker/internal/serializer"
//	"io"
//	"strconv"
//)
//
//type MessageContext struct {
//	ClientId   string
//	PayloadMap map[string][]byte
//	client     io.WriteCloser
//	Serializer *serializer.LineSeparatedSerializer
//}
//
//func convertToMessage(b []byte, client *Client) *MessageContext {
//
//	newLineB := []byte("\n")
//	colonB := []byte(":")
//
//	tcpMsg := &MessageContext{
//		PayloadMap: make(map[string][]byte),
//	}
//
//	// split by new line
//	partsByNewLine := bytes.Split(b, newLineB)
//
//	for _, part := range partsByNewLine {
//		partsByColon := bytes.Split(part, colonB)
//
//		if len(partsByColon) != 2 {
//			log.Warnf("bad payload data, discarding, message: %s", string(part))
//		}
//
//		tcpMsg.PayloadMap[string(partsByColon[0])] = partsByColon[1]
//	}
//
//	tcpMsg.ClientId = client.ClientId
//	tcpMsg.client = client
//	tcpMsg.Serializer = serializer.NewLineSeparatedSerializer()
//
//	return tcpMsg
//}
//
//func (m *MessageContext) ReadByteArr(key string) ([]byte, bool) {
//	value := m.PayloadMap[key]
//
//	if value == nil {
//		return nil, false
//	}
//
//	return value, true
//}
//
//func (m *MessageContext) ReadStr(key string) (string, bool) {
//	value := m.PayloadMap[key]
//
//	if value == nil {
//		return "", false
//	}
//
//	return string(value), true
//}
//
//func (m *MessageContext) ReadInt(key string) (int, bool) {
//	value := m.PayloadMap[key]
//
//	if value == nil {
//		return 0, false
//	}
//
//	num, err := strconv.Atoi(string(value))
//
//	if err != nil {
//		return 0, false
//	}
//
//	return num, true
//}
//
//func (m *MessageContext) GetMessageId() (string, bool) {
//	return m.ReadStr("msgId")
//}
//
//func (m *MessageContext) GetMessageType() (string, bool) {
//	return m.ReadStr("type")
//}
//
//func (m *MessageContext) SendAck() error {
//
//	msgId, ok := m.GetMessageId()
//
//	if !ok {
//		return errors.New("the message id isn't provided")
//	}
//
//	m.Serializer.WriteStr("type", "ACK")
//	m.Serializer.WriteStr("msgId", msgId)
//
//	log.Infof("sending ack to client %s", m.ClientId)
//
//	_, err := m.client.Write(m.Serializer.GetMessageBytes())
//
//	if err != nil {
//		return err
//	}
//
//	//for _, b := range m.Serializer.Bytes {
//	//	_, err = m.client.Write(b)
//	//
//	//	if err != nil {
//	//		return err
//	//	}
//	//
//	//}
//
//	return nil
//}
//
//func (m *MessageContext) Close() {
//	_ = m.client.Close()
//}
