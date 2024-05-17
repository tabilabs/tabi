package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "captains"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

const (
	prefixParams = iota + 1
	prefixNodeNextSeq
	prefixNode
	prefixNodeByOwner
	prefixDivision
	prefixDivisionByNode
	prefixEpochEmission
	prefixGlobalClaimedEmission
	prefixNodeClaimedEmission
	prefixNodeCumulativeEmissionByEpoch
	prefixClaimableComputingPower
	prefixGlobalComputingPowerOnEpoch
	prefixNodeComputingPowerOnEpoch
	prefixGlobalPledgeOnEpoch
	prefixOwnerPledgeOnEpoch
	prefixCurrEpoch
	prefixReportDigestOnEpoch
	prefixReportBatchOnEpoch
	prefixEndOnEpoch
	prefixStandByOver
)

var (
	ParamsKey                        = []byte{prefixParams}
	CurrEpochKey                     = []byte{prefixCurrEpoch}
	NodeKey                          = []byte{prefixNode}
	NodeNextSequenceKey              = []byte{prefixNodeNextSeq}
	NodeByOwnerKey                   = []byte{prefixNodeByOwner}
	DivisionKey                      = []byte{prefixDivision}
	DivisionByNode                   = []byte{prefixDivisionByNode}
	EpochEmissionKey                 = []byte{prefixEpochEmission}
	GlobalClaimedEmissionKey         = []byte{prefixGlobalClaimedEmission}
	NodeClaimedEmissionKey           = []byte{prefixNodeClaimedEmission}
	NodeCumulativeEmissionByEpochKey = []byte{prefixNodeCumulativeEmissionByEpoch}
	ClaimableComputingPowerKey       = []byte{prefixClaimableComputingPower}
	GlobalComputingPowerOnEpochKey   = []byte{prefixGlobalComputingPowerOnEpoch}
	NodeComputingPowerOnEpochKey     = []byte{prefixNodeComputingPowerOnEpoch}
	GlobalPledgeOnEpochKey           = []byte{prefixGlobalPledgeOnEpoch}
	OwnerPledgeOnEpochKey            = []byte{prefixOwnerPledgeOnEpoch}
	ReportDigestOnEpochKey           = []byte{prefixReportDigestOnEpoch}
	ReportBatchOnEpochKey            = []byte{prefixReportBatchOnEpoch}
	EndOnEpochKey                    = []byte{prefixEndOnEpoch}
	StandByOverKey                   = []byte{prefixStandByOver}

	Delimiter   = []byte{0x00}
	PlaceHolder = []byte{0x01}
)

// NodeStoreKey returns the byte representation of the node key
// Items are stored with the following key: values
// <prefix_key><node_id> -> <node_info_bz>
func NodeStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeKey)+len(nodeID))
	copy(key, NodeKey)
	copy(key[len(NodeKey):], nodeID)
	return key
}

// NodeByOwnerStoreKey returns the byte representation of the node owner
// Items are stored with the following key: values
// <prefix_key><owner><delimiter><node_id> -> <place_holder>
func NodeByOwnerStoreKey(owner sdk.AccAddress, nodeID string) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(NodeByOwnerKey)+len(owner)+len(Delimiter)+len(nodeID))
	copy(key, NodeByOwnerKey)
	copy(key[len(NodeByOwnerKey):], owner)
	copy(key[len(NodeByOwnerKey)+len(owner):], Delimiter)
	copy(key[len(NodeByOwnerKey)+len(owner)+len(Delimiter):], nodeID)
	return key
}

// NodeByOwnerPrefixStoreKey returns the byte representation of the node by owner prefix key
// Items are stored with the following key
// <prefix_key><owner><delimiter>
func NodeByOwnerPrefixStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(NodeByOwnerKey)+len(owner)+len(Delimiter))
	copy(key, NodeByOwnerKey)
	copy(key[len(NodeByOwnerKey):], owner)
	copy(key[len(NodeByOwnerKey)+len(owner):], Delimiter)
	return key
}

// DivisionStoreKey returns the byte representation of the divisions key
// Items are stored with the following key: values
// <prefix_key><division_id> -> <divisions_info_bz>
func DivisionStoreKey(divisionID string) []byte {
	key := make([]byte, len(DivisionKey)+len(divisionID))
	copy(key, DivisionKey)
	copy(key[len(DivisionKey):], divisionID)
	return key
}

// DivisionByNodeStoreKey returns the byte representation of the division by node key
// Items are stored with the following key: values
// <prefix_key><division_id><delimiter><node_id> -> <place_holder>
func DivisionByNodeStoreKey(divisionID, nodeID string) []byte {
	key := make([]byte, len(DivisionByNode)+len(divisionID)+len(Delimiter)+len(nodeID))
	copy(key, DivisionByNode)
	copy(key[len(DivisionByNode):], nodeID)
	copy(key[len(DivisionByNode)+len(nodeID):], Delimiter)
	copy(key[len(DivisionByNode)+len(nodeID)+len(Delimiter):], divisionID)
	return key
}

// EpochEmissionStoreKey returns the byte representation of the emission sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <emission>
func EpochEmissionStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(EpochEmissionKey)+len(epochBz))
	copy(key, EpochEmissionKey)
	copy(key[len(EpochEmissionKey):], epochBz)
	return key
}

// NodeCumulativeEmissionByEpochStoreKey returns the byte representation of cumulative emission by node on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id><delimiter><node_id> -> <emission>
func NodeCumulativeEmissionByEpochStoreKey(epochID uint64, nodeID string) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(NodeCumulativeEmissionByEpochKey)+len(epochBz)+len(Delimiter)+len(nodeID))
	copy(key, NodeCumulativeEmissionByEpochKey)
	copy(key[len(NodeCumulativeEmissionByEpochKey):], epochBz)
	copy(key[len(NodeCumulativeEmissionByEpochKey)+len(epochBz):], Delimiter)
	copy(key[len(NodeCumulativeEmissionByEpochKey)+len(epochBz)+len(Delimiter):], nodeID)
	return key
}

// ParseNodeCumulativeEmissionByEpochPrefixStoreKey returns the epochID and nodeID from the key
// <epoch_id><delimiter><node_id> are split by delimiter
func ParseNodeCumulativeEmissionByEpochPrefixStoreKey(key []byte) (uint64, string) {
	splits := bytes.Split(key, Delimiter)
	if len(splits) != 2 {
		panic(fmt.Errorf("unexpected key: %s", string(key)))
	}
	epochID := sdk.BigEndianToUint64(splits[0])
	nodeID := string(splits[1])
	return epochID, nodeID
}

// NodeClaimedEmissionStoreKey returns the byte representation of the historical emission last claimed key
// Items are stored with the following key: values
// <prefix_key><node_id> -> <emission>
func NodeClaimedEmissionStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeClaimedEmissionKey)+len(nodeID))
	copy(key, NodeClaimedEmissionKey)
	copy(key[len(NodeClaimedEmissionKey):], nodeID)
	return key
}

// ClaimableComputingPowerStoreKey returns the byte representation of claimable computing power key
// Items are stored with the following key: values
// <prefix_key><owner> -> <computing_power>
func ClaimableComputingPowerStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(ClaimableComputingPowerKey)+len(owner))
	copy(key, ClaimableComputingPowerKey)
	copy(key[len(ClaimableComputingPowerKey):], owner)
	return key
}

// GlobalComputingPowerOnEpochStoreKey returns the byte representation of the computing power sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <computing_power>
func GlobalComputingPowerOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(GlobalComputingPowerOnEpochKey)+len(epochBz))
	copy(key, GlobalComputingPowerOnEpochKey)
	copy(key[len(GlobalComputingPowerOnEpochKey):], epochBz)
	return key
}

// NodeComputingPowerOnEpochStoreKey returns the byte representation of the computing power by node on epoch key
// Items are stored with the following key: values
// <prefix_key><node_id><delimiter><epoch_id> -> <computing_power>
func NodeComputingPowerOnEpochStoreKey(epochID uint64, nodeID string) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(NodeComputingPowerOnEpochKey)+len(nodeID)+len(Delimiter)+len(epochBz))
	copy(key, NodeComputingPowerOnEpochKey)
	copy(key[len(NodeComputingPowerOnEpochKey):], nodeID)
	copy(key[len(NodeComputingPowerOnEpochKey)+len(nodeID):], Delimiter)
	copy(key[len(NodeComputingPowerOnEpochKey)+len(nodeID)+len(Delimiter):], epochBz)
	return key
}

// ParseNodeComputingPowerOnEpochPrefixStoreKey return the epoch id and node id of the computing power on epoch key.
func ParseNodeComputingPowerOnEpochPrefixStoreKey(key []byte) (uint64, string) {
	splits := bytes.Split(key, Delimiter)
	if len(splits) != 2 {
		panic(fmt.Errorf("unexpected key: %s", string(key)))
	}
	nodeID := string(splits[0])
	epochID := sdk.BigEndianToUint64(splits[1])
	return epochID, nodeID
}

// GlobalPledgeOnEpochStoreKey returns the byte representation of the pledge amount sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <pledge_amount>
func GlobalPledgeOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(GlobalPledgeOnEpochKey)+len(epochBz))
	copy(key, GlobalPledgeOnEpochKey)
	copy(key[len(GlobalPledgeOnEpochKey):], epochBz)
	return key
}

// OwnerPledgeOnEpochStoreKey returns the byte representation of the owner pledge on epoch key
// <prefix_key><epoch_id><delimiter><owner> -> <pledge>
func OwnerPledgeOnEpochStoreKey(owner sdk.AccAddress, epochID uint64) []byte {
	owner = address.MustLengthPrefix(owner)
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(OwnerPledgeOnEpochKey)+len(owner)+len(epochBz))
	copy(key, OwnerPledgeOnEpochKey)
	copy(key[len(OwnerPledgeOnEpochKey):], owner)
	copy(key[len(OwnerPledgeOnEpochKey)+len(owner):], Delimiter)
	copy(key[len(OwnerPledgeOnEpochKey)+len(owner)+len(Delimiter):], epochBz)
	return key
}

// ParseOwnerPledgeOnEpochPrefixStoreKey returns the epoch id and node id of the owner pledge on epoch key
func ParseOwnerPledgeOnEpochPrefixStoreKey(key []byte) (uint64, string) {
	splits := bytes.Split(key, Delimiter)
	if len(splits) != 2 {
		panic(fmt.Errorf("unexpected key: %s", string(key)))
	}
	epochID := sdk.BigEndianToUint64(splits[0])
	owner := string(splits[1])
	return epochID, owner
}

// ReportDigestOnEpochStoreKey returns the byte representation of the digest on epoch key
// <prefix_key><epoch_id> -> <digest>
func ReportDigestOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(ReportDigestOnEpochKey)+len(epochBz))
	copy(key, ReportDigestOnEpochKey)
	copy(key[len(ReportDigestOnEpochKey):], epochBz)
	return key
}

// ReportBatchOnEpochStoreKey returns the byte representation of the report batch on epoch key
// <prefix_key><epoch_id><delimiter><batch_id> -> <node_count_this_batch>
func ReportBatchOnEpochStoreKey(epochID, batchID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	batchBz := sdk.Uint64ToBigEndian(batchID)
	key := make([]byte, len(ReportBatchOnEpochKey)+len(epochBz)+len(Delimiter)+len(batchBz))
	copy(key, ReportBatchOnEpochKey)
	copy(key[len(ReportBatchOnEpochKey):], epochBz)
	copy(key[len(ReportBatchOnEpochKey)+len(epochBz):], Delimiter)
	copy(key[len(ReportBatchOnEpochKey)+len(epochBz)+len(Delimiter):], batchBz)
	return key
}

// ReportBatchOnEpochPrefixStoreKey returns the byte representation of the report batch on epoch prefix key
// <prefix_key><epoch_id><delimiter>
func ReportBatchOnEpochPrefixStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(ReportBatchOnEpochKey)+len(epochBz)+len(Delimiter))
	copy(key, ReportBatchOnEpochKey)
	copy(key[len(ReportBatchOnEpochKey):], epochBz)
	copy(key[len(ReportBatchOnEpochKey)+len(epochBz):], Delimiter)
	return key
}

// EndOnEpochStoreKey returns the byte representation of the end on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <end>
func EndOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(EndOnEpochKey)+len(epochBz))
	copy(key, EndOnEpochKey)
	copy(key[len(EndOnEpochKey):], epochBz)
	return key
}
