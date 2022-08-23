package router

import (
	"fmt"
	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
	"strconv"
	"sync"
)

func NewStaticRouter() *StaticRouter {
	return &StaticRouter{
		// Servers: make(map[string]Recipient),
		// Nodes:   make(map[string]Recipient),
	}
}

type StaticRouter struct {
	// Servers map[string]Recipient
	// Nodes   map[string]Recipient

	// user sync.Map
	// sockets sync.Map
	// nodes sync.Map

	sessions sync.Map
}

func (r *StaticRouter) AddSession(sch, nodeid string, s transport.Session) {

}

func (s *StaticRouter) Forward(from transport.Session, msg *protocal.ForwardMessageWarp) error {
	var sessionKey string

	switch dest := (msg.Dest.Endpoint).(type) {
	case *protocal.ForwardEndpoint_Uid:
		sessionKey = "uid://" + strconv.FormatInt(dest.Uid, 10)
	case *protocal.ForwardEndpoint_Socketid:
		sessionKey = "socket://" + dest.Socketid
	default:

	}

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
	// name := protoreflect.FullName(msg.MsgName)
	// server := string(name.Parent())
	// var d Recipient
	// if msg.Nodeid != "" {
	// 	d = s.Nodes[msg.Nodeid]
	// } else {
	// 	d = s.Servers[server]
	// }
	// if d == nil {
	// 	return
	// }
	// if err := d.OnMessage(sess, msg); err != nil {
	// 	fmt.Println(err)
	// }
}
