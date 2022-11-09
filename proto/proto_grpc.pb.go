// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: proto/proto.proto

package ring

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

// RingClient is the client API for Ring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RingClient interface {
	Contact(ctx context.Context, in *Request, opts ...grpc.CallOption) (*EmptyMessage, error)
}

type ringClient struct {
	cc grpc.ClientConnInterface
}

func NewRingClient(cc grpc.ClientConnInterface) RingClient {
	return &ringClient{cc}
}

func (c *ringClient) Contact(ctx context.Context, in *Request, opts ...grpc.CallOption) (*EmptyMessage, error) {
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, "/ring.Ring/contact", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RingServer is the server API for Ring service.
// All implementations must embed UnimplementedRingServer
// for forward compatibility
type RingServer interface {
	Contact(context.Context, *Request) (*EmptyMessage, error)
	mustEmbedUnimplementedRingServer()
}

// UnimplementedRingServer must be embedded to have forward compatible implementations.
type UnimplementedRingServer struct {
}

func (UnimplementedRingServer) Contact(context.Context, *Request) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Contact not implemented")
}
func (UnimplementedRingServer) mustEmbedUnimplementedRingServer() {}

// UnsafeRingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RingServer will
// result in compilation errors.
type UnsafeRingServer interface {
	mustEmbedUnimplementedRingServer()
}

func RegisterRingServer(s grpc.ServiceRegistrar, srv RingServer) {
	s.RegisterService(&Ring_ServiceDesc, srv)
}

func _Ring_Contact_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServer).Contact(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ring.Ring/contact",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServer).Contact(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Ring_ServiceDesc is the grpc.ServiceDesc for Ring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ring.Ring",
	HandlerType: (*RingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "contact",
			Handler:    _Ring_Contact_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proto.proto",
}