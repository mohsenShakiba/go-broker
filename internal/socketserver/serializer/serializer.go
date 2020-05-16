package serializer

import (
	"encoding/json"
	"errors"
)

type Serializer interface {
	Serialize(msg interface{}) ([]byte, error)
	Deserialize(payload []byte, msg interface{}) error
}

func NewJsonSerializer() Serializer {
	return &jsonSerializer{}
}

type jsonSerializer struct {
}

func (s *jsonSerializer) Serialize(msg interface{}) ([]byte, error) {

	if msg == nil {
		return nil, errors.New("nil message cannot be deserialized")
	}

	b, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *jsonSerializer) Deserialize(payload []byte, msg interface{}) error {

	if msg == nil {
		return errors.New("nil message cannot be deserialized")
	}

	err := json.Unmarshal(payload, msg)

	return err
}
