package handler

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	"hotwave/registry"
	"hotwave/servers/gateway/gate"
	"hotwave/servers/gateway/gate/codec"
	"hotwave/servers/gateway/proto"
)

type Gater struct {
	proto.UnimplementedGatewayServer

	// sessions sync.Map //map[string]*gate.Session
	user2session sync.Map //map[uint64]*gate.Session
	// router   router.Router
}

func NewGater() *Gater {
	g := &Gater{}
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

	converMsg := &proto.UserMessageWraper{
		Name:         msg.Name,
		Meta:         msg.Meta,
		Body:         msg.Body,
		UserId:       session.UID(),
		SerialNumber: fmt.Sprintf("%s-%d", session.ID(), session.SerialNumber()),
	}

	if msg.Name == "gateway" {
		g.OnUserMessage(converMsg)
	} else {
		var target *registry.Node = nil
		if len(msg.Nodeid) == 0 {

		} else {
			// target = frame.GetNode(msg.Nodeid)
		}

		// store
		var client proto.GateAdapterClient
		if target == nil {
			conn, err := grpc.Dial(target.Address)
			if err != nil {
				return
			}

			client = proto.NewGateAdapterClient(conn)

		}
		client.UserMessage(context.Background(), converMsg)
	}

	//TODO:
	// session.Send(msg)

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

func (g *Gater) OnUserMessage(msg *proto.UserMessageWraper) {
	// var err error
	if msg.Name == "gateway" {
		fmt.Println(msg.Body)
	}
}

func (g *Gater) SendMessageToUse(ctx context.Context, in *proto.SendMessageToUserRequest) (*proto.SendMessageToUserResponse, error) {
	out := &proto.SendMessageToUserResponse{}
	v, has := g.user2session.Load(in.Uid)
	if !has {
		return out, nil
	}
	session, ok := v.(gate.Session)
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	msg := &codec.Message{}
	err := session.Send(msg)

	return out, err
}

func (g *Gater) OnLogin(session gate.Session, in *proto.LoginRequest) error {

	//session.Send(&proto.LoginResponse{})
	return nil
}
