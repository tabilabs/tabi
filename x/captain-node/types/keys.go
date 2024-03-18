package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// nolint
const (
	// ModuleName defines the module name
	ModuleName = "captain-node"

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

	// KeyCallers is the key used to store the callers address in the keeper.
	KeyCallers = "callers"
)

var (
	ParamsKey           = []byte{0x00}
	NodeKey             = []byte{0x01}
	NodeByOwner         = []byte{0x02}
	OwnerKey            = []byte{0x03}
	DivisionKey         = []byte{0x04}
	DivisionTotalSupply = []byte{0x05}
	OwnerHolding        = []byte{0x06}
	PowerOnPeriod       = []byte{0x07}

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

// NodeStoreKey returns the byte representation of the node key
func NodeStoreKey(nodeID string) []byte {
	key := make([]byte, len(NodeKey)+len(nodeID))
	copy(key, NodeKey)
	copy(key[len(NodeKey):], nodeID)
	return key
}

// OwnerStoreKey returns the byte representation of the owner key
func OwnerStoreKey(nodeId string) []byte {

	key := make([]byte, len(OwnerKey)+len(nodeId))
	copy(key, OwnerKey)
	copy(key[len(OwnerKey):], nodeId)
	return key
}

// DivisionStoreKey returns the byte representation of the divisions key
func DivisionStoreKey(divisionsID string) []byte {
	key := make([]byte, len(DivisionKey)+len(divisionsID))
	copy(key, DivisionKey)
	copy(key[len(DivisionKey):], divisionsID)
	return key
}

// DivisionTotalSupplyKey returns the byte representation of the divisions total supply key
func DivisionTotalSupplyKey(divisionsID string) []byte {
	key := make([]byte, len(DivisionTotalSupply)+len(divisionsID))
	copy(key, DivisionTotalSupply)
	copy(key[len(DivisionTotalSupply):], divisionsID)
	return key
}

// NodeByOwnerStoreKey returns the byte representation of the node owner
// Items are stored with the following key: values
// 0x02<owner>
func NodeByOwnerStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(NodeByOwner)+len(owner))
	copy(key, NodeByOwner)
	copy(key[len(NodeByOwner):], owner)
	return key
}

func OwnerHoldingTotalSupplyStoreKey(owner sdk.AccAddress) []byte {
	owner = address.MustLengthPrefix(owner)
	key := make([]byte, len(OwnerHolding)+len(owner))
	copy(key, OwnerHolding)
	copy(key[len(OwnerHolding):], owner)
	return key
}

func NodePowerOnPeriodStoreKey(nodeID string) []byte {
	key := make([]byte, len(PowerOnPeriod)+len(nodeID))
	copy(key, PowerOnPeriod)
	copy(key[len(PowerOnPeriod):], nodeID)
	return key
}
