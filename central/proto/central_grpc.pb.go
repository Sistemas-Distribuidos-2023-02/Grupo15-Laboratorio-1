// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/central.proto

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

// CentralClient is the client API for Central service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CentralClient interface {
	Keygen(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type centralClient struct {
	cc grpc.ClientConnInterface
}

func NewCentralClient(cc grpc.ClientConnInterface) CentralClient {
	return &centralClient{cc}
}

func (c *centralClient) Keygen(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/central.Central/Keygen", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CentralServer is the server API for Central service.
// All implementations must embed UnimplementedCentralServer
// for forward compatibility
type CentralServer interface {
	Keygen(context.Context, *Message) (*Message, error)
	mustEmbedUnimplementedCentralServer()
}

// UnimplementedCentralServer must be embedded to have forward compatible implementations.
type UnimplementedCentralServer struct {
}

func (UnimplementedCentralServer) Keygen(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keygen not implemented")
}
func (UnimplementedCentralServer) mustEmbedUnimplementedCentralServer() {}

// UnsafeCentralServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CentralServer will
// result in compilation errors.
type UnsafeCentralServer interface {
	mustEmbedUnimplementedCentralServer()
}

func RegisterCentralServer(s grpc.ServiceRegistrar, srv CentralServer) {
	s.RegisterService(&Central_ServiceDesc, srv)
}

func _Central_Keygen_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralServer).Keygen(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/central.Central/Keygen",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralServer).Keygen(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// Central_ServiceDesc is the grpc.ServiceDesc for Central service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Central_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "central.Central",
	HandlerType: (*CentralServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Keygen",
			Handler:    _Central_Keygen_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/central.proto",
}