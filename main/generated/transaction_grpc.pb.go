// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: transaction.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TransactionService_HandleTransaction_FullMethodName = "/TransactionService.TransactionService/HandleTransaction"
)

// TransactionServiceClient is the client API for TransactionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionServiceClient interface {
	HandleTransaction(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error)
}

type transactionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionServiceClient(cc grpc.ClientConnInterface) TransactionServiceClient {
	return &transactionServiceClient{cc}
}

func (c *transactionServiceClient) HandleTransaction(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (*TransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, TransactionService_HandleTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionServiceServer is the server API for TransactionService service.
// All implementations must embed UnimplementedTransactionServiceServer
// for forward compatibility.
type TransactionServiceServer interface {
	HandleTransaction(context.Context, *TransactionRequest) (*TransactionResponse, error)
	mustEmbedUnimplementedTransactionServiceServer()
}

// UnimplementedTransactionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTransactionServiceServer struct{}

func (UnimplementedTransactionServiceServer) HandleTransaction(context.Context, *TransactionRequest) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleTransaction not implemented")
}
func (UnimplementedTransactionServiceServer) mustEmbedUnimplementedTransactionServiceServer() {}
func (UnimplementedTransactionServiceServer) testEmbeddedByValue()                            {}

// UnsafeTransactionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionServiceServer will
// result in compilation errors.
type UnsafeTransactionServiceServer interface {
	mustEmbedUnimplementedTransactionServiceServer()
}

func RegisterTransactionServiceServer(s grpc.ServiceRegistrar, srv TransactionServiceServer) {
	// If the following call pancis, it indicates UnimplementedTransactionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TransactionService_ServiceDesc, srv)
}

func _TransactionService_HandleTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).HandleTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TransactionService_HandleTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).HandleTransaction(ctx, req.(*TransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TransactionService_ServiceDesc is the grpc.ServiceDesc for TransactionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TransactionService.TransactionService",
	HandlerType: (*TransactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleTransaction",
			Handler:    _TransactionService_HandleTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction.proto",
}
