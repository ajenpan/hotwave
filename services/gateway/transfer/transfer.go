package transfer

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	protocol "hotwave/services/gateway/proto"
)

type Adpater interface {
	// OnUserMessage(msg *protocol.UserMessageWraper)
	OnUserAsyncMessage(msg *protocol.AsyncMessageWraper)
}

func NewTransfer() *Transfer {
	ret := &Transfer{
		exit: make(chan chan error),
	}

	grpcServer := grpc.NewServer()
	// protocol.RegisterGateAdapterServer(grpcServer, ret)
	ret.grpcServer = grpcServer

	return ret
}

type Transfer struct {
	// protocol.UnimplementedGateAdapterServer

	sync.RWMutex
	grpcServer *grpc.Server
	exit       chan chan error
	adapter    Adpater
}

func (t *Transfer) OnUserAsyncMessage(ctx context.Context, in *protocol.AsyncMessageWraper) (*protocol.SteamClosed, error) {
	t.adapter.OnUserAsyncMessage(in)
	return nil, nil
}
