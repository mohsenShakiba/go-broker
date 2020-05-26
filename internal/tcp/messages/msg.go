package messages

// Message contain information about an incoming tcp messages
type Message struct {
	Type   string
	MsgId  string
	Fields map[string][]byte
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

func (m *Message) WriteStr(key string, value string) {
	m.Fields[key] = []byte(value)
}

func (m *Message) WriteByte(key string, value []byte) {
	m.Fields[key] = value
}
