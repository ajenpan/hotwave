// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: frame/proto/frame.proto

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

// NodeBaseClient is the client API for NodeBase service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeBaseClient interface {
	UserMessage(ctx context.Context, opts ...grpc.CallOption) (NodeBase_UserMessageClient, error)
	EventMessage(ctx context.Context, opts ...grpc.CallOption) (NodeBase_EventMessageClient, error)
}

type nodeBaseClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeBaseClient(cc grpc.ClientConnInterface) NodeBaseClient {
	return &nodeBaseClient{cc}
}

func (c *nodeBaseClient) UserMessage(ctx context.Context, opts ...grpc.CallOption) (NodeBase_UserMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeBase_ServiceDesc.Streams[0], "/frame.NodeBase/UserMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeBaseUserMessageClient{stream}
	return x, nil
}

type NodeBase_UserMessageClient interface {
	Send(*UserMessageWraper) error
	CloseAndRecv() (*SteamClosed, error)
	grpc.ClientStream
}

type nodeBaseUserMessageClient struct {
	grpc.ClientStream
}

func (x *nodeBaseUserMessageClient) Send(m *UserMessageWraper) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeBaseUserMessageClient) CloseAndRecv() (*SteamClosed, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SteamClosed)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeBaseClient) EventMessage(ctx context.Context, opts ...grpc.CallOption) (NodeBase_EventMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &NodeBase_ServiceDesc.Streams[1], "/frame.NodeBase/EventMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeBaseEventMessageClient{stream}
	return x, nil
}

type NodeBase_EventMessageClient interface {
	Send(*EventMessageWraper) error
	CloseAndRecv() (*SteamClosed, error)
	grpc.ClientStream
}

type nodeBaseEventMessageClient struct {
	grpc.ClientStream
}

func (x *nodeBaseEventMessageClient) Send(m *EventMessageWraper) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeBaseEventMessageClient) CloseAndRecv() (*SteamClosed, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SteamClosed)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeBaseServer is the server API for NodeBase service.
// All implementations must embed UnimplementedNodeBaseServer
// for forward compatibility
type NodeBaseServer interface {
	UserMessage(NodeBase_UserMessageServer) error
	EventMessage(NodeBase_EventMessageServer) error
	mustEmbedUnimplementedNodeBaseServer()
}

// UnimplementedNodeBaseServer must be embedded to have forward compatible implementations.
type UnimplementedNodeBaseServer struct {
}

func (UnimplementedNodeBaseServer) UserMessage(NodeBase_UserMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method UserMessage not implemented")
}
func (UnimplementedNodeBaseServer) EventMessage(NodeBase_EventMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method EventMessage not implemented")
}
func (UnimplementedNodeBaseServer) mustEmbedUnimplementedNodeBaseServer() {}

// UnsafeNodeBaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeBaseServer will
// result in compilation errors.
type UnsafeNodeBaseServer interface {
	mustEmbedUnimplementedNodeBaseServer()
}

func RegisterNodeBaseServer(s grpc.ServiceRegistrar, srv NodeBaseServer) {
	s.RegisterService(&NodeBase_ServiceDesc, srv)
}

func _NodeBase_UserMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeBaseServer).UserMessage(&nodeBaseUserMessageServer{stream})
}

type NodeBase_UserMessageServer interface {
	SendAndClose(*SteamClosed) error
	Recv() (*UserMessageWraper, error)
	grpc.ServerStream
}

type nodeBaseUserMessageServer struct {
	grpc.ServerStream
}

func (x *nodeBaseUserMessageServer) SendAndClose(m *SteamClosed) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeBaseUserMessageServer) Recv() (*UserMessageWraper, error) {
	m := new(UserMessageWraper)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _NodeBase_EventMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeBaseServer).EventMessage(&nodeBaseEventMessageServer{stream})
}

type NodeBase_EventMessageServer interface {
	SendAndClose(*SteamClosed) error
	Recv() (*EventMessageWraper, error)
	grpc.ServerStream
}

type nodeBaseEventMessageServer struct {
	grpc.ServerStream
}

func (x *nodeBaseEventMessageServer) SendAndClose(m *SteamClosed) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeBaseEventMessageServer) Recv() (*EventMessageWraper, error) {
	m := new(EventMessageWraper)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NodeBase_ServiceDesc is the grpc.ServiceDesc for NodeBase service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeBase_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "frame.NodeBase",
	HandlerType: (*NodeBaseServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UserMessage",
			Handler:       _NodeBase_UserMessage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "EventMessage",
			Handler:       _NodeBase_EventMessage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "frame/proto/frame.proto",
}