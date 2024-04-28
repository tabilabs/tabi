// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tabi/captains/v1/report.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type ReportType int32

const (
	ReportType_REPORT_TYPE_UNSPECIFIED ReportType = 0
	ReportType_REPORT_TYPE_DIGEST      ReportType = 1
	ReportType_REPORT_TYPE_BATCH       ReportType = 2
	ReportType_REPORT_TYPE_END         ReportType = 3
)

var ReportType_name = map[int32]string{
	0: "REPORT_TYPE_UNSPECIFIED",
	1: "REPORT_TYPE_DIGEST",
	2: "REPORT_TYPE_BATCH",
	3: "REPORT_TYPE_END",
}

var ReportType_value = map[string]int32{
	"REPORT_TYPE_UNSPECIFIED": 0,
	"REPORT_TYPE_DIGEST":      1,
	"REPORT_TYPE_BATCH":       2,
	"REPORT_TYPE_END":         3,
}

func (x ReportType) String() string {
	return proto.EnumName(ReportType_name, int32(x))
}

func (ReportType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2b04da73fb1305c0, []int{0}
}

type ReportDigest struct {
	// epoch_id is the epoch id of the report
	EpochId uint64 `protobuf:"varint,1,opt,name=epoch_id,json=epochId,proto3" json:"epoch_id,omitempty"`
	// total_node_count is the total number of batches in the report
	TotalBatchCount uint64 `protobuf:"varint,2,opt,name=total_batch_count,json=totalBatchCount,proto3" json:"total_batch_count,omitempty"`
	// total_node_count is the total number of nodes in the report
	TotalNodeCount uint64 `protobuf:"varint,3,opt,name=total_node_count,json=totalNodeCount,proto3" json:"total_node_count,omitempty"`
	// maximum_node_count_per_batch is the maximum number of nodes per batch
	MaximumNodeCountPerBatch uint64 `protobuf:"varint,4,opt,name=maximum_node_count_per_batch,json=maximumNodeCountPerBatch,proto3" json:"maximum_node_count_per_batch,omitempty"`
	// global_on_operation_ratio is the operation ratio of global nodes
	GlobalOnOperationRatio github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=global_on_operation_ratio,json=globalOnOperationRatio,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"global_on_operation_ratio"`
}

func (m *ReportDigest) Reset()         { *m = ReportDigest{} }
func (m *ReportDigest) String() string { return proto.CompactTextString(m) }
func (*ReportDigest) ProtoMessage()    {}
func (*ReportDigest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b04da73fb1305c0, []int{0}
}
func (m *ReportDigest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReportDigest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReportDigest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReportDigest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportDigest.Merge(m, src)
}
func (m *ReportDigest) XXX_Size() int {
	return m.Size()
}
func (m *ReportDigest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportDigest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportDigest proto.InternalMessageInfo

func (m *ReportDigest) GetEpochId() uint64 {
	if m != nil {
		return m.EpochId
	}
	return 0
}

func (m *ReportDigest) GetTotalBatchCount() uint64 {
	if m != nil {
		return m.TotalBatchCount
	}
	return 0
}

func (m *ReportDigest) GetTotalNodeCount() uint64 {
	if m != nil {
		return m.TotalNodeCount
	}
	return 0
}

func (m *ReportDigest) GetMaximumNodeCountPerBatch() uint64 {
	if m != nil {
		return m.MaximumNodeCountPerBatch
	}
	return 0
}

// ReportBatch marks the a batch of nodes.
type ReportBatch struct {
	// epoch_id is the epoch id of the report
	EpochId uint64 `protobuf:"varint,1,opt,name=epoch_id,json=epochId,proto3" json:"epoch_id,omitempty"`
	// batch_id is the batch id of the report
	BatchId uint64 `protobuf:"varint,2,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	// node_count is the number of nodes in the batch
	NodeCount uint64 `protobuf:"varint,3,opt,name=node_count,json=nodeCount,proto3" json:"node_count,omitempty"`
	// node_ids is the list of node ids in the batch
	NodeIds []string `protobuf:"bytes,4,rep,name=node_ids,json=nodeIds,proto3" json:"node_ids,omitempty"`
}

func (m *ReportBatch) Reset()         { *m = ReportBatch{} }
func (m *ReportBatch) String() string { return proto.CompactTextString(m) }
func (*ReportBatch) ProtoMessage()    {}
func (*ReportBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b04da73fb1305c0, []int{1}
}
func (m *ReportBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReportBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReportBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReportBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportBatch.Merge(m, src)
}
func (m *ReportBatch) XXX_Size() int {
	return m.Size()
}
func (m *ReportBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportBatch.DiscardUnknown(m)
}

var xxx_messageInfo_ReportBatch proto.InternalMessageInfo

func (m *ReportBatch) GetEpochId() uint64 {
	if m != nil {
		return m.EpochId
	}
	return 0
}

func (m *ReportBatch) GetBatchId() uint64 {
	if m != nil {
		return m.BatchId
	}
	return 0
}

func (m *ReportBatch) GetNodeCount() uint64 {
	if m != nil {
		return m.NodeCount
	}
	return 0
}

func (m *ReportBatch) GetNodeIds() []string {
	if m != nil {
		return m.NodeIds
	}
	return nil
}

// ReportEnd marks the end of commiting a report.
type ReportEnd struct {
	Epoch uint64 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (m *ReportEnd) Reset()         { *m = ReportEnd{} }
func (m *ReportEnd) String() string { return proto.CompactTextString(m) }
func (*ReportEnd) ProtoMessage()    {}
func (*ReportEnd) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b04da73fb1305c0, []int{2}
}
func (m *ReportEnd) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReportEnd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReportEnd.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReportEnd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportEnd.Merge(m, src)
}
func (m *ReportEnd) XXX_Size() int {
	return m.Size()
}
func (m *ReportEnd) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportEnd.DiscardUnknown(m)
}

var xxx_messageInfo_ReportEnd proto.InternalMessageInfo

func (m *ReportEnd) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func init() {
	proto.RegisterEnum("tabi.captains.v1.ReportType", ReportType_name, ReportType_value)
	proto.RegisterType((*ReportDigest)(nil), "tabi.captains.v1.ReportDigest")
	proto.RegisterType((*ReportBatch)(nil), "tabi.captains.v1.ReportBatch")
	proto.RegisterType((*ReportEnd)(nil), "tabi.captains.v1.ReportEnd")
}

func init() { proto.RegisterFile("tabi/captains/v1/report.proto", fileDescriptor_2b04da73fb1305c0) }

var fileDescriptor_2b04da73fb1305c0 = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xe3, 0x24, 0x25, 0xc9, 0x80, 0xa8, 0xbb, 0x94, 0xe2, 0x14, 0xea, 0x86, 0x1c, 0x50,
	0xa8, 0x54, 0x5b, 0x15, 0x57, 0x84, 0x44, 0x62, 0x03, 0xbe, 0x24, 0x91, 0x6b, 0x0e, 0x70, 0x59,
	0xd9, 0xde, 0x95, 0x63, 0x11, 0x7b, 0x2d, 0x7b, 0x53, 0xda, 0x03, 0xef, 0xc0, 0xc3, 0x70, 0xe0,
	0x11, 0x7a, 0xac, 0x38, 0x21, 0x0e, 0x15, 0x4a, 0x5e, 0x04, 0x79, 0xd7, 0x2d, 0x96, 0x90, 0xb8,
	0x78, 0x77, 0xe6, 0xff, 0x77, 0x7e, 0xeb, 0xd3, 0xc0, 0x01, 0xf7, 0x83, 0xd8, 0x0c, 0xfd, 0x8c,
	0xfb, 0x71, 0x5a, 0x98, 0x67, 0x27, 0x66, 0x4e, 0x33, 0x96, 0x73, 0x23, 0xcb, 0x19, 0x67, 0x48,
	0x2d, 0x65, 0xe3, 0x46, 0x36, 0xce, 0x4e, 0xf6, 0x77, 0x23, 0x16, 0x31, 0x21, 0x9a, 0xe5, 0x4d,
	0xfa, 0xf6, 0xfb, 0x21, 0x2b, 0x12, 0x56, 0x60, 0x29, 0xc8, 0x42, 0x4a, 0xc3, 0xef, 0x4d, 0xb8,
	0xe7, 0x8a, 0x99, 0x56, 0x1c, 0xd1, 0x82, 0xa3, 0x3e, 0x74, 0x69, 0xc6, 0xc2, 0x05, 0x8e, 0x89,
	0xa6, 0x0c, 0x94, 0x51, 0xdb, 0xed, 0x88, 0xda, 0x21, 0xe8, 0x08, 0x76, 0x38, 0xe3, 0xfe, 0x12,
	0x07, 0x3e, 0x0f, 0x17, 0x38, 0x64, 0xab, 0x94, 0x6b, 0x4d, 0xe1, 0xd9, 0x16, 0xc2, 0xb8, 0xec,
	0x4f, 0xca, 0x36, 0x1a, 0x81, 0x2a, 0xbd, 0x29, 0x23, 0xb4, 0xb2, 0xb6, 0x84, 0xf5, 0xbe, 0xe8,
	0x4f, 0x19, 0xa1, 0xd2, 0xf9, 0x0a, 0x9e, 0x24, 0xfe, 0x79, 0x9c, 0xac, 0x92, 0x9a, 0x17, 0x67,
	0x34, 0x97, 0x31, 0x5a, 0x5b, 0xbc, 0xd2, 0x2a, 0xcf, 0xed, 0xbb, 0x39, 0xcd, 0x45, 0x1c, 0xfa,
	0x0c, 0xfd, 0x68, 0xc9, 0x02, 0x7f, 0x89, 0x59, 0x8a, 0x59, 0x46, 0x73, 0x9f, 0xc7, 0x2c, 0xc5,
	0xe2, 0xd0, 0xb6, 0x06, 0xca, 0xa8, 0x37, 0x7e, 0x79, 0x79, 0x7d, 0xd8, 0xf8, 0x75, 0x7d, 0xf8,
	0x2c, 0x8a, 0xf9, 0x62, 0x15, 0x18, 0x21, 0x4b, 0x2a, 0x0a, 0xd5, 0x71, 0x5c, 0x90, 0x4f, 0x26,
	0xbf, 0xc8, 0x68, 0x61, 0x58, 0x34, 0xfc, 0xf1, 0xed, 0x18, 0x2a, 0x48, 0x16, 0x0d, 0xdd, 0x3d,
	0x39, 0x7e, 0x96, 0xce, 0x6e, 0x86, 0xbb, 0xe5, 0x77, 0xf8, 0x05, 0xee, 0x4a, 0x72, 0xf2, 0x3f,
	0xfe, 0x03, 0xae, 0x0f, 0x5d, 0x89, 0x2c, 0x26, 0x15, 0xaf, 0x8e, 0xa8, 0x1d, 0x82, 0x0e, 0x00,
	0xfe, 0x21, 0xd4, 0x4b, 0x6f, 0xe1, 0xf4, 0xa1, 0x2b, 0xe4, 0x98, 0x14, 0x5a, 0x7b, 0xd0, 0x1a,
	0xf5, 0xdc, 0x4e, 0x59, 0x3b, 0xa4, 0x18, 0x3e, 0x85, 0x9e, 0x8c, 0xb7, 0x53, 0x82, 0x76, 0x61,
	0x4b, 0x84, 0x55, 0xc9, 0xb2, 0x38, 0x4a, 0x00, 0xa4, 0xc5, 0xbb, 0xc8, 0x28, 0x7a, 0x0c, 0x8f,
	0x5c, 0x7b, 0x3e, 0x73, 0x3d, 0xec, 0x7d, 0x98, 0xdb, 0xf8, 0xfd, 0xf4, 0x74, 0x6e, 0x4f, 0x9c,
	0x37, 0x8e, 0x6d, 0xa9, 0x0d, 0xb4, 0x07, 0xa8, 0x2e, 0x5a, 0xce, 0x5b, 0xfb, 0xd4, 0x53, 0x15,
	0xf4, 0x10, 0x76, 0xea, 0xfd, 0xf1, 0x6b, 0x6f, 0xf2, 0x4e, 0x6d, 0xa2, 0x07, 0xb0, 0x5d, 0x6f,
	0xdb, 0x53, 0x4b, 0x6d, 0x8d, 0x27, 0x97, 0x6b, 0x5d, 0xb9, 0x5a, 0xeb, 0xca, 0xef, 0xb5, 0xae,
	0x7c, 0xdd, 0xe8, 0x8d, 0xab, 0x8d, 0xde, 0xf8, 0xb9, 0xd1, 0x1b, 0x1f, 0x9f, 0xd7, 0xc0, 0x97,
	0x3b, 0xbb, 0xf4, 0x83, 0x42, 0x5c, 0xcc, 0xf3, 0xbf, 0xdb, 0x2d, 0xf8, 0x07, 0x77, 0xc4, 0x5e,
	0xbe, 0xf8, 0x13, 0x00, 0x00, 0xff, 0xff, 0xe7, 0x0d, 0x0e, 0x5c, 0xfb, 0x02, 0x00, 0x00,
}

func (m *ReportDigest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportDigest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReportDigest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.GlobalOnOperationRatio.Size()
		i -= size
		if _, err := m.GlobalOnOperationRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintReport(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.MaximumNodeCountPerBatch != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.MaximumNodeCountPerBatch))
		i--
		dAtA[i] = 0x20
	}
	if m.TotalNodeCount != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.TotalNodeCount))
		i--
		dAtA[i] = 0x18
	}
	if m.TotalBatchCount != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.TotalBatchCount))
		i--
		dAtA[i] = 0x10
	}
	if m.EpochId != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.EpochId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ReportBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReportBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NodeIds) > 0 {
		for iNdEx := len(m.NodeIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.NodeIds[iNdEx])
			copy(dAtA[i:], m.NodeIds[iNdEx])
			i = encodeVarintReport(dAtA, i, uint64(len(m.NodeIds[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if m.NodeCount != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.NodeCount))
		i--
		dAtA[i] = 0x18
	}
	if m.BatchId != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.BatchId))
		i--
		dAtA[i] = 0x10
	}
	if m.EpochId != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.EpochId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ReportEnd) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportEnd) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReportEnd) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Epoch != 0 {
		i = encodeVarintReport(dAtA, i, uint64(m.Epoch))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintReport(dAtA []byte, offset int, v uint64) int {
	offset -= sovReport(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ReportDigest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EpochId != 0 {
		n += 1 + sovReport(uint64(m.EpochId))
	}
	if m.TotalBatchCount != 0 {
		n += 1 + sovReport(uint64(m.TotalBatchCount))
	}
	if m.TotalNodeCount != 0 {
		n += 1 + sovReport(uint64(m.TotalNodeCount))
	}
	if m.MaximumNodeCountPerBatch != 0 {
		n += 1 + sovReport(uint64(m.MaximumNodeCountPerBatch))
	}
	l = m.GlobalOnOperationRatio.Size()
	n += 1 + l + sovReport(uint64(l))
	return n
}

func (m *ReportBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EpochId != 0 {
		n += 1 + sovReport(uint64(m.EpochId))
	}
	if m.BatchId != 0 {
		n += 1 + sovReport(uint64(m.BatchId))
	}
	if m.NodeCount != 0 {
		n += 1 + sovReport(uint64(m.NodeCount))
	}
	if len(m.NodeIds) > 0 {
		for _, s := range m.NodeIds {
			l = len(s)
			n += 1 + l + sovReport(uint64(l))
		}
	}
	return n
}

func (m *ReportEnd) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Epoch != 0 {
		n += 1 + sovReport(uint64(m.Epoch))
	}
	return n
}

func sovReport(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozReport(x uint64) (n int) {
	return sovReport(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ReportDigest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReport
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
			return fmt.Errorf("proto: ReportDigest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportDigest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochId", wireType)
			}
			m.EpochId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBatchCount", wireType)
			}
			m.TotalBatchCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalBatchCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalNodeCount", wireType)
			}
			m.TotalNodeCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalNodeCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaximumNodeCountPerBatch", wireType)
			}
			m.MaximumNodeCountPerBatch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaximumNodeCountPerBatch |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalOnOperationRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
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
				return ErrInvalidLengthReport
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReport
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GlobalOnOperationRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReport(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReport
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
func (m *ReportBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReport
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
			return fmt.Errorf("proto: ReportBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochId", wireType)
			}
			m.EpochId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchId", wireType)
			}
			m.BatchId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BatchId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeCount", wireType)
			}
			m.NodeCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NodeCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeIds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
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
				return ErrInvalidLengthReport
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReport
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeIds = append(m.NodeIds, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipReport(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReport
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
func (m *ReportEnd) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReport
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
			return fmt.Errorf("proto: ReportEnd: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportEnd: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
			}
			m.Epoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Epoch |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipReport(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReport
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
func skipReport(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReport
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
					return 0, ErrIntOverflowReport
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
					return 0, ErrIntOverflowReport
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
				return 0, ErrInvalidLengthReport
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupReport
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthReport
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthReport        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReport          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupReport = fmt.Errorf("proto: unexpected end of group")
)
