package packer

type Packer interface {

	// Pack method will append the payload length prefix to the payload
	Pack(p []byte) []byte

	// Unpack will remove the message length prefix and return the actual payload
	Unpack(p []byte) []byte

	// MessageLength will return the length of a payload
	PayloadLength(p []byte) int
}
