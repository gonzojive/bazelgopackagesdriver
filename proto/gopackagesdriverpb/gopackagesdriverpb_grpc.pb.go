// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gopackagesdriverpb

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

// GoPackagesDriverServiceClient is the client API for GoPackagesDriverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoPackagesDriverServiceClient interface {
	LoadPackages(ctx context.Context, in *LoadPackagesRequest, opts ...grpc.CallOption) (*LoadPackagesResponse, error)
}

type goPackagesDriverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGoPackagesDriverServiceClient(cc grpc.ClientConnInterface) GoPackagesDriverServiceClient {
	return &goPackagesDriverServiceClient{cc}
}

func (c *goPackagesDriverServiceClient) LoadPackages(ctx context.Context, in *LoadPackagesRequest, opts ...grpc.CallOption) (*LoadPackagesResponse, error) {
	out := new(LoadPackagesResponse)
	err := c.cc.Invoke(ctx, "/gopackagesdriver.GoPackagesDriverService/LoadPackages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoPackagesDriverServiceServer is the server API for GoPackagesDriverService service.
// All implementations must embed UnimplementedGoPackagesDriverServiceServer
// for forward compatibility
type GoPackagesDriverServiceServer interface {
	LoadPackages(context.Context, *LoadPackagesRequest) (*LoadPackagesResponse, error)
	mustEmbedUnimplementedGoPackagesDriverServiceServer()
}

// UnimplementedGoPackagesDriverServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGoPackagesDriverServiceServer struct {
}

func (UnimplementedGoPackagesDriverServiceServer) LoadPackages(context.Context, *LoadPackagesRequest) (*LoadPackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPackages not implemented")
}
func (UnimplementedGoPackagesDriverServiceServer) mustEmbedUnimplementedGoPackagesDriverServiceServer() {
}

// UnsafeGoPackagesDriverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoPackagesDriverServiceServer will
// result in compilation errors.
type UnsafeGoPackagesDriverServiceServer interface {
	mustEmbedUnimplementedGoPackagesDriverServiceServer()
}

func RegisterGoPackagesDriverServiceServer(s grpc.ServiceRegistrar, srv GoPackagesDriverServiceServer) {
	s.RegisterService(&GoPackagesDriverService_ServiceDesc, srv)
}

func _GoPackagesDriverService_LoadPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadPackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoPackagesDriverServiceServer).LoadPackages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gopackagesdriver.GoPackagesDriverService/LoadPackages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoPackagesDriverServiceServer).LoadPackages(ctx, req.(*LoadPackagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GoPackagesDriverService_ServiceDesc is the grpc.ServiceDesc for GoPackagesDriverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoPackagesDriverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gopackagesdriver.GoPackagesDriverService",
	HandlerType: (*GoPackagesDriverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoadPackages",
			Handler:    _GoPackagesDriverService_LoadPackages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gopackagesdriver.proto",
}
