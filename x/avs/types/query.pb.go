// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/avs/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryAVSInfoReq struct {
	AVSAddres string `protobuf:"bytes,1,opt,name=avs_address,json=avsAddress,proto3" json:"avs_address,omitempty"`
}

func (m *QueryAVSInfoReq) Reset()         { *m = QueryAVSInfoReq{} }
func (m *QueryAVSInfoReq) String() string { return proto.CompactTextString(m) }
func (*QueryAVSInfoReq) ProtoMessage()    {}
func (*QueryAVSInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d9df5f5965203ddc, []int{0}
}
func (m *QueryAVSInfoReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAVSInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAVSInfoReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAVSInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAVSInfoReq.Merge(m, src)
}
func (m *QueryAVSInfoReq) XXX_Size() int {
	return m.Size()
}
func (m *QueryAVSInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAVSInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAVSInfoReq proto.InternalMessageInfo

func (m *QueryAVSInfoReq) GetAVSAddres() string {
	if m != nil {
		return m.AVSAddres
	}
	return ""
}

type QueryAVSInfoResponse struct {
	Info *AVSInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (m *QueryAVSInfoResponse) Reset()         { *m = QueryAVSInfoResponse{} }
func (m *QueryAVSInfoResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAVSInfoResponse) ProtoMessage()    {}
func (*QueryAVSInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d9df5f5965203ddc, []int{1}
}
func (m *QueryAVSInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAVSInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAVSInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAVSInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAVSInfoResponse.Merge(m, src)
}
func (m *QueryAVSInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAVSInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAVSInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAVSInfoResponse proto.InternalMessageInfo

func (m *QueryAVSInfoResponse) GetInfo() *AVSInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryAVSInfoReq)(nil), "exocore.avs.QueryAVSInfoReq")
	proto.RegisterType((*QueryAVSInfoResponse)(nil), "exocore.avs.QueryAVSInfoResponse")
}

func init() { proto.RegisterFile("exocore/avs/query.proto", fileDescriptor_d9df5f5965203ddc) }

var fileDescriptor_d9df5f5965203ddc = []byte{
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0xad, 0xc8, 0x4f,
	0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0x2c, 0x2b, 0xd6, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x86, 0x4a, 0xe8, 0x25, 0x96, 0x15, 0x4b, 0x89, 0xa4, 0xe7, 0xa7,
	0xe7, 0x83, 0xc5, 0xf5, 0x41, 0x2c, 0x88, 0x12, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54,
	0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc, 0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62,
	0xa8, 0xac, 0x08, 0xb2, 0xc9, 0x25, 0x15, 0x10, 0x51, 0x25, 0x47, 0x2e, 0xfe, 0x40, 0x90, 0x2d,
	0x8e, 0x61, 0xc1, 0x9e, 0x79, 0x69, 0xf9, 0x41, 0xa9, 0x85, 0x42, 0x7a, 0x5c, 0xdc, 0x89, 0x65,
	0xc5, 0xf1, 0x89, 0x29, 0x29, 0x45, 0xa9, 0xc5, 0xc5, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x4e,
	0xbc, 0x8f, 0xee, 0xc9, 0x73, 0x3a, 0x86, 0x05, 0x3b, 0x82, 0x45, 0x83, 0xb8, 0x12, 0xcb, 0x8a,
	0x21, 0xcc, 0x62, 0x25, 0x07, 0x2e, 0x11, 0x54, 0x23, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85,
	0x34, 0xb8, 0x58, 0x32, 0xf3, 0xd2, 0xf2, 0xc1, 0x06, 0x70, 0x1b, 0x89, 0xe8, 0x21, 0x79, 0x40,
	0x0f, 0xa6, 0x16, 0xac, 0xc2, 0xa8, 0x9b, 0x91, 0x8b, 0x15, 0x6c, 0x84, 0x50, 0x23, 0x23, 0x17,
	0x0f, 0xb2, 0x61, 0x42, 0x32, 0x28, 0xda, 0xd0, 0x9c, 0x2a, 0xa5, 0x88, 0x47, 0x16, 0xe2, 0x0a,
	0x25, 0x83, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0x69, 0x09, 0x69, 0xe8, 0xbb, 0x42, 0x94, 0xfa, 0xa5,
	0x96, 0x94, 0xe7, 0x17, 0x65, 0xeb, 0x23, 0x07, 0x07, 0xb2, 0x4e, 0x27, 0xf7, 0x13, 0x8f, 0xe4,
	0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f,
	0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b,
	0xce, 0xcf, 0xc5, 0x65, 0x5a, 0x05, 0x24, 0x78, 0x2b, 0x0b, 0x52, 0x8b, 0x93, 0xd8, 0xc0, 0x41,
	0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xa5, 0x03, 0xeb, 0x15, 0xd4, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	QueryAVSInfo(ctx context.Context, in *QueryAVSInfoReq, opts ...grpc.CallOption) (*QueryAVSInfoResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) QueryAVSInfo(ctx context.Context, in *QueryAVSInfoReq, opts ...grpc.CallOption) (*QueryAVSInfoResponse, error) {
	out := new(QueryAVSInfoResponse)
	err := c.cc.Invoke(ctx, "/exocore.avs.Query/QueryAVSInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	QueryAVSInfo(context.Context, *QueryAVSInfoReq) (*QueryAVSInfoResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) QueryAVSInfo(ctx context.Context, req *QueryAVSInfoReq) (*QueryAVSInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAVSInfo not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_QueryAVSInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAVSInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryAVSInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exocore.avs.Query/QueryAVSInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryAVSInfo(ctx, req.(*QueryAVSInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "exocore.avs.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryAVSInfo",
			Handler:    _Query_QueryAVSInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exocore/avs/query.proto",
}

func (m *QueryAVSInfoReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAVSInfoReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAVSInfoReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AVSAddres) > 0 {
		i -= len(m.AVSAddres)
		copy(dAtA[i:], m.AVSAddres)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.AVSAddres)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAVSInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAVSInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAVSInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Info != nil {
		{
			size, err := m.Info.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryAVSInfoReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AVSAddres)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAVSInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Info != nil {
		l = m.Info.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryAVSInfoReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAVSInfoReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAVSInfoReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AVSAddres", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AVSAddres = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAVSInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAVSInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAVSInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Info", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Info == nil {
				m.Info = &AVSInfo{}
			}
			if err := m.Info.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
