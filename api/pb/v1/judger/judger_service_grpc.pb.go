// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: judger/judger_service.proto

package pb_jg

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

// CodeClient is the client API for Code service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeClient interface {
	// 运行代码
	RunCode(ctx context.Context, in *RunCodeRequest, opts ...grpc.CallOption) (*RunCodeResponse, error)
	// 判题
	JudgeCode(ctx context.Context, in *JudgeCodeRequest, opts ...grpc.CallOption) (*JudgeCodeResponse, error)
}

type codeClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeClient(cc grpc.ClientConnInterface) CodeClient {
	return &codeClient{cc}
}

func (c *codeClient) RunCode(ctx context.Context, in *RunCodeRequest, opts ...grpc.CallOption) (*RunCodeResponse, error) {
	out := new(RunCodeResponse)
	err := c.cc.Invoke(ctx, "/v1.judger.Code/RunCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *codeClient) JudgeCode(ctx context.Context, in *JudgeCodeRequest, opts ...grpc.CallOption) (*JudgeCodeResponse, error) {
	out := new(JudgeCodeResponse)
	err := c.cc.Invoke(ctx, "/v1.judger.Code/JudgeCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeServer is the server API for Code service.
// All implementations must embed UnimplementedCodeServer
// for forward compatibility
type CodeServer interface {
	// 运行代码
	RunCode(context.Context, *RunCodeRequest) (*RunCodeResponse, error)
	// 判题
	JudgeCode(context.Context, *JudgeCodeRequest) (*JudgeCodeResponse, error)
	mustEmbedUnimplementedCodeServer()
}

// UnimplementedCodeServer must be embedded to have forward compatible implementations.
type UnimplementedCodeServer struct {
}

func (UnimplementedCodeServer) RunCode(context.Context, *RunCodeRequest) (*RunCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunCode not implemented")
}
func (UnimplementedCodeServer) JudgeCode(context.Context, *JudgeCodeRequest) (*JudgeCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JudgeCode not implemented")
}
func (UnimplementedCodeServer) mustEmbedUnimplementedCodeServer() {}

// UnsafeCodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeServer will
// result in compilation errors.
type UnsafeCodeServer interface {
	mustEmbedUnimplementedCodeServer()
}

func RegisterCodeServer(s grpc.ServiceRegistrar, srv CodeServer) {
	s.RegisterService(&Code_ServiceDesc, srv)
}

func _Code_RunCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeServer).RunCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.judger.Code/RunCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeServer).RunCode(ctx, req.(*RunCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Code_JudgeCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JudgeCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeServer).JudgeCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.judger.Code/JudgeCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeServer).JudgeCode(ctx, req.(*JudgeCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Code_ServiceDesc is the grpc.ServiceDesc for Code service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Code_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.judger.Code",
	HandlerType: (*CodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunCode",
			Handler:    _Code_RunCode_Handler,
		},
		{
			MethodName: "JudgeCode",
			Handler:    _Code_JudgeCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "judger/judger_service.proto",
}
