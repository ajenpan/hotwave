package grpcgate

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	log "hotwave/logger"
	gwProto "hotwave/service/gateway/proto"
	"hotwave/transport"
)

type ClientGate struct {
	transport.SyncMapSocketMeta

	NodeID        string
	NodeName      string
	GrpcConn      *grpc.ClientConn
	GetewayClient gwProto.GatewayClient
	proxyc        gwProto.Gateway_ProxyClient

	status transport.SessionStat

	OnConnStatusFunc  func(c *ClientGate, ss transport.SessionStat)
	OnUserMessageFunc func(s *UserSocket, msg *gwProto.ToServerMessage)
}

func (c *ClientGate) Close() {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) != transport.Connected {
		return
	}
	c.proxyc.CloseSend()
}

func (c *ClientGate) ID() string {
	return c.NodeID
}

func (c *ClientGate) LocalAddr() string {
	return ""
}

func (c *ClientGate) RemoteAddr() string {
	return ""
}

func (c *ClientGate) OnConnStatus(ss transport.SessionStat) {
	atomic.SwapInt32((*int32)(&c.status), int32(ss))
	log.Info("grpc-client-OnConnStatus", ss)
	// switch ss {
	// case transport.Connected:
	// case transport.Disconnected:
	// 	go c.Reconnect()
	// default:
	// }

	if c.OnConnStatusFunc != nil {
		c.OnConnStatusFunc(c, ss)
	}
}

func (s *ClientGate) Reconnect() {
	err := s.Connect()
	if err != nil {
		log.Error(err)
		time.Sleep(2 * time.Second)
		go s.Reconnect()
	}
}

func (c *ClientGate) Connect() error {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) == transport.Connected {
		return nil
	}
	md := metadata.New(map[string]string{"nodeid": c.NodeID, "nodename": c.NodeName})
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

			c.onMessage(in)
		}
		log.Error("connect dist error: ", recvErr)
	}()
	return nil
}

func (c *ClientGate) onMessage(in *gwProto.ToServerMessage) {
	log.Info("gate-client-Recv: ", in.Name)

	if c.OnUserMessageFunc != nil {
		c.OnUserMessageFunc(&UserSocket{RemoteSocketID: in.FromSocketid, UID: in.FromUid, client: c}, in)
	}
}

func (c *ClientGate) sendMessage(warp *gwProto.ToClientMessage) error {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) != transport.Connected {
		return io.EOF
	}
	log.Info("gate-client-Send: ", warp.Name)
	return c.proxyc.Send(warp)
}
