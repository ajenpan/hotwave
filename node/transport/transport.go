package transport

import "hotwave/marshal"

type Message struct {
	Body      interface{}
	Marshaler marshal.Marshaler
}

func (m *Message) Encode() ([]byte, error) {
	return m.Marshaler.Marshal(m.Body)
}

func (m *Message) Decode(data []byte) error {
	return m.Marshaler.Unmarshal(data, m.Body)
}

func (m *Message) ContentType() string {
	return m.Marshaler.ContentType(m)
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	LocalAddr() string
	RemoteAddr() string
}

type Client interface {
	Socket
}

type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}
