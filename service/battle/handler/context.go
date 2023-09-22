package handler

import (
	"context"

	"hotwave/transport/tcp"
)

type userInfoKey struct{}
type tcpSocketKey struct{}

var UserInfoKey = &userInfoKey{}
var TcpSocketKey = &tcpSocketKey{}

func GetUserInfo(ctx context.Context) *tcp.UserInfo {
	return ctx.Value(UserInfoKey).(*tcp.UserInfo)
}

func WithUserInfo(ctx context.Context, uinfo *tcp.UserInfo) context.Context {
	return context.WithValue(ctx, UserInfoKey, uinfo)
}

func GetTcpSocket(ctx context.Context) *tcp.Socket {
	return ctx.Value(TcpSocketKey).(*tcp.Socket)
}

func WithTcpSocket(ctx context.Context, s *tcp.Socket) context.Context {
	return context.WithValue(ctx, TcpSocketKey, s)
}
