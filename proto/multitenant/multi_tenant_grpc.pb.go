// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/multi_tenant.proto

package multitenant

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

// MultiTenantClient is the client API for MultiTenant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultiTenantClient interface {
	CreateTransaction(ctx context.Context, in *CrtTransaction, opts ...grpc.CallOption) (*TransactionResponse, error)
}

type multiTenantClient struct {
	cc grpc.ClientConnInterface
}

func NewMultiTenantClient(cc grpc.ClientConnInterface) MultiTenantClient {
	return &multiTenantClient{cc}
}

func (c *multiTenantClient) CreateTransaction(ctx context.Context, in *CrtTransaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/multi_tenant.MultiTenant/CreateTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MultiTenantServer is the server API for MultiTenant service.
// All implementations must embed UnimplementedMultiTenantServer
// for forward compatibility
type MultiTenantServer interface {
	CreateTransaction(context.Context, *CrtTransaction) (*TransactionResponse, error)
	mustEmbedUnimplementedMultiTenantServer()
}

// UnimplementedMultiTenantServer must be embedded to have forward compatible implementations.
type UnimplementedMultiTenantServer struct {
}

func (UnimplementedMultiTenantServer) CreateTransaction(context.Context, *CrtTransaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransaction not implemented")
}
func (UnimplementedMultiTenantServer) mustEmbedUnimplementedMultiTenantServer() {}

// UnsafeMultiTenantServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultiTenantServer will
// result in compilation errors.
type UnsafeMultiTenantServer interface {
	mustEmbedUnimplementedMultiTenantServer()
}

func RegisterMultiTenantServer(s grpc.ServiceRegistrar, srv MultiTenantServer) {
	s.RegisterService(&MultiTenant_ServiceDesc, srv)
}

func _MultiTenant_CreateTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CrtTransaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiTenantServer).CreateTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/multi_tenant.MultiTenant/CreateTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiTenantServer).CreateTransaction(ctx, req.(*CrtTransaction))
	}
	return interceptor(ctx, in, info, handler)
}

// MultiTenant_ServiceDesc is the grpc.ServiceDesc for MultiTenant service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultiTenant_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "multi_tenant.MultiTenant",
	HandlerType: (*MultiTenantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTransaction",
			Handler:    _MultiTenant_CreateTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/multi_tenant.proto",
}
