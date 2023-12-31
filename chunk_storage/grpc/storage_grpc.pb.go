// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: storage.proto

package grpc

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

const (
	StorageService_Save_FullMethodName         = "/proto.StorageService/Save"
	StorageService_Get_FullMethodName          = "/proto.StorageService/Get"
	StorageService_Erase_FullMethodName        = "/proto.StorageService/Erase"
	StorageService_GetFreeSpace_FullMethodName = "/proto.StorageService/GetFreeSpace"
)

// StorageServiceClient is the client API for StorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageServiceClient interface {
	Save(ctx context.Context, opts ...grpc.CallOption) (StorageService_SaveClient, error)
	Get(ctx context.Context, in *Get_Request, opts ...grpc.CallOption) (StorageService_GetClient, error)
	Erase(ctx context.Context, in *Erase_Request, opts ...grpc.CallOption) (*Erase_Response, error)
	GetFreeSpace(ctx context.Context, in *GetFreeSpace_Request, opts ...grpc.CallOption) (*GetFreeSpace_Response, error)
}

type storageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageServiceClient(cc grpc.ClientConnInterface) StorageServiceClient {
	return &storageServiceClient{cc}
}

func (c *storageServiceClient) Save(ctx context.Context, opts ...grpc.CallOption) (StorageService_SaveClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[0], StorageService_Save_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceSaveClient{stream}
	return x, nil
}

type StorageService_SaveClient interface {
	Send(*Save_Request) error
	CloseAndRecv() (*Save_Response, error)
	grpc.ClientStream
}

type storageServiceSaveClient struct {
	grpc.ClientStream
}

func (x *storageServiceSaveClient) Send(m *Save_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageServiceSaveClient) CloseAndRecv() (*Save_Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Save_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) Get(ctx context.Context, in *Get_Request, opts ...grpc.CallOption) (StorageService_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[1], StorageService_Get_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &storageServiceGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StorageService_GetClient interface {
	Recv() (*Get_Response, error)
	grpc.ClientStream
}

type storageServiceGetClient struct {
	grpc.ClientStream
}

func (x *storageServiceGetClient) Recv() (*Get_Response, error) {
	m := new(Get_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageServiceClient) Erase(ctx context.Context, in *Erase_Request, opts ...grpc.CallOption) (*Erase_Response, error) {
	out := new(Erase_Response)
	err := c.cc.Invoke(ctx, StorageService_Erase_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) GetFreeSpace(ctx context.Context, in *GetFreeSpace_Request, opts ...grpc.CallOption) (*GetFreeSpace_Response, error) {
	out := new(GetFreeSpace_Response)
	err := c.cc.Invoke(ctx, StorageService_GetFreeSpace_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServiceServer is the server API for StorageService service.
// All implementations must embed UnimplementedStorageServiceServer
// for forward compatibility
type StorageServiceServer interface {
	Save(StorageService_SaveServer) error
	Get(*Get_Request, StorageService_GetServer) error
	Erase(context.Context, *Erase_Request) (*Erase_Response, error)
	GetFreeSpace(context.Context, *GetFreeSpace_Request) (*GetFreeSpace_Response, error)
	mustEmbedUnimplementedStorageServiceServer()
}

// UnimplementedStorageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServiceServer struct {
}

func (UnimplementedStorageServiceServer) Save(StorageService_SaveServer) error {
	return status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedStorageServiceServer) Get(*Get_Request, StorageService_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedStorageServiceServer) Erase(context.Context, *Erase_Request) (*Erase_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Erase not implemented")
}
func (UnimplementedStorageServiceServer) GetFreeSpace(context.Context, *GetFreeSpace_Request) (*GetFreeSpace_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFreeSpace not implemented")
}
func (UnimplementedStorageServiceServer) mustEmbedUnimplementedStorageServiceServer() {}

// UnsafeStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServiceServer will
// result in compilation errors.
type UnsafeStorageServiceServer interface {
	mustEmbedUnimplementedStorageServiceServer()
}

func RegisterStorageServiceServer(s grpc.ServiceRegistrar, srv StorageServiceServer) {
	s.RegisterService(&StorageService_ServiceDesc, srv)
}

func _StorageService_Save_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServiceServer).Save(&storageServiceSaveServer{stream})
}

type StorageService_SaveServer interface {
	SendAndClose(*Save_Response) error
	Recv() (*Save_Request, error)
	grpc.ServerStream
}

type storageServiceSaveServer struct {
	grpc.ServerStream
}

func (x *storageServiceSaveServer) SendAndClose(m *Save_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageServiceSaveServer) Recv() (*Save_Request, error) {
	m := new(Save_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StorageService_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Get_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).Get(m, &storageServiceGetServer{stream})
}

type StorageService_GetServer interface {
	Send(*Get_Response) error
	grpc.ServerStream
}

type storageServiceGetServer struct {
	grpc.ServerStream
}

func (x *storageServiceGetServer) Send(m *Get_Response) error {
	return x.ServerStream.SendMsg(m)
}

func _StorageService_Erase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Erase_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).Erase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageService_Erase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).Erase(ctx, req.(*Erase_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_GetFreeSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFreeSpace_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).GetFreeSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageService_GetFreeSpace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).GetFreeSpace(ctx, req.(*GetFreeSpace_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// StorageService_ServiceDesc is the grpc.ServiceDesc for StorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.StorageService",
	HandlerType: (*StorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Erase",
			Handler:    _StorageService_Erase_Handler,
		},
		{
			MethodName: "GetFreeSpace",
			Handler:    _StorageService_GetFreeSpace_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Save",
			Handler:       _StorageService_Save_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Get",
			Handler:       _StorageService_Get_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage.proto",
}
