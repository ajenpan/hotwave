package router

import (
	"fmt"

	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
	"sync"
)

func NewStaticRouter() *StaticRouter {
	return &StaticRouter{}
}

type StaticRouter struct {
	sessions sync.Map
}

func (r *StaticRouter) StoreSession(stype string, s transport.Session) {
	r.sessions.Store(stype+s.ID(), s)
}

func (r *StaticRouter) RemoveSession(stype string, s transport.Session) {
	r.sessions.Delete(stype + s.ID())
}

func (s *StaticRouter) Forward(msg *protocal.ForwardMessageWarp) error {
	if msg.Dest == nil {
		return fmt.Errorf("dest is nil")
	}
	sessionKey := msg.Dest.EndpointType + msg.Dest.EndpointId
	if len(sessionKey) == 0 {
		return fmt.Errorf("key is empty")
	}
	v, has := s.sessions.Load(sessionKey)
	if !has {
		return fmt.Errorf("not found")
	}
	to, ok := v.(transport.Session)
	if !ok {
		return fmt.Errorf("value is  not session")
	}
	return to.Send(msg)
}
