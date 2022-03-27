package handler

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/grpc"

	// "google.golang.org/protobuf/runtime/protoimpl"
	// "google.golang.org/protobuf/types/dynamicpb"
	protobuf "google.golang.org/protobuf/proto"
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

	g.gateCallTable = ExtractProtoFile(proto.File_servers_gateway_proto_gate_proto, g)

	g.gateCallTable.Range(func(key string, value *Method) bool {
		logger.Infof("handler gate message: %s", key)
		return true
	})

	for _, v := range []string{
		string(protobuf.MessageName(&proto.LoginRequest{})),
		string(protobuf.MessageName(&proto.EchoMessage{})),
	} {
		g.allowList.Store(v, true)
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
			g.onNodeWatch(res)
		}
	}()

	return g
}

type Gater struct {
	proto.UnimplementedGatewayServer
	// proto.UnimplementedGateServer
	// proto.UnimplementedGateAdapterServer

	user2session sync.Map //map[uint64]*gate.Session

	tclientLock sync.RWMutex
	tclient     map[string]proto.GateAdapterClient

	router *router.Router
	frame  *frame.Frame

	allowList sync.Map

	gateCallTable *CallTable
}

func (g *Gater) onNodeWatch(res *registry.Result) {
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

	g.tclientLock.RLock()
	client, ok := g.tclient[node.Id]
	if ok {
		g.tclientLock.RUnlock()
		return client
	}
	g.tclientLock.RUnlock()

	conn, err := grpc.Dial(node.Address)
	if err != nil {
		return nil
	}
	client = proto.NewGateAdapterClient(conn)

	g.tclientLock.Lock()
	defer g.tclientLock.Unlock()

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

	if uid == 0 {
		allowChekcer := func(msgName string) bool {
			_, has := g.allowList.Load(msgName)
			return has
		}

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

	if string(serverName) == "gate" {
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

	g.gateCallTable.Get(msg.Name)

	s, has := g.user2session.Load(msg.UserId)
	if !has {
		return
	}
	s = s.(gate.Session)

	method := g.gateCallTable.Get(msg.Name)
	if method == nil {
		return
	}

	reqestV := reflect.New(method.Req)

	reqest := reqestV.Interface().(protobuf.Message)
	err := protobuf.Unmarshal(msg.Body, reqest)
	if err != nil {
		logger.Error(err)
		return
	}

	callResult := method.Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(reqest)})

	if len(callResult) > 0 && !callResult[0].IsNil() {
		err, ok := callResult[0].Interface().(error)
		if ok {
			if err != nil {
				logger.Error(err)
			}
		}
	}

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

func (g *Gater) Login(session gate.Session, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	out := &proto.LoginResponse{
		Flag: proto.LoginResponse_Success,
	}

	switch c := in.Checker.(type) {
	case *proto.LoginRequest_AccountInfo:
		if strings.HasPrefix(c.AccountInfo.Account, "test") {
			return out, nil
		}
	case *proto.LoginRequest_SessionInfo:
		{
			return out, nil
		}
	}

	out.Flag = proto.LoginResponse_UnknowError
	return out, nil
}

func (g *Gater) Echo(session gate.Session, in *proto.EchoMessage) (*proto.EchoMessage, error) {
	return in, nil
}
