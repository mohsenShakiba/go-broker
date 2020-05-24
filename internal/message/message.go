package message

type Message struct {
	Type   string
	Fields map[string][]byte
}
