package router

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	protocal "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type Recipient interface {
	OnMessage(s transport.Session, msg *protocal.GateClientMessage) error
}

func NewStaticRouter() *StaticRouter {
	return &StaticRouter{
		Servers: make(map[string]Recipient),
		Nodes:   make(map[string]Recipient),
	}
}

type StaticRouter struct {
	Servers map[string]Recipient
	Nodes   map[string]Recipient
}

func (r *StaticRouter) Add(name, nodeid string, d Recipient) {
	r.Nodes[nodeid] = d
	r.Servers[name] = d
}

func (s *StaticRouter) OnRouteMessage(sess transport.Session, msg *protocal.GateClientMessage) {
	name := protoreflect.FullName(msg.Name)
	server := string(name.Parent())

	var d Recipient = s.Servers[server]

	switch endpoint := msg.Endpoint.(type) {
	case *protocal.GateClientMessage_Nodeid:
		d = s.Nodes[endpoint.Nodeid]
	}

	if d == nil {
		return
	}

	if err := d.OnMessage(sess, msg); err != nil {
		fmt.Println(err)
	}
}
