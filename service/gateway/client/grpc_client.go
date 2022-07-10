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
	"hotwave/session"
)

type GRPCClient struct {
	NodeID   string
	NodeName string
	Conn     *grpc.ClientConn
	proxyc   gwProto.Gateway_ProxyClient

	status session.SessionStat

	OnMessage func(*GRPCClient, *gwProto.ToServerMessage)
}

func (c *GRPCClient) Close() error {
	c.proxyc.CloseSend()
	return nil
}

func (c *GRPCClient) OnConnStatus(ss session.SessionStat) {
	atomic.SwapInt32((*int32)(&c.status), int32(ss))
	log.Info("grpc-client-OnConnStatus", ss)
	switch ss {
	case session.SessionStatConnected:
	case session.SessionStatDisconnected:
		go reconnect(c)
	default:
	}
}

func reconnect(s *GRPCClient) {
	time.Sleep(2 * time.Second)
	err := s.Connect()
	if err != nil {
		log.Error(err)
		go reconnect(s)
	}
}

func (c *GRPCClient) Connect() error {
	if session.SessionStat(atomic.LoadInt32((*int32)(&c.status))) == session.SessionStatConnected {
		return nil
	}

	md := metadata.New(map[string]string{"nodeid": c.NodeID, "nodename": c.NodeName})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	proxyc, err := gwProto.NewGatewayClient(c.Conn).Proxy(ctx)
	if err != nil {
		return err
	}
	c.proxyc = proxyc
	go func() {
		c.OnConnStatus(session.SessionStatConnected)
		defer c.OnConnStatus(session.SessionStatDisconnected)
		var recvErr error
		for {
			in, err := proxyc.Recv()
			if err != nil {
				recvErr = err
				break
			}
			if c.OnMessage != nil {
				c.OnMessage(c, in)
			}
		}
		log.Error("connect dist error: ", recvErr)
	}()
	return nil
}

func (c *GRPCClient) Send(uid int64, msg proto.Message) error {
	if session.SessionStat(atomic.LoadInt32((*int32)(&c.status))) != session.SessionStatConnected {
		return io.EOF
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	warp := gwProto.ToClientMessage{
		Uid:  uid,
		Name: string(proto.MessageName(msg)),
		Data: data,
	}
	return c.proxyc.Send(&warp)
}
