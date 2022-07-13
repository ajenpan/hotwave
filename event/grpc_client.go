package event

import (
	"context"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	evproto "hotwave/event/proto"
	log "hotwave/logger"
	"hotwave/transport"
)

type GRPCClient struct {
	NodeID string
	Conn   *grpc.ClientConn
	Topics []string
	Recv   Recver

	client evproto.Event_SubscribeClient
	status transport.SessionStat
}

func (c *GRPCClient) Close() error {
	return c.client.CloseSend()
}

func (c *GRPCClient) onConnStatus(ss transport.SessionStat) {
	atomic.SwapInt32((*int32)(&c.status), int32(ss))
	log.Info("grpc-client-OnConnStatus", ss)
	switch ss {
	case transport.Connected:
	case transport.Disconnected:
		go c.Reconnect()
	default:
	}
}

func (c *GRPCClient) Reconnect() {
	err := c.Connect()
	if err != nil {
		log.Error(err)
		time.Sleep(2 * time.Second)
		go c.Reconnect()
	}
}

func (c *GRPCClient) Connect() error {
	if transport.SessionStat(atomic.LoadInt32((*int32)(&c.status))) == transport.Connected {
		return nil
	}

	ctx := context.Background()
	req := &evproto.SubscribeRequest{
		Topics: c.Topics,
	}

	proxyc, err := evproto.NewEventClient(c.Conn).Subscribe(ctx, req)
	if err != nil {
		return err
	}
	c.client = proxyc
	go func() {
		c.onConnStatus(transport.Connected)
		defer c.onConnStatus(transport.Disconnected)
		var recvErr error
		for {
			in, err := proxyc.Recv()
			if err != nil {
				recvErr = err
				break
			}

			if c.Recv != nil {
				c.Recv.OnEvent(in)
			}
		}
		log.Error("connect dist error: ", recvErr)
	}()
	return nil
}
