package node

//TODO:

type Message struct {
	Body interface{}
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
