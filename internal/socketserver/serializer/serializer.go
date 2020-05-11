package serializer

type Serializer interface {
	Serialize(msg interface{}) ([]byte, error)
	Deserialize(payload []byte)
}
