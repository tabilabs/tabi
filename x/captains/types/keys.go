package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// nolint
const (
	// ModuleName defines the module name
	ModuleName = "captains"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// KeyNextMTSequence is the key used to store the next CaptainNode sequence in the keeper
	KeyNextNodeSequence = "nextNodeSequence"
)

var (
	ParamsKey                          = []byte{0x00}
	NodeKey                            = []byte{0x01}
	NodeByOwnerKey                     = []byte{0x02}
	DivisionKey                        = []byte{0x03}
	DivisionByNode                     = []byte{0x04}
	HistoricalEmissionSumOnEpochKey    = []byte{0x05}
	HistoricalEmissionByNodeOnEpochKey = []byte{0x06}
	HistoricalEmissionLastClaimedKey   = []byte{0x07}
	ComputingPowerClaimableKey         = []byte{0x08}
	ComputingPowerSumOnEpochKey        = []byte{0x09}
	ComputingPowerByNodeOnEpochKey     = []byte{0x0A}
	PledgeAmountSumOnEpochKey          = []byte{0x0B}

	ReportValidationSumOnEpochKey = []byte{0x0C}
	ComputingPowerCalcCountKey    = []byte{0x0D}
	RewardCalcCountKey            = []byte{0x0E}

	Delimiter   = []byte{0x00}
	PlaceHolder = []byte{0x01}
)

// NodeStoreKey returns the byte representation of the node key
// Items are stored with the following key: values
// 0x01<node_id> -> <node_info_bz>
func NodeStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeKey)+len(nodeID))
	copy(key, NodeKey)
	copy(key[len(NodeKey):], nodeID)
	return key
}

// NodeByOwnerStoreKey returns the byte representation of the node owner
// Items are stored with the following key: values
// 0x02<owner><delimiter><node_id> -> <place_holder>
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
// 0x02<owner><delimiter>
func NodeByOwnerPrefixStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(NodeByOwnerKey)+len(owner))
	copy(key, NodeByOwnerKey)
	copy(key[len(NodeByOwnerKey):], owner)
	return key
}

// DivisionStoreKey returns the byte representation of the divisions key
// Items are stored with the following key: values
// 0x03<division_id> -> <divisions_info_bz>
func DivisionStoreKey(divisionID string) []byte {
	key := make([]byte, len(DivisionKey)+len(divisionID))
	copy(key, DivisionKey)
	copy(key[len(DivisionKey):], divisionID)
	return key
}

// DivisionByNodeStoreKey returns the byte representation of the division by node key
// Items are stored with the following key: values
// 0x04<division_id><delimiter><node_id> -> <place_holder>
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
// 0x05<epoch_id> -> <emission>
func HistoricalEmissionSumOnEpochStoreKey(epochID string) []byte {
	key := make([]byte, len(HistoricalEmissionSumOnEpochKey)+len(epochID))
	copy(key, HistoricalEmissionSumOnEpochKey)
	copy(key[len(HistoricalEmissionSumOnEpochKey):], epochID)
	return key
}

// HistoricalEmissionByNodeOnEpochStoreKey returns the byte representation of the historical emission by node on epoch key
// Items are stored with the following key: values
// 0x06<epoch_id><delimiter><node_id> -> <emission>
func HistoricalEmissionByNodeOnEpochStoreKey(epochID, nodeID string) []byte {
	key := make([]byte, len(HistoricalEmissionByNodeOnEpochKey)+len(epochID)+len(Delimiter)+len(nodeID))
	copy(key, HistoricalEmissionByNodeOnEpochKey)
	copy(key[len(HistoricalEmissionByNodeOnEpochKey):], epochID)
	copy(key[len(HistoricalEmissionByNodeOnEpochKey)+len(epochID):], Delimiter)
	copy(key[len(HistoricalEmissionByNodeOnEpochKey)+len(epochID)+len(Delimiter):], nodeID)
	return key
}

// HistoricalEmissionLastClaimedStoreKey returns the byte representation of the historical emission last claimed key
// Items are stored with the following key: values
// 0x07<node_id> -> <emission>
func HistoricalEmissionLastClaimedStoreKey(nodeID string) []byte {
	key := make([]byte, len(HistoricalEmissionLastClaimedKey)+len(nodeID))
	copy(key, HistoricalEmissionLastClaimedKey)
	copy(key[len(HistoricalEmissionLastClaimedKey):], nodeID)
	return key
}

// ComputingPowerClaimableStoreKey returns the byte representation of the computing power claimable key
// Items are stored with the following key: values
// 0x08<owner> -> <computing_power>
func ComputingPowerClaimableStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, len(ComputingPowerClaimableKey)+len(owner))
	copy(key, ComputingPowerClaimableKey)
	copy(key[len(ComputingPowerClaimableKey):], owner)
	return key
}

// ComputingPowerSumOnEpochStoreKey returns the byte representation of the computing power sum on epoch key
// Items are stored with the following key: values
// 0x09<epoch_id> -> <computing_power>
func ComputingPowerSumOnEpochStoreKey(epochID string) []byte {
	key := make([]byte, len(ComputingPowerSumOnEpochKey)+len(epochID))
	copy(key, ComputingPowerSumOnEpochKey)
	copy(key[len(ComputingPowerSumOnEpochKey):], epochID)
	return key
}

// ComputingPowerByNodeOnEpochStoreKey returns the byte representation of the computing power by node on epoch key
// Items are stored with the following key: values
// 0x0A<epoch_id><delimiter><node_id> -> <computing_power>
func ComputingPowerByNodeOnEpochStoreKey(epochID, nodeID string) []byte {
	key := make([]byte, len(ComputingPowerByNodeOnEpochKey)+len(epochID)+len(Delimiter)+len(nodeID))
	copy(key, ComputingPowerByNodeOnEpochKey)
	copy(key[len(ComputingPowerByNodeOnEpochKey):], epochID)
	copy(key[len(ComputingPowerByNodeOnEpochKey)+len(epochID):], Delimiter)
	copy(key[len(ComputingPowerByNodeOnEpochKey)+len(epochID)+len(Delimiter):], nodeID)
	return key
}

// PledgeAmountSumOnEpochStoreKey returns the byte representation of the pledge amount sum on epoch key
// Items are stored with the following key: values
// 0x0B<epoch_id> -> <pledge_amount>
func PledgeAmountSumOnEpochStoreKey(epochID string) []byte {
	key := make([]byte, len(PledgeAmountSumOnEpochKey)+len(epochID))
	copy(key, PledgeAmountSumOnEpochKey)
	copy(key[len(PledgeAmountSumOnEpochKey):], epochID)
	return key
}

// ReportValidationSumOnEpochStoreKey returns the byte representation of the report validation sum on epoch key
// Items are stored with the following key: values
// 0x0C<epoch_id> -> <report_validation>
func ReportValidationSumOnEpochStoreKey(epochID string) []byte {
	key := make([]byte, len(ReportValidationSumOnEpochKey)+len(epochID))
	copy(key, ReportValidationSumOnEpochKey)
	copy(key[len(ReportValidationSumOnEpochKey):], epochID)
	return key
}

// ComputingPowerCalcCountStoreKey returns the byte representation of the computing power calculation count key
// Items are stored with the following key: values
// 0x0D<epoch_id> -> <computing_power_calc_count>
func ComputingPowerCalcCountStoreKey(epochID string) []byte {
	key := make([]byte, len(ComputingPowerCalcCountKey)+len(epochID))
	copy(key, ComputingPowerCalcCountKey)
	copy(key[len(ComputingPowerCalcCountKey):], epochID)
	return key
}

// RewardCalcCountStoreKey returns the byte representation of the reward calculation count key
// Items are stored with the following key: values
// 0x0E<epoch_id> -> <reward_calc_count>
func RewardCalcCountStoreKey(epochID string) []byte {
	key := make([]byte, len(RewardCalcCountKey)+len(epochID))
	copy(key, RewardCalcCountKey)
	copy(key[len(RewardCalcCountKey):], epochID)
	return key
}
