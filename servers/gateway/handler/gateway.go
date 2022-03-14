package handler

import (
	"context"
	"fmt"
	"hotwave/frame"
	frameproto "hotwave/frame/proto"
	"hotwave/servers/gateway/gate"
	"hotwave/servers/gateway/gate/codec"
	"hotwave/servers/gateway/proto"
)

type Gater struct {
	proto.UnimplementedGatewayServer

	sessions map[string]*gate.Session
	users    map[int64]*gate.Session
	// router   router.Router
}

func NewGater() *Gater {
	g := &Gater{
		sessions: make(map[string]*gate.Session),
		users:    make(map[int64]*gate.Session),
	}
	return g
}

func (g *Gater) SendSessionErrorAndClose(session gate.Session, err error) {

	session.Close()
}

func (g *Gater) OnGateMessage(session gate.Session, msg *codec.Message) {
	uid := session.UID()

	inAllowList := true
	if inAllowList {

	} else {
		if uid == 0 {
			g.SendSessionErrorAndClose(session, fmt.Errorf("auth is required"))
			return
		}
	}

	session.Send(msg)

	// router.Route(msg.Route, msg.Type, msg.Head, msg.Body)
	// session
	// rawTargetID, has := session.GetMeta(msg.Endpoint)
	// if !has {
	// 	session.SetMeta(msg.Endpoint, rawTargetID)
	// }
	// switch msg.Type {
	// case codec.Request:
	// case codec.Event:
	// case codec.Async:
	// case codec.Response:
	// default:
	// }
}

func (g *Gater) OnGateConnStat(session gate.Session, status gate.SocketStat) {

}

func (g *Gater) OnUserMessage(user frame.User, msg *frameproto.UserMessageWraper) {

}

func (g *Gater) OnNodeEvent(nodeid string, msg *frameproto.EventMessageWraper) {

}

func (g *Gater) SendMessageToUse(ctx context.Context, in *proto.SendMessageToUserRequest) (*proto.SendMessageToUserResponse, error) {
	out := &proto.SendMessageToUserResponse{}
	return out, nil
}
