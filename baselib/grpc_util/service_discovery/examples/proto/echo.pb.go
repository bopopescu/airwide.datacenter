// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	echo.proto

It has these top-level messages:
	EchoReq
	EchoRsp
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type EchoReq struct {
	EchoData string `protobuf:"bytes,1,opt,name=echo_data,json=echoData" json:"echo_data,omitempty"`
}

func (m *EchoReq) Reset()                    { *m = EchoReq{} }
func (m *EchoReq) String() string            { return proto1.CompactTextString(m) }
func (*EchoReq) ProtoMessage()               {}
func (*EchoReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoReq) GetEchoData() string {
	if m != nil {
		return m.EchoData
	}
	return ""
}

type EchoRsp struct {
	EchoData string `protobuf:"bytes,1,opt,name=echo_data,json=echoData" json:"echo_data,omitempty"`
}

func (m *EchoRsp) Reset()                    { *m = EchoRsp{} }
func (m *EchoRsp) String() string            { return proto1.CompactTextString(m) }
func (*EchoRsp) ProtoMessage()               {}
func (*EchoRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EchoRsp) GetEchoData() string {
	if m != nil {
		return m.EchoData
	}
	return ""
}

func init() {
	proto1.RegisterType((*EchoReq)(nil), "proto.EchoReq")
	proto1.RegisterType((*EchoRsp)(nil), "proto.EchoRsp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EchoService service

type EchoServiceClient interface {
	Echo(ctx context.Context, in *EchoReq, opts ...grpc.CallOption) (*EchoRsp, error)
}

type echoServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoServiceClient(cc *grpc.ClientConn) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) Echo(ctx context.Context, in *EchoReq, opts ...grpc.CallOption) (*EchoRsp, error) {
	out := new(EchoRsp)
	err := grpc.Invoke(ctx, "/proto.EchoService/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EchoService service

type EchoServiceServer interface {
	Echo(context.Context, *EchoReq) (*EchoRsp, error)
}

func RegisterEchoServiceServer(s *grpc.Server, srv EchoServiceServer) {
	s.RegisterService(&_EchoService_serviceDesc, srv)
}

func _EchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Echo(ctx, req.(*EchoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}

func init() { proto1.RegisterFile("echo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 119 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4d, 0xce, 0xc8,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x6a, 0x5c, 0xec, 0xae, 0xc9,
	0x19, 0xf9, 0x41, 0xa9, 0x85, 0x42, 0xd2, 0x5c, 0x9c, 0x20, 0xf9, 0xf8, 0x94, 0xc4, 0x92, 0x44,
	0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x0e, 0x90, 0x80, 0x4b, 0x62, 0x49, 0x22, 0x5c, 0x5d,
	0x71, 0x01, 0x5e, 0x75, 0x46, 0xe6, 0x5c, 0xdc, 0x20, 0x75, 0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9,
	0xa9, 0x42, 0x1a, 0x5c, 0x2c, 0x20, 0xae, 0x10, 0x1f, 0xc4, 0x56, 0x3d, 0xa8, 0x5d, 0x52, 0x28,
	0xfc, 0xe2, 0x02, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x80, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x52,
	0x99, 0x06, 0x45, 0xa4, 0x00, 0x00, 0x00,
}
