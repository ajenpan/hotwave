package router

import (
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type Router interface {
	Forward(msg *protocal.GateMessage, src transport.Session) error
}
