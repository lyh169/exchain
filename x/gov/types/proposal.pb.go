// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cosmos/gov/v1/tx2.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	types1 "github.com/okex/exchain/libs/cosmos-sdk/types"
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

// MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary
// proposal Content.
type ProtobufMsgSubmitProposal struct {
	Content        *types.Any           `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	InitialDeposit []types1.CoinAdapter `protobuf:"bytes,2,rep,name=initial_deposit,json=initialDeposit,proto3" json:"initial_deposit"`
	Proposer       string               `protobuf:"bytes,3,opt,name=proposer,proto3" json:"proposer,omitempty"`
}

func (m *ProtobufMsgSubmitProposal) Reset()         { *m = ProtobufMsgSubmitProposal{} }
func (m *ProtobufMsgSubmitProposal) String() string { return proto.CompactTextString(m) }
func (*ProtobufMsgSubmitProposal) ProtoMessage()    {}
func (*ProtobufMsgSubmitProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1bba2ff2655b54a, []int{0}
}
func (m *ProtobufMsgSubmitProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProtobufMsgSubmitProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProtobufMsgSubmitProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProtobufMsgSubmitProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProtobufMsgSubmitProposal.Merge(m, src)
}
func (m *ProtobufMsgSubmitProposal) XXX_Size() int {
	return m.Size()
}
func (m *ProtobufMsgSubmitProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_ProtobufMsgSubmitProposal.DiscardUnknown(m)
}

var xxx_messageInfo_ProtobufMsgSubmitProposal proto.InternalMessageInfo

func (m *ProtobufMsgSubmitProposal) GetContent() *types.Any {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *ProtobufMsgSubmitProposal) GetInitialDeposit() []types1.CoinAdapter {
	if m != nil {
		return m.InitialDeposit
	}
	return nil
}

func (m *ProtobufMsgSubmitProposal) GetProposer() string {
	if m != nil {
		return m.Proposer
	}
	return ""
}

func init() {
	proto.RegisterType((*ProtobufMsgSubmitProposal)(nil), "cosmos.gov.v1.ProtobufMsgSubmitProposal")
}

func init() { proto.RegisterFile("cosmos/gov/v1/tx2.proto", fileDescriptor_e1bba2ff2655b54a) }

var fileDescriptor_e1bba2ff2655b54a = []byte{
	// 359 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0x4e, 0xfa, 0x40,
	0x10, 0xc7, 0xdb, 0x1f, 0xbf, 0xf8, 0xa7, 0x04, 0x4d, 0x1a, 0x12, 0x0b, 0x87, 0xd2, 0x78, 0x22,
	0x1a, 0x76, 0x53, 0xf4, 0xe4, 0x0d, 0xd4, 0xa3, 0x91, 0xc0, 0xcd, 0x0b, 0xe9, 0x9f, 0x75, 0xdd,
	0x48, 0x77, 0x9a, 0xee, 0xd2, 0xc0, 0xd5, 0x27, 0xf0, 0x51, 0x3c, 0xf8, 0x10, 0x1c, 0x89, 0x27,
	0x4f, 0x46, 0xe1, 0xe0, 0x6b, 0x98, 0x76, 0x17, 0x0e, 0x9c, 0x3a, 0xd3, 0xef, 0x67, 0xe6, 0x3b,
	0x3b, 0x63, 0x9d, 0x44, 0x20, 0x12, 0x10, 0x98, 0x42, 0x8e, 0x73, 0x1f, 0xcb, 0x59, 0x17, 0xa5,
	0x19, 0x48, 0xb0, 0x6b, 0x4a, 0x40, 0x14, 0x72, 0x94, 0xfb, 0xcd, 0x96, 0xe6, 0xc2, 0x40, 0x10,
	0x9c, 0xfb, 0x21, 0x91, 0x81, 0x8f, 0x23, 0x60, 0x5c, 0xf3, 0xcd, 0x9d, 0x46, 0x45, 0x99, 0x12,
	0xea, 0x14, 0x28, 0x94, 0x21, 0x2e, 0x22, 0xfd, 0xb7, 0xa1, 0xf0, 0xb1, 0x12, 0xb4, 0x97, 0x96,
	0x28, 0x00, 0x9d, 0x10, 0x5c, 0x66, 0xe1, 0xf4, 0x11, 0x07, 0x7c, 0xbe, 0x63, 0x92, 0x08, 0x5a,
	0x98, 0x24, 0x82, 0x2a, 0xe1, 0xf4, 0xc7, 0xb4, 0x1a, 0x03, 0xcd, 0xdf, 0x09, 0x3a, 0x9a, 0x86,
	0x09, 0x93, 0x83, 0x0c, 0x52, 0x10, 0xc1, 0xc4, 0x46, 0xd6, 0x7e, 0x04, 0x5c, 0x12, 0x2e, 0x1d,
	0xd3, 0x33, 0xdb, 0xd5, 0x6e, 0x1d, 0x29, 0x0f, 0xb4, 0xf1, 0x40, 0x3d, 0x3e, 0x1f, 0x6e, 0x20,
	0xfb, 0xde, 0x3a, 0x66, 0x9c, 0x49, 0x16, 0x4c, 0xc6, 0x31, 0x49, 0x41, 0x30, 0xe9, 0xfc, 0xf3,
	0x2a, 0xed, 0x6a, 0xd7, 0x43, 0x7a, 0xd2, 0x62, 0x0d, 0x48, 0xaf, 0x01, 0x5d, 0x03, 0xe3, 0xbd,
	0x38, 0x48, 0x25, 0xc9, 0xfa, 0xff, 0x17, 0x5f, 0x2d, 0x63, 0x78, 0xa4, 0xcb, 0x6f, 0x54, 0xb5,
	0x7d, 0x69, 0x1d, 0xa4, 0xe5, 0x30, 0x24, 0x73, 0x2a, 0x9e, 0xd9, 0x3e, 0xec, 0x3b, 0x1f, 0xef,
	0x9d, 0xba, 0x6e, 0xd6, 0x8b, 0xe3, 0x8c, 0x08, 0x31, 0x92, 0x19, 0xe3, 0x74, 0xb8, 0x25, 0xaf,
	0x6a, 0x2f, 0xbf, 0x6f, 0x67, 0xdb, 0xb4, 0x7f, 0xbb, 0x58, 0xb9, 0xe6, 0x72, 0xe5, 0x9a, 0xdf,
	0x2b, 0xd7, 0x7c, 0x5d, 0xbb, 0xc6, 0x72, 0xed, 0x1a, 0x9f, 0x6b, 0xd7, 0x78, 0x38, 0xa7, 0x4c,
	0x3e, 0x4d, 0x43, 0x14, 0x41, 0xa2, 0x57, 0xa9, 0x3f, 0x1d, 0x11, 0x3f, 0xe3, 0x59, 0x79, 0x13,
	0x39, 0x4f, 0x89, 0x28, 0x2e, 0xb7, 0x57, 0xbe, 0xf9, 0xe2, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xb3,
	0x38, 0xf2, 0xf0, 0xfa, 0x01, 0x00, 0x00,
}

func (m *ProtobufMsgSubmitProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProtobufMsgSubmitProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProtobufMsgSubmitProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Proposer) > 0 {
		i -= len(m.Proposer)
		copy(dAtA[i:], m.Proposer)
		i = encodeVarintTx2(dAtA, i, uint64(len(m.Proposer)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.InitialDeposit) > 0 {
		for iNdEx := len(m.InitialDeposit) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.InitialDeposit[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx2(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Content != nil {
		{
			size, err := m.Content.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx2(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx2(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx2(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProtobufMsgSubmitProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Content != nil {
		l = m.Content.Size()
		n += 1 + l + sovTx2(uint64(l))
	}
	if len(m.InitialDeposit) > 0 {
		for _, e := range m.InitialDeposit {
			l = e.Size()
			n += 1 + l + sovTx2(uint64(l))
		}
	}
	l = len(m.Proposer)
	if l > 0 {
		n += 1 + l + sovTx2(uint64(l))
	}
	return n
}

func sovTx2(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx2(x uint64) (n int) {
	return sovTx2(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProtobufMsgSubmitProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx2
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
			return fmt.Errorf("proto: ProtobufMsgSubmitProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProtobufMsgSubmitProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx2
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
				return ErrInvalidLengthTx2
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx2
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Content == nil {
				m.Content = &types.Any{}
			}
			if err := m.Content.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialDeposit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx2
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
				return ErrInvalidLengthTx2
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx2
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InitialDeposit = append(m.InitialDeposit, types1.CoinAdapter{})
			if err := m.InitialDeposit[len(m.InitialDeposit)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proposer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx2
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
				return ErrInvalidLengthTx2
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx2
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proposer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx2(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx2
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
func skipTx2(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx2
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
					return 0, ErrIntOverflowTx2
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
					return 0, ErrIntOverflowTx2
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
				return 0, ErrInvalidLengthTx2
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx2
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx2
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx2        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx2          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx2 = fmt.Errorf("proto: unexpected end of group")
)
