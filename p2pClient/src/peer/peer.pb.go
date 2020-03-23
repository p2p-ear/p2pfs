// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer.proto

package peer

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PingMessage struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=Ok,proto3" json:"Ok,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingMessage) Reset()         { *m = PingMessage{} }
func (m *PingMessage) String() string { return proto.CompactTextString(m) }
func (*PingMessage) ProtoMessage()    {}
func (*PingMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_055ae5a865fc1c9e, []int{0}
}

func (m *PingMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingMessage.Unmarshal(m, b)
}
func (m *PingMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingMessage.Marshal(b, m, deterministic)
}
func (m *PingMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingMessage.Merge(m, src)
}
func (m *PingMessage) XXX_Size() int {
	return xxx_messageInfo_PingMessage.Size(m)
}
func (m *PingMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PingMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PingMessage proto.InternalMessageInfo

func (m *PingMessage) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_055ae5a865fc1c9e, []int{1}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type WriteRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WriteRequest) Reset()         { *m = WriteRequest{} }
func (m *WriteRequest) String() string { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()    {}
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_055ae5a865fc1c9e, []int{2}
}

func (m *WriteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteRequest.Unmarshal(m, b)
}
func (m *WriteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteRequest.Marshal(b, m, deterministic)
}
func (m *WriteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteRequest.Merge(m, src)
}
func (m *WriteRequest) XXX_Size() int {
	return xxx_messageInfo_WriteRequest.Size(m)
}
func (m *WriteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WriteRequest proto.InternalMessageInfo

func (m *WriteRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *WriteRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReadRequest struct {
	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	//Maximum amount of bytes to read (inf if 0)
	Maxbytes             uint64   `protobuf:"varint,2,opt,name=Maxbytes,proto3" json:"Maxbytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReadRequest) Reset()         { *m = ReadRequest{} }
func (m *ReadRequest) String() string { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()    {}
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_055ae5a865fc1c9e, []int{3}
}

func (m *ReadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadRequest.Unmarshal(m, b)
}
func (m *ReadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadRequest.Marshal(b, m, deterministic)
}
func (m *ReadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadRequest.Merge(m, src)
}
func (m *ReadRequest) XXX_Size() int {
	return xxx_messageInfo_ReadRequest.Size(m)
}
func (m *ReadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReadRequest proto.InternalMessageInfo

func (m *ReadRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ReadRequest) GetMaxbytes() uint64 {
	if m != nil {
		return m.Maxbytes
	}
	return 0
}

type ReadReply struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReadReply) Reset()         { *m = ReadReply{} }
func (m *ReadReply) String() string { return proto.CompactTextString(m) }
func (*ReadReply) ProtoMessage()    {}
func (*ReadReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_055ae5a865fc1c9e, []int{4}
}

func (m *ReadReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadReply.Unmarshal(m, b)
}
func (m *ReadReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadReply.Marshal(b, m, deterministic)
}
func (m *ReadReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadReply.Merge(m, src)
}
func (m *ReadReply) XXX_Size() int {
	return xxx_messageInfo_ReadReply.Size(m)
}
func (m *ReadReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadReply.DiscardUnknown(m)
}

var xxx_messageInfo_ReadReply proto.InternalMessageInfo

func (m *ReadReply) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*PingMessage)(nil), "peer.PingMessage")
	proto.RegisterType((*Empty)(nil), "peer.Empty")
	proto.RegisterType((*WriteRequest)(nil), "peer.WriteRequest")
	proto.RegisterType((*ReadRequest)(nil), "peer.ReadRequest")
	proto.RegisterType((*ReadReply)(nil), "peer.ReadReply")
}

func init() {
	proto.RegisterFile("peer.proto", fileDescriptor_055ae5a865fc1c9e)
}

var fileDescriptor_055ae5a865fc1c9e = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4b, 0x03, 0x31,
	0x14, 0x84, 0x37, 0x4b, 0xaa, 0xed, 0xdb, 0xa2, 0xf8, 0x4e, 0x25, 0x20, 0x96, 0x9c, 0x8a, 0xc8,
	0x1e, 0x14, 0xbc, 0x79, 0xd3, 0x63, 0x6d, 0x89, 0x07, 0xcf, 0xa9, 0x3e, 0xca, 0xd2, 0xd6, 0xc6,
	0x24, 0x8a, 0xf9, 0x25, 0xfe, 0x5d, 0x49, 0x82, 0x6b, 0x40, 0xe8, 0x6d, 0x32, 0xc9, 0x0c, 0xf9,
	0x06, 0xc0, 0x10, 0xd9, 0xd6, 0xd8, 0xbd, 0xdf, 0x23, 0x8f, 0x5a, 0x9e, 0x43, 0xb3, 0xec, 0xde,
	0xd6, 0x73, 0x72, 0x4e, 0xaf, 0x09, 0x4f, 0xa0, 0x5e, 0x6c, 0x26, 0x6c, 0xca, 0x66, 0x43, 0x55,
	0x2f, 0x36, 0xf2, 0x18, 0x06, 0x0f, 0x3b, 0xe3, 0x83, 0xbc, 0x85, 0xf1, 0xb3, 0xed, 0x3c, 0x29,
	0x7a, 0xff, 0x20, 0xe7, 0x11, 0x81, 0x3f, 0xea, 0x1d, 0xa5, 0xa7, 0x23, 0x95, 0x74, 0xf4, 0xee,
	0xb5, 0xd7, 0x93, 0x7a, 0xca, 0x66, 0x63, 0x95, 0xb4, 0xbc, 0x83, 0x46, 0x91, 0x7e, 0x3d, 0x14,
	0x13, 0x30, 0x9c, 0xeb, 0xaf, 0x55, 0xf0, 0xe4, 0x52, 0x94, 0xab, 0xfe, 0x2c, 0x2f, 0x60, 0x94,
	0xe3, 0x66, 0x1b, 0xfa, 0x7e, 0xf6, 0xd7, 0x7f, 0xfd, 0xcd, 0xa0, 0x59, 0x12, 0xd9, 0x27, 0xb2,
	0x9f, 0xdd, 0x0b, 0x61, 0x0b, 0x3c, 0xf2, 0xe0, 0x59, 0x9b, 0x50, 0x0b, 0x36, 0xf1, 0xdf, 0x92,
	0x15, 0x5e, 0xc2, 0x20, 0x71, 0x21, 0xe6, 0xdb, 0x12, 0x52, 0x34, 0xd9, 0xcb, 0x0b, 0x54, 0x78,
	0x05, 0x3c, 0x7e, 0xe6, 0xb7, 0xbb, 0xe0, 0x12, 0xa7, 0xa5, 0x65, 0xb6, 0x41, 0x56, 0xab, 0xa3,
	0x34, 0xf3, 0xcd, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbd, 0x19, 0x5c, 0x49, 0x74, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PeerServiceClient is the client API for PeerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PeerServiceClient interface {
	Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*Empty, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadReply, error)
}

type peerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPeerServiceClient(cc grpc.ClientConnInterface) PeerServiceClient {
	return &peerServiceClient{cc}
}

func (c *peerServiceClient) Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*PingMessage, error) {
	out := new(PingMessage)
	err := c.cc.Invoke(ctx, "/peer.PeerService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerServiceClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/peer.PeerService/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerServiceClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadReply, error) {
	out := new(ReadReply)
	err := c.cc.Invoke(ctx, "/peer.PeerService/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerServiceServer is the server API for PeerService service.
type PeerServiceServer interface {
	Ping(context.Context, *PingMessage) (*PingMessage, error)
	Write(context.Context, *WriteRequest) (*Empty, error)
	Read(context.Context, *ReadRequest) (*ReadReply, error)
}

// UnimplementedPeerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPeerServiceServer struct {
}

func (*UnimplementedPeerServiceServer) Ping(ctx context.Context, req *PingMessage) (*PingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedPeerServiceServer) Write(ctx context.Context, req *WriteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (*UnimplementedPeerServiceServer) Read(ctx context.Context, req *ReadRequest) (*ReadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}

func RegisterPeerServiceServer(s *grpc.Server, srv PeerServiceServer) {
	s.RegisterService(&_PeerService_serviceDesc, srv)
}

func _PeerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/peer.PeerService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServiceServer).Ping(ctx, req.(*PingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerService_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServiceServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/peer.PeerService/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServiceServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerService_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServiceServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/peer.PeerService/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServiceServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PeerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "peer.PeerService",
	HandlerType: (*PeerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _PeerService_Ping_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _PeerService_Write_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _PeerService_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peer.proto",
}