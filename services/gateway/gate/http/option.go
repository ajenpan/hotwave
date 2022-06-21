package websocket

import (
	"hotwave/services/gateway/gate"
)

type Option func(*Options)

type Options struct {
	Address string
	Adapter gate.MethodAdapter
}
