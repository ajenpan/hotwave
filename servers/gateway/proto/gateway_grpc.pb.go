// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: servers/gateway/proto/gateway.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GatewayClient is the client API for Gateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayClient interface {
	SendMessageToUse(ctx context.Context, in *SendMessageToUserRequest, opts ...grpc.CallOption) (*SendMessageToUserResponse, error)
}

type gatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayClient(cc grpc.ClientConnInterface) GatewayClient {
	return &gatewayClient{cc}
}

func (c *gatewayClient) SendMessageToUse(ctx context.Context, in *SendMessageToUserRequest, opts ...grpc.CallOption) (*SendMessageToUserResponse, error) {
	out := new(SendMessageToUserResponse)
	err := c.cc.Invoke(ctx, "/gateway.Gateway/SendMessageToUse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServer is the server API for Gateway service.
// All implementations must embed UnimplementedGatewayServer
// for forward compatibility
type GatewayServer interface {
	SendMessageToUse(context.Context, *SendMessageToUserRequest) (*SendMessageToUserResponse, error)
	mustEmbedUnimplementedGatewayServer()
}

// UnimplementedGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedGatewayServer struct {
}

func (UnimplementedGatewayServer) SendMessageToUse(context.Context, *SendMessageToUserRequest) (*SendMessageToUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessageToUse not implemented")
}
func (UnimplementedGatewayServer) mustEmbedUnimplementedGatewayServer() {}

// UnsafeGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServer will
// result in compilation errors.
type UnsafeGatewayServer interface {
	mustEmbedUnimplementedGatewayServer()
}

func RegisterGatewayServer(s grpc.ServiceRegistrar, srv GatewayServer) {
	s.RegisterService(&Gateway_ServiceDesc, srv)
}

func _Gateway_SendMessageToUse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageToUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).SendMessageToUse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.Gateway/SendMessageToUse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).SendMessageToUse(ctx, req.(*SendMessageToUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Gateway_ServiceDesc is the grpc.ServiceDesc for Gateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.Gateway",
	HandlerType: (*GatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessageToUse",
			Handler:    _Gateway_SendMessageToUse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "servers/gateway/proto/gateway.proto",
}

// GateAdpaterClient is the client API for GateAdpater service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GateAdpaterClient interface {
	UserMessage(ctx context.Context, opts ...grpc.CallOption) (GateAdpater_UserMessageClient, error)
}

type gateAdpaterClient struct {
	cc grpc.ClientConnInterface
}

func NewGateAdpaterClient(cc grpc.ClientConnInterface) GateAdpaterClient {
	return &gateAdpaterClient{cc}
}

func (c *gateAdpaterClient) UserMessage(ctx context.Context, opts ...grpc.CallOption) (GateAdpater_UserMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &GateAdpater_ServiceDesc.Streams[0], "/gateway.GateAdpater/UserMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &gateAdpaterUserMessageClient{stream}
	return x, nil
}

type GateAdpater_UserMessageClient interface {
	Send(*UserMessageWraper) error
	CloseAndRecv() (*SteamClosed, error)
	grpc.ClientStream
}

type gateAdpaterUserMessageClient struct {
	grpc.ClientStream
}

func (x *gateAdpaterUserMessageClient) Send(m *UserMessageWraper) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gateAdpaterUserMessageClient) CloseAndRecv() (*SteamClosed, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SteamClosed)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GateAdpaterServer is the server API for GateAdpater service.
// All implementations must embed UnimplementedGateAdpaterServer
// for forward compatibility
type GateAdpaterServer interface {
	UserMessage(GateAdpater_UserMessageServer) error
	mustEmbedUnimplementedGateAdpaterServer()
}

// UnimplementedGateAdpaterServer must be embedded to have forward compatible implementations.
type UnimplementedGateAdpaterServer struct {
}

func (UnimplementedGateAdpaterServer) UserMessage(GateAdpater_UserMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method UserMessage not implemented")
}
func (UnimplementedGateAdpaterServer) mustEmbedUnimplementedGateAdpaterServer() {}

// UnsafeGateAdpaterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GateAdpaterServer will
// result in compilation errors.
type UnsafeGateAdpaterServer interface {
	mustEmbedUnimplementedGateAdpaterServer()
}

func RegisterGateAdpaterServer(s grpc.ServiceRegistrar, srv GateAdpaterServer) {
	s.RegisterService(&GateAdpater_ServiceDesc, srv)
}

func _GateAdpater_UserMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GateAdpaterServer).UserMessage(&gateAdpaterUserMessageServer{stream})
}

type GateAdpater_UserMessageServer interface {
	SendAndClose(*SteamClosed) error
	Recv() (*UserMessageWraper, error)
	grpc.ServerStream
}

type gateAdpaterUserMessageServer struct {
	grpc.ServerStream
}

func (x *gateAdpaterUserMessageServer) SendAndClose(m *SteamClosed) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gateAdpaterUserMessageServer) Recv() (*UserMessageWraper, error) {
	m := new(UserMessageWraper)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GateAdpater_ServiceDesc is the grpc.ServiceDesc for GateAdpater service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GateAdpater_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.GateAdpater",
	HandlerType: (*GateAdpaterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UserMessage",
			Handler:       _GateAdpater_UserMessage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "servers/gateway/proto/gateway.proto",
}
