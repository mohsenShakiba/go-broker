package message

type Message struct {
	Type   string
	Fields map[string][]byte
}

func (m *Message) ReadMsgId() (string, bool) {
	return m.ReadStr("msgId")
}

func (m *Message) ReadStr(key string) (string, bool) {
	value := m.Fields[key]

	if value == nil {
		return "", false
	}

	return string(value), true
}

func (m *Message) ReadByteArr(key string) ([]byte, bool) {
	value, ok := m.Fields[key]

	if !ok {
		return nil, false
	}

	return value, true
}
