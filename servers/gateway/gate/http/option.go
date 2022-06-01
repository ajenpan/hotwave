package websocket

import (
	"hotwave/servers/gateway/gate"
)

type Option func(*Options)

type Options struct {
	Address string
	Adapter gate.MethodAdapter
}
