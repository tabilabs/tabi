// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tabi/V2Incentives/v1/genesis.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// V2GenesisState defines the module's genesis state.
type V2GenesisState struct {
	// V2Params are the V2Incentives module parameters
	V2Params V2Params `protobuf:"bytes,1,opt,name=V2Params,proto3" json:"V2Params"`
	// V2Incentives is a slice of active V2Incentives
	V2Incentives []V2Incentive `protobuf:"bytes,2,rep,name=V2Incentives,proto3" json:"V2Incentives"`
	// gas_meters is a slice of active V2GasMeters
	V2GasMeters []V2GasMeter `protobuf:"bytes,3,rep,name=gas_meters,json=V2GasMeters,proto3" json:"gas_meters"`
}

func (m *V2GenesisState) Reset()         { *m = V2GenesisState{} }
func (m *V2GenesisState) String() string { return proto.CompactTextString(m) }
func (*V2GenesisState) ProtoMessage()    {}
func (*V2GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bb1f7c7e8ad160b, []int{0}
}

func (m *V2GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *V2GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_V2GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *V2GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_V2GenesisState.Merge(m, src)
}

func (m *V2GenesisState) XXX_Size() int {
	return m.Size()
}

func (m *V2GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_V2GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_V2GenesisState proto.InternalMessageInfo

func (m *V2GenesisState) GetV2Params() V2Params {
	if m != nil {
		return m.V2Params
	}
	return V2Params{}
}

func (m *V2GenesisState) GetV2Incentives() []V2Incentive {
	if m != nil {
		return m.V2Incentives
	}
	return nil
}

func (m *V2GenesisState) GetV2GasMeters() []V2GasMeter {
	if m != nil {
		return m.V2GasMeters
	}
	return nil
}

// V2Params defines the V2Incentives module V2Params
type V2Params struct {
	// enable_V2Incentives is the parameter to enable V2Incentives
	EnableIncentives bool `protobuf:"varint,1,opt,name=enable_V2Incentives,json=EnableIncentives,proto3" json:"enable_V2Incentives,omitempty"`
	// allocation_limit is the maximum percentage an V2Incentive can allocate per denomination
	AllocationLimit github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=allocation_limit,json=allocationLimit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"allocation_limit"`
	// V2Incentives_epoch_identifier for the epochs module hooks
	IncentivesEpochIdentifier string `protobuf:"bytes,3,opt,name=V2Incentives_epoch_identifier,json=IncentivesEpochIdentifier,proto3" json:"V2Incentives_epoch_identifier,omitempty"`
	// reward_scaler is the scaling factor for capping rewards
	RewardScaler github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=reward_scaler,json=rewardScaler,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_scaler"`
}

func (m *V2Params) Reset()         { *m = V2Params{} }
func (m *V2Params) String() string { return proto.CompactTextString(m) }
func (*V2Params) ProtoMessage()    {}
func (*V2Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bb1f7c7e8ad160b, []int{1}
}

func (m *V2Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *V2Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_V2Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *V2Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_V2Params.Merge(m, src)
}

func (m *V2Params) XXX_Size() int {
	return m.Size()
}

func (m *V2Params) XXX_DiscardUnknown() {
	xxx_messageInfo_V2Params.DiscardUnknown(m)
}

var xxx_messageInfo_V2Params proto.InternalMessageInfo

func (m *V2Params) GetEnableIncentives() bool {
	if m != nil {
		return m.EnableIncentives
	}
	return false
}

func (m *V2Params) GetIncentivesEpochIdentifier() string {
	if m != nil {
		return m.IncentivesEpochIdentifier
	}
	return ""
}

func init() {
	proto.RegisterType((*V2GenesisState)(nil), "tabi.V2Incentives.v1.V2GenesisState")
	proto.RegisterType((*V2Params)(nil), "tabi.V2Incentives.v1.V2Params")
}

func init() {
	proto.RegisterFile("tabi/V2Incentives/v1/genesis.proto", fileDescriptor_7bb1f7c7e8ad160b)
}

var fileDescriptor_7bb1f7c7e8ad160b = []byte{
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x6c, 0x92, 0xcf, 0x4a, 0xf3, 0x40,
	0x14, 0xc5, 0x6f, 0xa6, 0x1f, 0xe5, 0xeb, 0xb4, 0x62, 0x0d, 0x2e, 0x6a, 0x55, 0x28, 0x75, 0x53,
	0x28, 0x24, 0xa4, 0x6e, 0xa4, 0x3b, 0xe3, 0x9f, 0x52, 0x50, 0x90, 0x74, 0xe7, 0x26, 0x4c, 0xd2,
	0x31, 0x1d, 0x48, 0x32, 0x21, 0x33, 0x56, 0x7d, 0x1d, 0x97, 0x3e, 0x82, 0x4b, 0x97, 0xae, 0x7c,
	0x05, 0x5d, 0xfa, 0x14, 0x32, 0x93, 0x96, 0x44, 0xec, 0xee, 0xde, 0x9c, 0xdf, 0x39, 0x77, 0x6e,
	0xb8, 0xb8, 0x27, 0x49, 0xc0, 0x6c, 0x96, 0x86, 0x34, 0x95, 0x6c, 0x49, 0x85, 0xbd, 0x74, 0xec,
	0x88, 0xa6, 0x54, 0x30, 0x61, 0x65, 0x39, 0x97, 0xdc, 0x34, 0x15, 0x61, 0x95, 0x84, 0xb5, 0x74,
	0xba, 0x47, 0x1b, 0x5c, 0x15, 0x42, 0x1b, 0xbb, 0xbb, 0x11, 0x8f, 0xb8, 0x2e, 0x6d, 0x55, 0x15,
	0x5f, 0xfb, 0x1f, 0x06, 0x6e, 0x4d, 0x8a, 0x01, 0x33, 0x49, 0x24, 0x35, 0x4f, 0x70, 0x3d, 0x23,
	0x39, 0x49, 0x84, 0x09, 0x1d, 0xa3, 0x67, 0x0c, 0x9a, 0xa3, 0xae, 0xf5, 0x77, 0xa4, 0x75, 0xa3,
	0x99, 0x31, 0xb8, 0xe0, 0xad, 0xf9, 0x33, 0x8c, 0x4b, 0xc6, 0x84, 0x0e, 0xea, 0xd5, 0x06, 0xcd,
	0xd1, 0xe1, 0x26, 0xf7, 0x74, 0xdd, 0xe9, 0x80, 0xaa, 0xed, 0x14, 0xe3, 0x88, 0x08, 0x3f, 0xa1,
	0x92, 0xe6, 0x2a, 0xa4, 0xa6, 0x43, 0x0e, 0x36, 0x85, 0x4c, 0x88, 0xb8, 0x56, 0x98, 0xce, 0x68,
	0x44, 0xab, 0x4e, 0xf4, 0x3f, 0x0d, 0x5c, 0x2f, 0x9e, 0x67, 0x3a, 0x78, 0x87, 0xa6, 0x24, 0x88,
	0xa9, 0xff, 0xeb, 0x65, 0x6a, 0xaf, 0xff, 0x23, 0x18, 0x83, 0xd7, 0x2e, 0xe4, 0x69, 0xa9, 0x3a,
	0xb8, 0x4d, 0xe2, 0x98, 0x87, 0x44, 0x32, 0x9e, 0xfa, 0x31, 0x4b, 0x98, 0xd4, 0xbb, 0x18, 0x83,
	0x86, 0x72, 0xb8, 0xe0, 0x6d, 0x97, 0xfa, 0x95, 0x96, 0xcf, 0xf1, 0x7e, 0x19, 0xef, 0xd3, 0x8c,
	0x87, 0x0b, 0x9f, 0xcd, 0x55, 0x7f, 0xc7, 0x68, 0xae, 0x97, 0x58, 0xb9, 0xbd, 0xbd, 0x12, 0xbc,
	0x50, 0xdc, 0xb4, 0xc4, 0x86, 0x78, 0x2b, 0xa7, 0x0f, 0x24, 0x9f, 0xfb, 0x22, 0x24, 0xb1, 0xf6,
	0xfd, 0xab, 0x4c, 0x6d, 0x15, 0xe2, 0x4c, 0x6b, 0xee, 0x25, 0x06, 0x17, 0x6e, 0x87, 0x11, 0x93,
	0x8b, 0xfb, 0xc0, 0x0a, 0x79, 0x62, 0xab, 0x1f, 0x14, 0x93, 0x40, 0xe8, 0xc2, 0x7e, 0xac, 0x5e,
	0x83, 0x7c, 0xca, 0xa8, 0x78, 0x46, 0xf0, 0x82, 0xe0, 0x15, 0xc1, 0x1b, 0x82, 0x77, 0x04, 0x5f,
	0x08, 0xbe, 0x11, 0x04, 0x75, 0x7d, 0x05, 0xc7, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd8, 0xeb,
	0xcf, 0x1a, 0x78, 0x02, 0x00, 0x00,
}

func (m *V2GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *V2GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *V2GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.V2GasMeters) > 0 {
		for iNdEx := len(m.V2GasMeters) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.V2GasMeters[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.V2Incentives) > 0 {
		for iNdEx := len(m.V2Incentives) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.V2Incentives[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.V2Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *V2Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *V2Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *V2Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.RewardScaler.Size()
		i -= size
		if _, err := m.RewardScaler.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.IncentivesEpochIdentifier) > 0 {
		i -= len(m.IncentivesEpochIdentifier)
		copy(dAtA[i:], m.IncentivesEpochIdentifier)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.IncentivesEpochIdentifier)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size := m.AllocationLimit.Size()
		i -= size
		if _, err := m.AllocationLimit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.EnableIncentives {
		i--
		if m.EnableIncentives {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func (m *V2GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.V2Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.V2Incentives) > 0 {
		for _, e := range m.V2Incentives {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.V2GasMeters) > 0 {
		for _, e := range m.V2GasMeters {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *V2Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EnableIncentives {
		n += 2
	}
	l = m.AllocationLimit.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = len(m.IncentivesEpochIdentifier)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.RewardScaler.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func (m *V2GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: V2GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: V2GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V2Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.V2Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V2Incentives", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.V2Incentives = append(m.V2Incentives, V2Incentive{})
			if err := m.V2Incentives[len(m.V2Incentives)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V2GasMeters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.V2GasMeters = append(m.V2GasMeters, V2GasMeter{})
			if err := m.V2GasMeters[len(m.V2GasMeters)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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

func (m *V2Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: V2Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: V2Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableIncentives", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.EnableIncentives = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllocationLimit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AllocationLimit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncentivesEpochIdentifier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IncentivesEpochIdentifier = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardScaler", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardScaler.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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

func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)