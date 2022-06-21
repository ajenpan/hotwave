package handler

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	frame "hotwave"
	"hotwave/logger"
	"hotwave/registry"
	"hotwave/services/gateway/gate"
	protocol "hotwave/services/gateway/proto"
	"hotwave/services/gateway/protostore"
	"hotwave/services/gateway/router"
	utlhandle "hotwave/util/handle"
)

func NewGater(frame *frame.Frame) *Gater {
	router := router.New(frame.Options().Registry)

	if err := router.Start(); err != nil {
		return nil
	}

	g := &Gater{
		frame:      frame,
		router:     router,
		protoStore: protostore.NewMomoryStore(),
	}

	g.gateCallTable = utlhandle.ExtractProtoFile(protocol.File_gate_proto, g)

	g.gateCallTable.Range(func(key string, value *utlhandle.Method) bool {
		logger.Infof("handler gate message: %s", key)
		return true
	})

	for _, v := range []string{
		strings.TrimSuffix(string(protobuf.MessageName(&protocol.LoginRequest{})), "Request"),
		strings.TrimSuffix(string(protobuf.MessageName(&protocol.EchoRequest{})), "Request"),
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
	// protocol.UnimplementedGatewayServer
	// proto.UnimplementedGateServer
	// proto.UnimplementedGateAdapterServer

	user2session sync.Map //map[uint64]*gate.Session

	tclientLock sync.RWMutex
	// tclient     map[string]protocol.GateAdapterClient

	router *router.Router
	frame  *frame.Frame

	allowList sync.Map

	gateCallTable *utlhandle.CallTable

	protoStore *protostore.MomoryStore
}

func (g *Gater) onNodeWatch(res *registry.Result) {
	// switch res.Action {
	// case registry.Create.String():
	// 	fallthrough
	// case registry.Update.String():
	// 	if res.Service.Name == "gateway" {
	// 		return
	// 	}
	// 	if res.Service == nil || len(res.Service.Nodes) == 0 {
	// 		return
	// 	}
	// 	node := res.Service.Nodes[0]
	// 	if len(node.Id) <= 2 {
	// 		return
	// 	}

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 	defer cancel()

	// 	conn, err := grpc.DialContext(ctx, node.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return
	// 	}

	// 	client := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	// 	steam, err := client.ServerReflectionInfo(context.Background())
	// 	if err != nil {
	// 		logger.Error(err)
	// 		return
	// 	}
	// 	defer steam.CloseSend()

	// 	if err := steam.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
	// 		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
	// 			FileContainingSymbol: res.Service.Name,
	// 		},
	// 	}); err != nil {
	// 		logger.Error(err)
	// 		return
	// 	}

	// 	resp, err := steam.Recv()
	// 	if err != nil {
	// 		logger.Error(err)
	// 	}

	// 	fmt.Println(resp)
	// }
	logger.Infof("onWatch: %s %s", res.Action, res.Service.Name)
}

func (g *Gater) GetGrpcConn(node *router.Node) *grpc.ClientConn {
	conn, err := grpc.Dial(node.Address)
	if err != nil {
		return nil
	}
	return conn
}

// func (g *Gater) GetGateAdapterClient(node *router.Node) protocol.GateAdapterClient {
// 	if node == nil {
// 		return nil
// 	}

// 	g.tclientLock.RLock()
// 	client, ok := g.tclient[node.Id]
// 	if ok {
// 		g.tclientLock.RUnlock()
// 		return client
// 	}
// 	g.tclientLock.RUnlock()

// 	conn, err := grpc.Dial(node.Address)
// 	if err != nil {
// 		return nil
// 	}
// 	client = protocol.NewGateAdapterClient(conn)

// 	g.tclientLock.Lock()
// 	defer g.tclientLock.Unlock()

// 	g.tclient[node.Id] = client
// 	return client
// }

func (g *Gater) SendSessionErrorAndClose(session gate.Session, err error) {
	logger.Error(err)
	session.Close()
}

func (g *Gater) GetNode(nodeId string, servername string) *router.Node {
	var ret *router.Node
	if nodeId != "" {
		ret = g.router.GetNode(nodeId)
	}
	if ret == nil {
		svrs, err := g.router.GetService(servername)
		if err != nil {
			logger.Error(err)
			return nil
		}
		//select ret
		if len(svrs) == 0 {
			logger.Error("GetService name %s, is empty", servername)
			return nil
		}
		// simple random
		// return svr.Nodes[0]
	}
	return ret
}

func (g *Gater) OnGateMethod(ctx context.Context, msg *protocol.ClientMessageWraper) (*protocol.ClientMessageWraper, error) {

	// msg.Method.
	// req,resp,err:=g.protoStore.NewTypeByMethod(msg.Method)
	// if err!=nil{
	// 	return nil,err
	// }

	// c:= &grpc.ClientConn{}
	// c.Invoke()
	// err=protobuf.Unmarshal(msg.Data,req)

	// var grr error

	// var dialCtx context.Context
	// var cancel context.CancelFunc
	// dialCtx, cancel = context.WithCancel(ctx)
	// defer cancel()

	// grpcDialOptions := []grpc.DialOption{
	// 	grpc.WithDefaultCallOptions(
	// 	),
	// }

	// grr := grpc.DialContext(dialCtx,"localhost",	grpc.ForceCodec(cf), )

	// cc, err := g.pool.getConn(dialCtx, address, grpcDialOptions...)
	// if err != nil {
	// 	return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error sending request: %v", err))
	// }
	// defer func() {
	// 	g.pool.release(address, cc, grr)
	// }()

	// grpcCallOptions := []grpc.CallOption{
	// 	grpc.ForceCodec(cf),
	// 	grpc.CallContentSubtype(cf.Name())}
	// if opts := g.getGrpcCallOptions(); opts != nil {
	// 	grpcCallOptions = append(grpcCallOptions, opts...)
	// }
	// err := cc.Invoke(ctx, methodToGRPC(req.Service(), req.Endpoint()), req.Body(), rsp, grpcCallOptions...)

	return nil, nil
}

func (g *Gater) OnGateAsync(session gate.Session, msg *protocol.ClientMessageWraper) {
	name := protoreflect.FullName(msg.Method)
	serverName := string(name.Parent())
	fmt.Println(serverName)
	uid := session.UID()

	if uid == 0 {
		allowChekcer := func(msgName string) bool {
			_, has := g.allowList.Load(msgName)
			return has
		}

		if !allowChekcer(msg.Method) {
			g.SendSessionErrorAndClose(session, fmt.Errorf("no permition to send name:%s", msg.Method))
			return
		}
	}

	// converMsg := &proto.ClientMessageWraper{
	// 	MsgName: msg.MsgName,
	// 	Meta:    msg.Meta,
	// 	Body:    msg.Body,
	// 	UserId:  uid,
	// }

	// //TODO : no special server
	// if serverName == "gate" {
	// 	g.OnUserMessage(session, converMsg)
	// 	return
	// }

	// var node *router.Node

	// if len(msg.Nodeid) != 0 {
	// 	node = g.router.GetNode(msg.Nodeid)
	// } else {
	// 	if v, has := session.GetMeta(fmt.Sprintf("server-%s", serverName)); has {
	// 		msg.Nodeid = v.(string)
	// 		node = g.router.GetNode(msg.Nodeid)
	// 	} else {
	// 		g.router.GetService(serverName)
	// 	}
	// }
	// client := g.GetGateAdapterClient(node)
	// if client == nil {
	// 	g.SendSessionErrorAndClose(session, fmt.Errorf("gateway: can't find server:%s", serverName))
	// 	return
	// }

	// _, err := client.UserEvent(context.Background(), converMsg)
	// if err != nil {
	// 	g.SendSessionErrorAndClose(session, fmt.Errorf("route msg to %v error", err))
	// }
}

func (g *Gater) OnGateConnStat(session gate.Session, status gate.SocketStat) {
	fmt.Printf("session:%s, connect state:%v \n", session.ID(), status)
	switch status {
	case gate.SocketStatConnected:
	case gate.SocketStatDisconnected:
	}
}

func (g *Gater) OnUserMessage(s gate.Session, msg *protocol.ClientMessageWraper) {
	logger.Infof("OnUserMessage: %s", msg.Method)

	g.gateCallTable.Get(msg.Method)

	method := g.gateCallTable.Get(msg.Method)
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
			s.Send(resp)
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

func (g *Gater) SendMessageToUse(ctx context.Context, in *protocol.SendMessageToUserRequest) (*protocol.SendMessageToUserResponse, error) {
	out := &protocol.SendMessageToUserResponse{}
	v, has := g.user2session.Load(in.Uid)
	if !has {
		return nil, fmt.Errorf("user %d not online", in.Uid)
	}
	session, ok := v.(gate.Session)
	if !ok {
		return nil, fmt.Errorf("invalid session")
	}

	msg := &protocol.ClientMessageWraper{}
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

func (g *Gater) Login(ctx context.Context, in *protocol.LoginRequest) (*protocol.LoginResponse, error) {
	out := &protocol.LoginResponse{
		Flag: protocol.LoginResponse_Success,
	}

	switch c := in.Checker.(type) {
	case *protocol.LoginRequest_AccountInfo:
		if strings.HasPrefix(c.AccountInfo.Account, "test") {

		}
	case *protocol.LoginRequest_SessionInfo:
		{
		}
	default:
		out.Flag = protocol.LoginResponse_UnknowError
		return out, nil
	}

	uid := uint64(1001)

	if s, ok := CtxWithSessionValue(ctx); ok {
		s.SetUID(uid)
		g.user2session.Store(uid, s)
	}

	return out, nil
}

func (g *Gater) Echo(ctx context.Context, in *protocol.EchoRequest) (*protocol.EchoResponse, error) {
	out := &protocol.EchoResponse{
		Data: in.Data,
	}
	return out, nil
}

func (g *Gater) ProxyServer(ctx context.Context, in *protocol.ProxyServerRequest) (*protocol.ProxyServerResponse, error) {
	node := g.router.GetNode(in.Nodeid)
	if node == nil {
		return nil, fmt.Errorf("node %s not found", in.Nodeid)
	}
	if len(in.ServerName) == 0 || len(in.Version) == 0 {
		return nil, fmt.Errorf("server name or version is empty")
	}
	key := fmt.Sprintf("%s-%s", in.ServerName, in.Version)
	err := g.protoStore.StoreProtoFiles(key, in.Files)
	if err != nil {
		return nil, err
	}
	return &protocol.ProxyServerResponse{}, nil
}
