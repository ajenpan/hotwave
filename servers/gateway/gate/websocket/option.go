package websocket

import (
	"time"

	"hotwave/servers/gateway/gate"
)

type Option func(*Options)

type Options struct {
	Address string
	// The interval on which to register
	HeatbeatInterval time.Duration
	Adapter          gate.AsyncAdapter
}
