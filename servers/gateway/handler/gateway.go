package handler

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	// "google.golang.org/protobuf/runtime/protoimpl"
	// "google.golang.org/protobuf/types/dynamicpb"

	"google.golang.org/protobuf/reflect/protoreflect"

	frame "hotwave"
	"hotwave/logger"
	"hotwave/registry"
	"hotwave/servers/gateway/gate"
	"hotwave/servers/gateway/gate/codec"
	"hotwave/servers/gateway/proto"
	"hotwave/servers/gateway/router"
)

func NewGater(frame *frame.Frame) *Gater {
	router := router.New(frame.Options().Registry)

	if err := router.Start(); err != nil {
		return nil
	}

	g := &Gater{
		frame:  frame,
		router: router,
	}

	w, err := frame.Options().Registry.Watch()
	if err != nil {
		return nil
	}

	go func() {
		for {
			res, err := w.Next()
			if err != nil {
				logger.Error(err)
				return
			}
			g.onWatch(res)
		}
	}()

	return g
}

type Gater struct {
	proto.UnimplementedGatewayServer

	user2session sync.Map //map[uint64]*gate.Session

	tclientLock sync.RWMutex
	tclient     map[string]proto.GateAdapterClient

	router *router.Router
	frame  *frame.Frame
}

func (g *Gater) onWatch(res *registry.Result) {
	switch res.Action {
	case registry.Create.String():
	case registry.Update.String():
	}

	// logger.Infof("onWatch: %s %s", res.Action, res.Service.Name)
}

func (g *Gater) GetGateAdapterClient(node *router.Node) proto.GateAdapterClient {
	if node == nil {
		return nil
	}

	g.tclientLock.Lock()
	defer g.tclientLock.Unlock()
	if client, ok := g.tclient[node.Id]; ok {
		return client
	}
	conn, err := grpc.Dial(node.Address)
	if err != nil {
		return nil
	}
	client := proto.NewGateAdapterClient(conn)
	g.tclient[node.Id] = client
	return client
}

func (g *Gater) SendSessionErrorAndClose(session gate.Session, err error) {
	logger.Error(err)
	session.Close()
}

func NewConn(node *router.Node) *grpc.ClientConn {
	c, err := grpc.Dial(node.Address)
	if err != nil {
		return nil
	}
	return c
}

func (g *Gater) OnGateMessage(session gate.Session, msg *codec.AsyncMessage) {
	name := protoreflect.FullName(msg.Name)
	serverName := name.Parent()

	uid := session.UID()

	allowChekcer := func(msgName string) bool {
		return true
	}

	if uid == 0 {
		if !allowChekcer(msg.Name) {
			g.SendSessionErrorAndClose(session, fmt.Errorf("no permition to send name:%s", msg.Name))
			return
		}
	}

	converMsg := &proto.UserMessageWraper{
		Name:   msg.Name,
		Meta:   msg.Meta,
		Body:   msg.Body,
		UserId: uid,
	}

	if g.frame.Options().Name == string(serverName) {
		g.OnUserMessage(converMsg)
		return
	}

	var node *router.Node

	if len(msg.Nodeid) != 0 {
		node = g.router.GetNode(msg.Nodeid)
	} else {
		if v, has := session.GetMeta(fmt.Sprintf("server-%s", serverName)); has {
			msg.Nodeid = v.(string)
			node = g.router.GetNode(msg.Nodeid)
		} else {
			// node = g.router.GetServices(msg.Name)
		}
	}
	client := g.GetGateAdapterClient(node)
	if client == nil {
		g.SendSessionErrorAndClose(session, fmt.Errorf("gateway: can't find server:%s", serverName))
		return
	}

	_, err := client.UserMessage(context.Background(), converMsg)
	if err != nil {
		g.SendSessionErrorAndClose(session, fmt.Errorf("route msg to %v error", err))
	}
}

func (g *Gater) OnGateConnStat(session gate.Session, status gate.SocketStat) {
	fmt.Printf("session:%s, connect state:%v \n", session.ID(), status)

	switch status {
	case gate.SocketStatConnected:
	case gate.SocketStatDisconnected:
	}
}

func (g *Gater) OnUserMessage(msg *proto.UserMessageWraper) {
	logger.Infof("OnUserMessage: %s", msg.Name)
}

func (g *Gater) SendMessageToUse(ctx context.Context, in *proto.SendMessageToUserRequest) (*proto.SendMessageToUserResponse, error) {
	out := &proto.SendMessageToUserResponse{}
	v, has := g.user2session.Load(in.Uid)
	if !has {
		return nil, fmt.Errorf("user %d not online", in.Uid)
	}
	session, ok := v.(gate.Session)
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	msg := &codec.AsyncMessage{}
	err := session.Send(msg)
	return out, err
}

func (g *Gater) OnLoginRequest(session gate.Session, in *proto.LoginRequest) error {
	return nil
}
