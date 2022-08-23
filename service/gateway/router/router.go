package router

import (
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type Router interface {
	Forward(from transport.Session, msg *protocal.ForwardMessageWarp) error
}

type Recipient interface {
	OnMessage(s transport.Session, msg *protocal.ForwardMessageWarp) error
}
