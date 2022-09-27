// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api.proto

package api

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

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminClient interface {
	ActivityCreate(ctx context.Context, in *ActivityCreateRequest, opts ...grpc.CallOption) (*ActivityCreateResponse, error)
	ActivityList(ctx context.Context, in *ActivityListRequest, opts ...grpc.CallOption) (*ActivityListResponse, error)
	ActivityListStream(ctx context.Context, in *ActivityListStreamRequest, opts ...grpc.CallOption) (Admin_ActivityListStreamClient, error)
	ActivityToday(ctx context.Context, in *ActivityTodayRequest, opts ...grpc.CallOption) (*ActivityTodayResponse, error)
	ActivityGet(ctx context.Context, in *ActivityGetRequest, opts ...grpc.CallOption) (*ActivityGetResponse, error)
	ActivityUpdate(ctx context.Context, in *ActivityUpdateRequest, opts ...grpc.CallOption) (*ActivityUpdateResponse, error)
	ActivityDelete(ctx context.Context, in *ActivityDeleteRequest, opts ...grpc.CallOption) (*ActivityDeleteResponse, error)
}

type adminClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) ActivityCreate(ctx context.Context, in *ActivityCreateRequest, opts ...grpc.CallOption) (*ActivityCreateResponse, error) {
	out := new(ActivityCreateResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ActivityList(ctx context.Context, in *ActivityListRequest, opts ...grpc.CallOption) (*ActivityListResponse, error) {
	out := new(ActivityListResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ActivityListStream(ctx context.Context, in *ActivityListStreamRequest, opts ...grpc.CallOption) (Admin_ActivityListStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Admin_ServiceDesc.Streams[0], "/api.Admin/ActivityListStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &adminActivityListStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Admin_ActivityListStreamClient interface {
	Recv() (*ActivityListStreamResponse, error)
	grpc.ClientStream
}

type adminActivityListStreamClient struct {
	grpc.ClientStream
}

func (x *adminActivityListStreamClient) Recv() (*ActivityListStreamResponse, error) {
	m := new(ActivityListStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *adminClient) ActivityToday(ctx context.Context, in *ActivityTodayRequest, opts ...grpc.CallOption) (*ActivityTodayResponse, error) {
	out := new(ActivityTodayResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityToday", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ActivityGet(ctx context.Context, in *ActivityGetRequest, opts ...grpc.CallOption) (*ActivityGetResponse, error) {
	out := new(ActivityGetResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ActivityUpdate(ctx context.Context, in *ActivityUpdateRequest, opts ...grpc.CallOption) (*ActivityUpdateResponse, error) {
	out := new(ActivityUpdateResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ActivityDelete(ctx context.Context, in *ActivityDeleteRequest, opts ...grpc.CallOption) (*ActivityDeleteResponse, error) {
	out := new(ActivityDeleteResponse)
	err := c.cc.Invoke(ctx, "/api.Admin/ActivityDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
// All implementations must embed UnimplementedAdminServer
// for forward compatibility
type AdminServer interface {
	ActivityCreate(context.Context, *ActivityCreateRequest) (*ActivityCreateResponse, error)
	ActivityList(context.Context, *ActivityListRequest) (*ActivityListResponse, error)
	ActivityListStream(*ActivityListStreamRequest, Admin_ActivityListStreamServer) error
	ActivityToday(context.Context, *ActivityTodayRequest) (*ActivityTodayResponse, error)
	ActivityGet(context.Context, *ActivityGetRequest) (*ActivityGetResponse, error)
	ActivityUpdate(context.Context, *ActivityUpdateRequest) (*ActivityUpdateResponse, error)
	ActivityDelete(context.Context, *ActivityDeleteRequest) (*ActivityDeleteResponse, error)
	mustEmbedUnimplementedAdminServer()
}

// UnimplementedAdminServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (UnimplementedAdminServer) ActivityCreate(context.Context, *ActivityCreateRequest) (*ActivityCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityCreate not implemented")
}
func (UnimplementedAdminServer) ActivityList(context.Context, *ActivityListRequest) (*ActivityListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityList not implemented")
}
func (UnimplementedAdminServer) ActivityListStream(*ActivityListStreamRequest, Admin_ActivityListStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ActivityListStream not implemented")
}
func (UnimplementedAdminServer) ActivityToday(context.Context, *ActivityTodayRequest) (*ActivityTodayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityToday not implemented")
}
func (UnimplementedAdminServer) ActivityGet(context.Context, *ActivityGetRequest) (*ActivityGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityGet not implemented")
}
func (UnimplementedAdminServer) ActivityUpdate(context.Context, *ActivityUpdateRequest) (*ActivityUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityUpdate not implemented")
}
func (UnimplementedAdminServer) ActivityDelete(context.Context, *ActivityDeleteRequest) (*ActivityDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivityDelete not implemented")
}
func (UnimplementedAdminServer) mustEmbedUnimplementedAdminServer() {}

// UnsafeAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServer will
// result in compilation errors.
type UnsafeAdminServer interface {
	mustEmbedUnimplementedAdminServer()
}

func RegisterAdminServer(s grpc.ServiceRegistrar, srv AdminServer) {
	s.RegisterService(&Admin_ServiceDesc, srv)
}

func _Admin_ActivityCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityCreate(ctx, req.(*ActivityCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ActivityList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityList(ctx, req.(*ActivityListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ActivityListStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ActivityListStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AdminServer).ActivityListStream(m, &adminActivityListStreamServer{stream})
}

type Admin_ActivityListStreamServer interface {
	Send(*ActivityListStreamResponse) error
	grpc.ServerStream
}

type adminActivityListStreamServer struct {
	grpc.ServerStream
}

func (x *adminActivityListStreamServer) Send(m *ActivityListStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Admin_ActivityToday_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityTodayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityToday(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityToday",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityToday(ctx, req.(*ActivityTodayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ActivityGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityGet(ctx, req.(*ActivityGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ActivityUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityUpdate(ctx, req.(*ActivityUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ActivityDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivityDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ActivityDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Admin/ActivityDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ActivityDelete(ctx, req.(*ActivityDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Admin_ServiceDesc is the grpc.ServiceDesc for Admin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Admin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ActivityCreate",
			Handler:    _Admin_ActivityCreate_Handler,
		},
		{
			MethodName: "ActivityList",
			Handler:    _Admin_ActivityList_Handler,
		},
		{
			MethodName: "ActivityToday",
			Handler:    _Admin_ActivityToday_Handler,
		},
		{
			MethodName: "ActivityGet",
			Handler:    _Admin_ActivityGet_Handler,
		},
		{
			MethodName: "ActivityUpdate",
			Handler:    _Admin_ActivityUpdate_Handler,
		},
		{
			MethodName: "ActivityDelete",
			Handler:    _Admin_ActivityDelete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ActivityListStream",
			Handler:       _Admin_ActivityListStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api.proto",
}
