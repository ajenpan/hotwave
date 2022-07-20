package client

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	log "hotwave/logger"
	gwProto "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type UserSession struct {
	client   *GRPCClient
	UID      int64
	SocketID string
}

func (s *UserSession) Send(msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	warp := &gwProto.ToUserMessage{
		ToUid:      s.UID,
		ToSocketid: s.SocketID,
		Name:       string(proto.MessageName(msg)),
		Data:       data,
	}

	return s.client.SendMessage(warp)
}

type GRPCClient struct {
	transport.SyncMapSocketMeta

	NodeID        string
	NodeType      string
	GrpcConn      *grpc.ClientConn
	GetewayClient gwProto.GatewayClient
	proxyc        gwProto.Gateway_ProxyClient

	status transport.SessionStat

	OnConnStatusFunc  func(c *GRPCClient, ss transport.SessionStat)
	OnUserMessageFunc func(s *UserSession, msg *gwProto.ToServerMessage)
}

func (c *GRPCClient) Close() {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) != transport.Connected {
		return
	}
	c.proxyc.CloseSend()
}

func (c *GRPCClient) ID() string {
	return c.NodeID
}
func (c *GRPCClient) LocalAddr() string {
	return ""
}
func (c *GRPCClient) RemoteAddr() string {
	return ""
}

func (c *GRPCClient) OnConnStatus(ss transport.SessionStat) {
	atomic.SwapInt32((*int32)(&c.status), int32(ss))
	log.Info("grpc-client-OnConnStatus", ss)

	if c.OnConnStatusFunc != nil {
		c.OnConnStatusFunc(c, ss)
	}
}

func (s *GRPCClient) Reconnect() {
	err := s.Connect()
	if err != nil {
		log.Error(err)
		time.Sleep(2 * time.Second)
		go s.Reconnect()
	}
}

func (c *GRPCClient) Connect() error {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) == transport.Connected {
		return nil
	}
	md := metadata.New(map[string]string{"nodeid": c.NodeID, "nodename": c.NodeType})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	c.GetewayClient = gwProto.NewGatewayClient(c.GrpcConn)
	proxyc, err := c.GetewayClient.Proxy(ctx)
	if err != nil {
		return err
	}
	c.proxyc = proxyc
	go func() {
		c.OnConnStatus(transport.Connected)
		defer c.OnConnStatus(transport.Disconnected)
		var recvErr error
		for {
			in, err := proxyc.Recv()
			if err != nil {
				recvErr = err
				break
			}
			log.Info("gate-client-Recv: ", in.Name)
			if c.OnUserMessageFunc != nil {
				c.OnUserMessageFunc(&UserSession{SocketID: in.FromSocketid, UID: in.FromUid, client: c}, in)
			}
		}
		log.Error("connect dist error: ", recvErr)
	}()
	return nil
}

func (c *GRPCClient) SendMessage(warp *gwProto.ToUserMessage) error {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) != transport.Connected {
		return io.EOF
	}
	log.Info("gate-client-Send: ", warp.Name)
	return c.proxyc.Send(warp)
}
