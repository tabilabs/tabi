package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// nolint
const (
	// ModuleName defines the module name
	ModuleName = "captainnode"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters = "parameters"

	// KeyNextMTSequence is the key used to store the next MT sequence in the keeper
	KeyNextNodeSequence = "nextNodeSequence"
)

var (
	ParamsKey              = []byte{0x00}
	NodeKey                = []byte{0x01}
	NodeByOwnerKey         = []byte{0x02}
	OwnerKey               = []byte{0x03}
	DivisionKey            = []byte{0x04}
	DivisionTotalSupplyKey = []byte{0x05}
	OwnerHoldingKey        = []byte{0x06}
	PowerOnPeriodKey       = []byte{0x07}
	ExperienceKey          = []byte{0x08}

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

// NodeStoreKey returns the byte representation of the node key
// Items are stored with the following key: values
// 0x01<nodeID>
func NodeStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeKey)+len(nodeID))
	copy(key, NodeKey)
	copy(key[len(NodeKey):], nodeID)
	return key
}

// NodeByOwnerStoreKey returns the byte representation of the node owner
// Items are stored with the following key: values
// 0x02<owner><Delimiter>
func NodeByOwnerStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)

	key := make([]byte, len(NodeByOwnerKey)+len(owner)+len(Delimiter))
	copy(key, NodeByOwnerKey)
	copy(key[len(NodeByOwnerKey):], owner)
	copy(key[len(NodeByOwnerKey)+len(owner):], Delimiter)
	return key
}

// OwnerStoreKey returns the byte representation of the owner key
// Items are stored with the following key: values
// 0x03<nodeID>
func OwnerStoreKey(nodeId string) []byte {

	key := make([]byte, len(OwnerKey)+len(nodeId))
	copy(key, OwnerKey)
	copy(key[len(OwnerKey):], nodeId)
	return key
}

// DivisionStoreKey returns the byte representation of the divisions key
// Items are stored with the following key: values
// 0x04<divisionsID>
func DivisionStoreKey(divisionsID string) []byte {
	key := make([]byte, len(DivisionKey)+len(divisionsID))
	copy(key, DivisionKey)
	copy(key[len(DivisionKey):], divisionsID)
	return key
}

// DivisionTotalSupplyStoreKey returns the byte representation of the divisions total supply key
// Items are stored with the following key: values
// 0x05<divisionsID>
func DivisionTotalSupplyStoreKey(divisionsID string) []byte {
	key := make([]byte, len(DivisionTotalSupplyKey)+len(divisionsID))
	copy(key, DivisionTotalSupplyKey)
	copy(key[len(DivisionTotalSupplyKey):], divisionsID)
	return key
}

// OwnerHoldingTotalSupplyStoreKey returns the byte representation of the owner holding total supply key
// Items are stored with the following key: values
// 0x06<owner>
func OwnerHoldingTotalSupplyStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(OwnerHoldingKey)+len(owner))
	copy(key, OwnerHoldingKey)
	copy(key[len(OwnerHoldingKey):], owner)
	return key
}

// NodePowerOnPeriodStoreKey returns the byte representation of the node power on period key
// Items are stored with the following key: values
// 0x07<nodeID>
func NodePowerOnPeriodStoreKey(nodeID string) []byte {
	key := make([]byte, len(PowerOnPeriodKey)+len(nodeID))
	copy(key, PowerOnPeriodKey)
	copy(key[len(PowerOnPeriodKey):], nodeID)
	return key
}

// ExperienceStoreKey returns the byte representation of the experience key
// Items are stored with the following key: values
// 0x08<owner>
func ExperienceStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(ExperienceKey)+len(owner))
	copy(key, ExperienceKey)
	copy(key[len(ExperienceKey):], owner)
	return key
}
