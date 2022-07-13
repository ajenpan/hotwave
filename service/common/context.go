package common

import (
	"context"

	"hotwave/transport"
)

type socketKey struct{}

func CtxWithSocket(ctx context.Context, s transport.Session) context.Context {
	return context.WithValue(ctx, socketKey{}, s)
}

func GetSocket(ctx context.Context) (transport.Session, bool) {
	s, ok := ctx.Value(socketKey{}).(transport.Session)
	return s, ok
}
