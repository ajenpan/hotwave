package transfer

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	protocol "hotwave/servers/gateway/proto"
)

type Adpater interface {
	OnUserMessage(msg *protocol.UserMessageWraper)
}

func NewTransfer() *Transfer {
	ret := &Transfer{
		exit: make(chan chan error),
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterGateAdapterServer(grpcServer, ret)
	ret.grpcServer = grpcServer

	return ret
}

type Transfer struct {
	protocol.UnimplementedGateAdapterServer

	sync.RWMutex
	grpcServer *grpc.Server
	exit       chan chan error
	adapter    Adpater
}

func (t *Transfer) UserMessage(ctx context.Context, in *protocol.UserMessageWraper) (*protocol.SteamClosed, error) {
	t.adapter.OnUserMessage(in)
	return nil, nil
}
