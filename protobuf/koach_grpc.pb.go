// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: koach.proto

package protobuf

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

// ObservabilityServiceClient is the client API for ObservabilityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObservabilityServiceClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	ListenAlert(ctx context.Context, in *ListenAlertRequest, opts ...grpc.CallOption) (ObservabilityService_ListenAlertClient, error)
}

type observabilityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewObservabilityServiceClient(cc grpc.ClientConnInterface) ObservabilityServiceClient {
	return &observabilityServiceClient{cc}
}

func (c *observabilityServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/proto.ObservabilityService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *observabilityServiceClient) ListenAlert(ctx context.Context, in *ListenAlertRequest, opts ...grpc.CallOption) (ObservabilityService_ListenAlertClient, error) {
	stream, err := c.cc.NewStream(ctx, &ObservabilityService_ServiceDesc.Streams[0], "/proto.ObservabilityService/ListenAlert", opts...)
	if err != nil {
		return nil, err
	}
	x := &observabilityServiceListenAlertClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ObservabilityService_ListenAlertClient interface {
	Recv() (*ListenAlertResponse, error)
	grpc.ClientStream
}

type observabilityServiceListenAlertClient struct {
	grpc.ClientStream
}

func (x *observabilityServiceListenAlertClient) Recv() (*ListenAlertResponse, error) {
	m := new(ListenAlertResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObservabilityServiceServer is the server API for ObservabilityService service.
// All implementations must embed UnimplementedObservabilityServiceServer
// for forward compatibility
type ObservabilityServiceServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	ListenAlert(*ListenAlertRequest, ObservabilityService_ListenAlertServer) error
	mustEmbedUnimplementedObservabilityServiceServer()
}

// UnimplementedObservabilityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedObservabilityServiceServer struct {
}

func (UnimplementedObservabilityServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedObservabilityServiceServer) ListenAlert(*ListenAlertRequest, ObservabilityService_ListenAlertServer) error {
	return status.Errorf(codes.Unimplemented, "method ListenAlert not implemented")
}
func (UnimplementedObservabilityServiceServer) mustEmbedUnimplementedObservabilityServiceServer() {}

// UnsafeObservabilityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObservabilityServiceServer will
// result in compilation errors.
type UnsafeObservabilityServiceServer interface {
	mustEmbedUnimplementedObservabilityServiceServer()
}

func RegisterObservabilityServiceServer(s grpc.ServiceRegistrar, srv ObservabilityServiceServer) {
	s.RegisterService(&ObservabilityService_ServiceDesc, srv)
}

func _ObservabilityService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObservabilityServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ObservabilityService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObservabilityServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObservabilityService_ListenAlert_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenAlertRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObservabilityServiceServer).ListenAlert(m, &observabilityServiceListenAlertServer{stream})
}

type ObservabilityService_ListenAlertServer interface {
	Send(*ListenAlertResponse) error
	grpc.ServerStream
}

type observabilityServiceListenAlertServer struct {
	grpc.ServerStream
}

func (x *observabilityServiceListenAlertServer) Send(m *ListenAlertResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ObservabilityService_ServiceDesc is the grpc.ServiceDesc for ObservabilityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ObservabilityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ObservabilityService",
	HandlerType: (*ObservabilityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _ObservabilityService_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenAlert",
			Handler:       _ObservabilityService_ListenAlert_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "koach.proto",
}
