package gate

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"hotwave/servers/gateway/gate/codec"
)

type Session interface {
	ID() string

	UID() uint64
	SetUID(uint64)

	SetMeta(string, interface{})
	GetMeta(string) (interface{}, bool)

	Send(*codec.AsyncMessage) error
	Close()

	sync.Locker
}

var sid int64 = 0

func NewSessionID(prex string) string {
	return fmt.Sprintf("%s_%d_%d", prex, atomic.AddInt64(&sid, 1), time.Now().Unix())
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

type AsyncAdapter interface {
	OnGateMessage(Session, *codec.AsyncMessage)
	OnGateConnStat(Session, SocketStat)
}
