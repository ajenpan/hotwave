package gate

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	protocol "hotwave/services/gateway/proto"

	protobuf "google.golang.org/protobuf/proto"
)

type Session interface {
	ID() string

	UID() uint64
	SetUID(uint64)

	SetMeta(string, interface{})
	GetMeta(string) (interface{}, bool)

	Send(protobuf.Message) error

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
	Async  MessageType = 2
)

type SocketStat = int32

const (
	SocketStatConnected    SocketStat = 1
	SocketStatDisconnected SocketStat = 2
)

type AsyncAdapter interface {
	OnGateAsync(Session, *protocol.ClientMessageWraper)
	OnGateConnStat(Session, SocketStat)
}

type MethodAdapter interface {
	OnGateMethod(context.Context, *protocol.ClientMessageWraper) (*protocol.ClientMessageWraper, error)
}

type GateAdapter interface {
	AsyncAdapter
	MethodAdapter
}
