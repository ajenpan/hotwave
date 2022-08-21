// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.1
// source: event/proto/event.proto

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

// EventClient is the client API for Event service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Event_SubscribeClient, error)
}

type eventClient struct {
	cc grpc.ClientConnInterface
}

func NewEventClient(cc grpc.ClientConnInterface) EventClient {
	return &eventClient{cc}
}

func (c *eventClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (Event_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Event_ServiceDesc.Streams[0], "/event.Event/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Event_SubscribeClient interface {
	Recv() (*EventMessage, error)
	grpc.ClientStream
}

type eventSubscribeClient struct {
	grpc.ClientStream
}

func (x *eventSubscribeClient) Recv() (*EventMessage, error) {
	m := new(EventMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventServer is the server API for Event service.
// All implementations must embed UnimplementedEventServer
// for forward compatibility
type EventServer interface {
	Subscribe(*SubscribeRequest, Event_SubscribeServer) error
	mustEmbedUnimplementedEventServer()
}

// UnimplementedEventServer must be embedded to have forward compatible implementations.
type UnimplementedEventServer struct {
}

func (UnimplementedEventServer) Subscribe(*SubscribeRequest, Event_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedEventServer) mustEmbedUnimplementedEventServer() {}

// UnsafeEventServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServer will
// result in compilation errors.
type UnsafeEventServer interface {
	mustEmbedUnimplementedEventServer()
}

func RegisterEventServer(s grpc.ServiceRegistrar, srv EventServer) {
	s.RegisterService(&Event_ServiceDesc, srv)
}

func _Event_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventServer).Subscribe(m, &eventSubscribeServer{stream})
}

type Event_SubscribeServer interface {
	Send(*EventMessage) error
	grpc.ServerStream
}

type eventSubscribeServer struct {
	grpc.ServerStream
}

func (x *eventSubscribeServer) Send(m *EventMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Event_ServiceDesc is the grpc.ServiceDesc for Event service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Event_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.Event",
	HandlerType: (*EventServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Event_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "event/proto/event.proto",
}
