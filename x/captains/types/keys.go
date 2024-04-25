package types

import (
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
	prefixEpoch
	prefixNodeNextSeq
	prefixNode
	prefixNodeByOwner
	prefixDivision
	prefixDivisionByNode
	prefixHistoricalEmissionSumOnEpoch
	prefixEmissionSumOnEpoch
	prefixNodeHistoricalEmissionOnEpoch
	prefixNodeHistoricalEmissionOnLastClaim
	prefixComputingPowerClaimable
	prefixComputingPowerSumOnEpoch
	prefixNodeComputingPowerOnEpoch
	prefixOwnerPledgeOnEpoch
	prefixPledgeSumOnEpoch
	prefixBatchCount
	prefixDigestOnEpoch
	prefixEndOnEpoch
)

var (
	ParamsKey = []byte{prefixParams}

	EpochKey = []byte{prefixEpoch}

	NodeKey             = []byte{prefixNode}
	NodeNextSequenceKey = []byte{prefixNodeNextSeq}
	NodeByOwnerKey      = []byte{prefixNodeByOwner}

	DivisionKey    = []byte{prefixDivision}
	DivisionByNode = []byte{prefixDivisionByNode}

	HistoricalEmissionSumOnEpochKey      = []byte{prefixHistoricalEmissionSumOnEpoch}
	EmissionSumOnEpochKey                = []byte{prefixEmissionSumOnEpoch}
	NodeHistoricalEmissionOnEpochKey     = []byte{prefixNodeHistoricalEmissionOnEpoch}
	NodeHistoricalEmissionOnLastClaimKey = []byte{prefixNodeHistoricalEmissionOnLastClaim}

	ComputingPowerClaimableKey   = []byte{prefixComputingPowerClaimable}
	ComputingPowerSumOnEpochKey  = []byte{prefixComputingPowerSumOnEpoch}
	NodeComputingPowerOnEpochKey = []byte{prefixNodeComputingPowerOnEpoch}

	OwnerPledgeOnEpochKey = []byte{prefixOwnerPledgeOnEpoch}
	PledgeSumOnEpochKey   = []byte{prefixPledgeSumOnEpoch}

	BatchCountOnEpochKey = []byte{prefixBatchCount}
	DigestOnEpochKey     = []byte{prefixDigestOnEpoch}
	EndOnEpochKey        = []byte{prefixEndOnEpoch}

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
	key := make([]byte, len(NodeByOwnerKey)+len(owner))
	copy(key, NodeByOwnerKey)
	copy(key[len(NodeByOwnerKey):], owner)
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

// HistoricalEmissionSumOnEpochStoreKey returns the byte representation of the historical emission sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <emission>
func HistoricalEmissionSumOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(HistoricalEmissionSumOnEpochKey)+len(epochBz))
	copy(key, HistoricalEmissionSumOnEpochKey)
	copy(key[len(HistoricalEmissionSumOnEpochKey):], epochBz)
	return key
}

// NodeHistoricalEmissionOnEpochStoreKey returns the byte representation of the historical emission by node on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id><delimiter><node_id> -> <emission>
func NodeHistoricalEmissionOnEpochStoreKey(epochID uint64, nodeID string) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(NodeHistoricalEmissionOnEpochKey)+len(epochBz)+len(Delimiter)+len(nodeID))
	copy(key, NodeHistoricalEmissionOnEpochKey)
	copy(key[len(NodeHistoricalEmissionOnEpochKey):], epochBz)
	copy(key[len(NodeHistoricalEmissionOnEpochKey)+len(epochBz):], Delimiter)
	copy(key[len(NodeHistoricalEmissionOnEpochKey)+len(epochBz)+len(Delimiter):], nodeID)
	return key
}

// NodeHistoricalEmissionOnLastClaimStoreKey returns the byte representation of the historical emission last claimed key
// Items are stored with the following key: values
// <prefix_key><node_id> -> <emission>
func NodeHistoricalEmissionOnLastClaimStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeHistoricalEmissionOnLastClaimKey)+len(nodeID))
	copy(key, NodeHistoricalEmissionOnLastClaimKey)
	copy(key[len(NodeHistoricalEmissionOnLastClaimKey):], nodeID)
	return key
}

// NodeClaimableComputingPowerStoreKey returns the byte representation of the computing power claimable key
// Items are stored with the following key: values
// <prefix_key><owner> -> <computing_power>
func NodeClaimableComputingPowerStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, len(ComputingPowerClaimableKey)+len(owner))
	copy(key, ComputingPowerClaimableKey)
	copy(key[len(ComputingPowerClaimableKey):], owner)
	return key
}

// ComputingPowerSumOnEpochStoreKey returns the byte representation of the computing power sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <computing_power>
func ComputingPowerSumOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(ComputingPowerSumOnEpochKey)+len(epochBz))
	copy(key, ComputingPowerSumOnEpochKey)
	copy(key[len(ComputingPowerSumOnEpochKey):], epochBz)
	return key
}

// NodeComputingPowerOnEpochStoreKey returns the byte representation of the computing power by node on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id><delimiter><node_id> -> <computing_power>
func NodeComputingPowerOnEpochStoreKey(epochID uint64, nodeID string) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(NodeComputingPowerOnEpochKey)+len(epochBz)+len(Delimiter)+len(nodeID))
	copy(key, NodeComputingPowerOnEpochKey)
	copy(key[len(NodeComputingPowerOnEpochKey):], epochBz)
	copy(key[len(NodeComputingPowerOnEpochKey)+len(epochBz):], Delimiter)
	copy(key[len(NodeComputingPowerOnEpochKey)+len(epochBz)+len(Delimiter):], nodeID)
	return key
}

// PledgeSumOnEpochStoreKey returns the byte representation of the pledge amount sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <pledge_amount>
func PledgeSumOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(PledgeSumOnEpochKey)+len(epochBz))
	copy(key, PledgeSumOnEpochKey)
	copy(key[len(PledgeSumOnEpochKey):], epochBz)
	return key
}

// EmissionSumOnEpochStoreKey returns the byte representation of the emission sum on epoch key
// Items are stored with the following key: values
// <prefix_key><epoch_id> -> <emission>
func EmissionSumOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(EmissionSumOnEpochKey)+len(epochBz))
	copy(key, EmissionSumOnEpochKey)
	copy(key[len(EmissionSumOnEpochKey):], epochBz)
	return key
}

// OwnerPledgeOnEpochStoreKey returns the byte representation of the owner pledge on epoch key
// <prefix_key><epoch_id><delimiter><node_id> -> <pledge>
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

func BatchCountOnEpochStoreKey(epochID, batchID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	batchBz := sdk.Uint64ToBigEndian(batchID)
	key := make([]byte, len(BatchCountOnEpochKey)+len(epochBz)+len(Delimiter)+len(batchBz))
	copy(key, BatchCountOnEpochKey)
	copy(key[len(BatchCountOnEpochKey):], epochBz)
	copy(key[len(BatchCountOnEpochKey)+len(epochBz):], Delimiter)
	copy(key[len(BatchCountOnEpochKey)+len(epochBz)+len(Delimiter):], epochBz)
	return key
}

func BatchCountOnEpochPrefixKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(BatchCountOnEpochKey)+len(epochBz)+len(Delimiter))
	copy(key, BatchCountOnEpochKey)
	copy(key[len(BatchCountOnEpochKey):], epochBz)
	copy(key[len(BatchCountOnEpochKey)+len(epochBz):], Delimiter)
	return key
}

func DigestOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(DigestOnEpochKey)+len(epochBz))
	copy(key, DigestOnEpochKey)
	copy(key[len(DigestOnEpochKey):], epochBz)
	return key
}

func EndOnEpochStoreKey(epochID uint64) []byte {
	epochBz := sdk.Uint64ToBigEndian(epochID)
	key := make([]byte, len(EmissionSumOnEpochKey)+len(epochBz))
	copy(key, EndOnEpochKey)
	copy(key[len(EndOnEpochKey):], epochBz)
	return key
}
