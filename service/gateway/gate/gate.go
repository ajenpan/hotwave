package gate

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"hotwave/service/gateway/proto"
	"hotwave/service/gateway/router"
	"hotwave/transport"
	"hotwave/transport/tcp"
)

type TcpGater struct {
	Name             string
	Route            router.Router
	Address          string
	HeatbeatInterval time.Duration
	tcpLinster       *tcp.Server
}

func (g *TcpGater) Start() error {
	g.tcpLinster = tcp.NewServer(tcp.ServerOptions{
		HeatbeatInterval: g.HeatbeatInterval,
		Address:          g.Address,
		OnMessage:        g.onGateMessage,
		OnConn:           g.onStat,
		NewIDFunc:        transport.NewSessionID,
	})
	return g.tcpLinster.Start()
}

func (g *TcpGater) Stop() {
	if g.tcpLinster == nil {
		return
	}
	g.tcpLinster.Stop()
}

func (g *TcpGater) onStat(socket transport.Session, status transport.SessionStat) {
	g.Route.OnSessionStat(socket, status)
}

func (g *TcpGater) onGateMessage(session transport.Session, iraw interface{}) {
	msg := &proto.GateMessage{}
	switch raw := iraw.(type) {
	case []byte:
		protobuf.Unmarshal(raw, msg)
	case *proto.GateMessage:
		msg = raw
	default:
		return
	}

	warp := &proto.ForwardMessageWarp{
		Src: &proto.ForwardEndpoint{
			EndpointType: g.Name,
			EndpointId:   session.ID(),
		},
		Dest:    &proto.ForwardEndpoint{},
		MsgName: msg.Name,
		Data:    msg.Body,
		Type:    msg.Type,
		Mime:    msg.Mime,
	}

	g.Route.Forward(session, warp)
}
