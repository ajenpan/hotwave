package gate

import "hotwave/servers/gateway/gate/codec"

type Session interface {
	ID() string
	UID() int64
	SetUID(int64)

	SetMeta(string, interface{})
	GetMeta(string) (interface{}, bool)

	Send(*codec.Message) error
	Close()
}

type MessageType = int

const (
	Method MessageType = 1
	Stream MessageType = 2
)

type SocketStat = int32

const (
	SocketStatConnected    SocketStat = 1
	SocketStatDisconnected SocketStat = 2
)

type OnMessageFunc func(*Session, *codec.Message)
type OnConnStatFunc func(*Session, SocketStat)

type Adapter interface {
	OnMessage(Session, *codec.Message)
	OnConnStat(Session, SocketStat)
}
