// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tabi/captain-node/v1/captain-node.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// Params defines captain-node module's parameters
type Params struct {
	MaxCaptains          uint32 `protobuf:"varint,1,opt,name=max_captains,json=maxCaptains,proto3" json:"max_captains,omitempty"`
	MinimumPowerOnPeriod uint32 `protobuf:"varint,2,opt,name=minimum_power_on_period,json=minimumPowerOnPeriod,proto3" json:"minimum_power_on_period,omitempty"`
	MaximumPowerOnPeriod uint32 `protobuf:"varint,3,opt,name=maximum_power_on_period,json=maximumPowerOnPeriod,proto3" json:"maximum_power_on_period,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_cdb21c7c329a6c9e, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMaxCaptains() uint32 {
	if m != nil {
		return m.MaxCaptains
	}
	return 0
}

func (m *Params) GetMinimumPowerOnPeriod() uint32 {
	if m != nil {
		return m.MinimumPowerOnPeriod
	}
	return 0
}

func (m *Params) GetMaximumPowerOnPeriod() uint32 {
	if m != nil {
		return m.MaximumPowerOnPeriod
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "tabi.captain_node.v1.Params")
}

func init() {
	proto.RegisterFile("tabi/captain-node/v1/captain-node.proto", fileDescriptor_cdb21c7c329a6c9e)
}

var fileDescriptor_cdb21c7c329a6c9e = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2f, 0x49, 0x4c, 0xca,
	0xd4, 0x4f, 0x4e, 0x2c, 0x28, 0x49, 0xcc, 0xcc, 0xd3, 0xcd, 0xcb, 0x4f, 0x49, 0xd5, 0x2f, 0x33,
	0x44, 0xe1, 0xeb, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0x89, 0x80, 0x14, 0xea, 0x41, 0x25, 0xe2,
	0xc1, 0x12, 0x65, 0x86, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x05, 0xfa, 0x20, 0x16, 0x44,
	0xad, 0xd2, 0x4c, 0x46, 0x2e, 0xb6, 0x80, 0xc4, 0xa2, 0xc4, 0xdc, 0x62, 0x21, 0x45, 0x2e, 0x9e,
	0xdc, 0xc4, 0x8a, 0x78, 0xa8, 0xbe, 0x62, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xde, 0x20, 0xee, 0xdc,
	0xc4, 0x0a, 0x67, 0xa8, 0x90, 0x90, 0x29, 0x97, 0x78, 0x6e, 0x66, 0x5e, 0x66, 0x6e, 0x69, 0x6e,
	0x7c, 0x41, 0x7e, 0x79, 0x6a, 0x51, 0x7c, 0x7e, 0x5e, 0x7c, 0x41, 0x6a, 0x51, 0x66, 0x7e, 0x8a,
	0x04, 0x13, 0x58, 0xb5, 0x08, 0x54, 0x3a, 0x00, 0x24, 0xeb, 0x9f, 0x17, 0x00, 0x96, 0x03, 0x6b,
	0x4b, 0xac, 0xc0, 0xaa, 0x8d, 0x19, 0xaa, 0x0d, 0x22, 0x8d, 0xa2, 0xcd, 0xc9, 0xfd, 0xc4, 0x23,
	0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2,
	0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x74, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4,
	0x92, 0xf3, 0x73, 0xf5, 0x41, 0x9e, 0xcd, 0x49, 0x4c, 0x2a, 0x06, 0x33, 0xf4, 0x2b, 0x50, 0x03,
	0xa8, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0xec, 0x57, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xe4, 0x84, 0x78, 0x1f, 0x42, 0x01, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaximumPowerOnPeriod != 0 {
		i = encodeVarintCaptainNode(dAtA, i, uint64(m.MaximumPowerOnPeriod))
		i--
		dAtA[i] = 0x18
	}
	if m.MinimumPowerOnPeriod != 0 {
		i = encodeVarintCaptainNode(dAtA, i, uint64(m.MinimumPowerOnPeriod))
		i--
		dAtA[i] = 0x10
	}
	if m.MaxCaptains != 0 {
		i = encodeVarintCaptainNode(dAtA, i, uint64(m.MaxCaptains))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintCaptainNode(dAtA []byte, offset int, v uint64) int {
	offset -= sovCaptainNode(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxCaptains != 0 {
		n += 1 + sovCaptainNode(uint64(m.MaxCaptains))
	}
	if m.MinimumPowerOnPeriod != 0 {
		n += 1 + sovCaptainNode(uint64(m.MinimumPowerOnPeriod))
	}
	if m.MaximumPowerOnPeriod != 0 {
		n += 1 + sovCaptainNode(uint64(m.MaximumPowerOnPeriod))
	}
	return n
}

func sovCaptainNode(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCaptainNode(x uint64) (n int) {
	return sovCaptainNode(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCaptainNode
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxCaptains", wireType)
			}
			m.MaxCaptains = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaptainNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxCaptains |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumPowerOnPeriod", wireType)
			}
			m.MinimumPowerOnPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaptainNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinimumPowerOnPeriod |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaximumPowerOnPeriod", wireType)
			}
			m.MaximumPowerOnPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaptainNode
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaximumPowerOnPeriod |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCaptainNode(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCaptainNode
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
func skipCaptainNode(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCaptainNode
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
					return 0, ErrIntOverflowCaptainNode
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
					return 0, ErrIntOverflowCaptainNode
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
				return 0, ErrInvalidLengthCaptainNode
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCaptainNode
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCaptainNode
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCaptainNode        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCaptainNode          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCaptainNode = fmt.Errorf("proto: unexpected end of group")
)
