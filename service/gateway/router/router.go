package router

import (
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type Router interface {
	Forward(from transport.Session, msg *protocal.ForwardMessageWarp)
	OnSessionStat(transport.Session, transport.SessionStat)
}
