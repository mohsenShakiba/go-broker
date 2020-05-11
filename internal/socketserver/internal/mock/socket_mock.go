package mock

import "io"

type socketMock struct {
	buffer   []byte
	IsClosed bool
}

func NewMockSocket() io.ReadWriteCloser {
	return &socketMock{
		buffer: make([]byte, 0),
	}
}

func (s *socketMock) Write(p []byte) (n int, err error) {
	s.buffer = p
	return len(s.buffer), nil
}

func (s *socketMock) Read(p []byte) (n int, err error) {
	copy(p, s.buffer)

	defer func() {
		s.buffer = make([]byte, 0)
	}()

	return len(s.buffer), nil
}

func (s *socketMock) Close() error {
	s.IsClosed = true
	return nil
}
