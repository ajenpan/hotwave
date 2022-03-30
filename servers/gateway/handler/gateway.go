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
	utlhandle "hotwave/util/handle"
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

	g.gateCallTable = utlhandle.ExtractProtoFile(proto.File_servers_gateway_proto_gate_proto, g)

	g.gateCallTable.Range(func(key string, value *utlhandle.Method) bool {
		logger.Infof("handler gate message: %s", key)
		return true
	})

	for _, v := range []string{
		strings.TrimSuffix(string(protobuf.MessageName(&proto.LoginRequest{})), "Request"),
		strings.TrimSuffix(string(protobuf.MessageName(&proto.EchoRequest{})), "Request"),
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

	gateCallTable *utlhandle.CallTable
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

func sendMsgToUser(s gate.Session, msg protobuf.Message) error {
	if s == nil {
		return fmt.Errorf("session is nil")
	}

	if msg == nil {
		return nil
	}
	body, err := protobuf.Marshal(msg)
	if err != nil {
		return err
	}

	warp := &codec.AsyncMessage{
		Body:    body,
		MsgName: string(protobuf.MessageName(msg)),
	}
	return s.Send(warp)
}

func (g *Gater) OnGateMessage(session gate.Session, msg *codec.AsyncMessage) {
	name := protoreflect.FullName(msg.MsgName)
	serverName := string(name.Parent())

	uid := session.UID()

	if uid == 0 {
		allowChekcer := func(msgName string) bool {
			_, has := g.allowList.Load(msgName)
			return has
		}

		if !allowChekcer(msg.MsgName) {
			g.SendSessionErrorAndClose(session, fmt.Errorf("no permition to send name:%s", msg.MsgName))
			return
		}
	}

	converMsg := &proto.UserMessageWraper{
		MsgName: msg.MsgName,
		Meta:    msg.Meta,
		Body:    msg.Body,
		UserId:  uid,
	}

	//TODO : no special server
	if serverName == "gate" {
		g.OnUserMessage(session, converMsg)
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
			g.router.GetService(serverName)
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

func (g *Gater) OnUserMessage(s gate.Session, msg *proto.UserMessageWraper) {
	logger.Infof("OnUserMessage: %s", msg.MsgName)

	g.gateCallTable.Get(msg.MsgName)

	method := g.gateCallTable.Get(msg.MsgName)
	if method == nil {
		return
	}

	reqestV := reflect.New(method.RequestType)
	reqest := reqestV.Interface().(protobuf.Message)

	if err := protobuf.Unmarshal(msg.Body, reqest); err != nil {
		logger.Error(err)
		return
	}
	ctx := CtxWithSession(context.Background(), s)
	callResult := method.Call(reflect.ValueOf(g), reflect.ValueOf(ctx), reflect.ValueOf(reqest))
	var respErr error
	if len(callResult) == 1 {
		if ierr, ok := callResult[0].Interface().(error); ok {
			respErr = ierr
		}
	} else if len(callResult) == 2 {
		if resp, ok := callResult[0].Interface().(protobuf.Message); ok {
			sendMsgToUser(s, resp)
		}
		if err, ok := callResult[1].Interface().(error); ok {
			respErr = err
		}
	} else {
		logger.Warn("call result is not 1 or 2")
	}

	if respErr != nil {
		logger.Errorf("call method error: %s", respErr)
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

type CtxSessionKey struct{}

func CtxWithSession(ctx context.Context, session gate.Session) context.Context {
	return context.WithValue(ctx, CtxSessionKey{}, session)
}

func CtxWithSessionValue(ctx context.Context) (gate.Session, bool) {
	v, ok := ctx.Value(CtxSessionKey{}).(gate.Session)
	return v, ok
}

func (g *Gater) Login(ctx context.Context, in *proto.LoginRequest) (*proto.LoginResponse, error) {
	out := &proto.LoginResponse{
		Flag: proto.LoginResponse_Success,
	}

	switch c := in.Checker.(type) {
	case *proto.LoginRequest_AccountInfo:
		if strings.HasPrefix(c.AccountInfo.Account, "test") {

		}
	case *proto.LoginRequest_SessionInfo:
		{
		}
	default:
		out.Flag = proto.LoginResponse_UnknowError
		return out, nil
	}

	uid := uint64(1001)

	if s, ok := CtxWithSessionValue(ctx); ok {
		s.SetUID(uid)
		g.user2session.Store(uid, s)
	}

	return out, nil
}

func (g *Gater) Echo(ctx context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	out := &proto.EchoResponse{
		Data: in.Data,
	}

	return out, nil
}
