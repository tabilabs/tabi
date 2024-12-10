package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the module
	ModuleName = "tokenconvert"

	// StoreKey is the store key for token-convert
	StoreKey = ModuleName

	RouterKey = ModuleName
)

// KVStore keys
var (
	StrategyKey       = []byte{0x01}
	VoucherKey        = []byte{0x02}
	VoucherByOwnerKey = []byte{0x03}
	VoucherSeqKey     = []byte{0x04}

	Delimiter   = []byte{0x00}
	PlaceHolder = []byte{0x01}
)

// StrategyStoreKey returns the byte representation of the strategy key
// Items are stored with key as follows:
// 0x01<strategy>
func StrategyStoreKey(strategy []byte) []byte {
	bz := make([]byte, len(StrategyKey)+len(strategy))

	copy(bz, StrategyKey)
	copy(bz[len(StrategyKey):], strategy)

	return bz
}

// VoucherStoreKey returns the byte representation of the voucher key
// Items are stored with key as follows:
// 0x02<voucherID>
func VoucherStoreKey(voucherID string) []byte {
	bz := make([]byte, len(VoucherKey)+len(voucherID))

	copy(bz, VoucherKey)
	copy(bz[len(VoucherKey):], voucherID)

	return bz
}

// VoucherByOwnerStoreKey returns the byte representation of the voucher by owner key
// Items are stored with key as follows:
// 0x03<owner><Delimiter(1 Byte)><voucherID>
func VoucherByOwnerStoreKey(owner sdktypes.AccAddress, voucherID string) []byte {
	address.MustLengthPrefix(owner)

	bz := make([]byte, len(VoucherByOwnerKey)+len(owner)+len(Delimiter)+len(voucherID))

	copy(bz, VoucherByOwnerKey)
	copy(bz[len(VoucherByOwnerKey):], owner)
	copy(bz[len(VoucherByOwnerKey)+len(owner):], Delimiter)
	copy(bz[len(VoucherByOwnerKey)+len(owner)+len(Delimiter):], voucherID)

	return bz
}

// VoucherByOwnerStorePrefixKey returns the byte representation of the voucher by owner prefix key
// Items are stored with key as follows:
// 0x03<owner><Delimiter(1 Byte)>
func VoucherByOwnerStorePrefixKey(owner sdktypes.AccAddress) []byte {
	address.MustLengthPrefix(owner)

	bz := make([]byte, len(VoucherByOwnerKey)+len(owner)+len(Delimiter))

	copy(bz, VoucherByOwnerKey)
	copy(bz[len(VoucherByOwnerKey):], owner)
	copy(bz[len(VoucherByOwnerKey)+len(owner):], Delimiter)

	return bz
}
