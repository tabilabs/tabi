package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the module
	ModuleName = "token-convert"

	// StoreKey is the store key for token-convert
	StoreKey = ModuleName
)

// KVStore keys
var (
	StrategySequenceKey = []byte{0x01}
	VoucherSequenceKey  = []byte{0x02}
	StrategyKey         = []byte{0x03}
	VoucherKey          = []byte{0x04}
	VoucherByOwnerKey   = []byte{0x05}

	Delimiter   = []byte{0x00}
	PlaceHolder = []byte{0x01}
)

// StrategyStoreKey returns the byte representation of the strategy key
// Items are stored with key as follows:
// 0x03<strategyID>
func StrategyStoreKey(strategyID int64) []byte {
	strategyIdBz := sdktypes.Uint64ToBigEndian(uint64(strategyID))

	bz := make([]byte, len(StrategyKey)+len(strategyIdBz))

	copy(bz, StrategyKey)
	copy(bz[len(StrategyKey):], strategyIdBz)

	return bz
}

// VoucherStoreKey returns the byte representation of the voucher key
// Items are stored with key as follows:
// 0x04<voucherID>
func VoucherStoreKey(voucherID int64) []byte {
	voucherIdBz := sdktypes.Uint64ToBigEndian(uint64(voucherID))

	bz := make([]byte, len(VoucherKey)+len(voucherIdBz))

	copy(bz, VoucherKey)
	copy(bz[len(VoucherKey):], voucherIdBz)

	return bz
}

// VoucherByOwnerStoreKey returns the byte representation of the voucher by owner key
// Items are stored with key as follows:
// 0x05<owner><Delimiter(1 Byte)><voucherID>
func VoucherByOwnerStoreKey(owner sdktypes.AccAddress, voucherID uint64) []byte {
	address.MustLengthPrefix(owner)
	voucherIdBz := sdktypes.Uint64ToBigEndian(voucherID)

	bz := make([]byte, len(VoucherByOwnerKey)+len(owner)+len(Delimiter)+len(voucherIdBz))

	copy(bz, VoucherByOwnerKey)
	copy(bz[len(VoucherByOwnerKey):], owner)
	copy(bz[len(VoucherByOwnerKey)+len(owner):], Delimiter)
	copy(bz[len(VoucherByOwnerKey)+len(owner)+len(Delimiter):], voucherIdBz)

	return bz
}

// VoucherByOwnerStorePrefixKey returns the byte representation of the voucher by owner prefix key
// Items are stored with key as follows:
// 0x05<owner><Delimiter(1 Byte)>
func VoucherByOwnerStorePrefixKey(owner sdktypes.AccAddress) []byte {
	address.MustLengthPrefix(owner)

	bz := make([]byte, len(VoucherByOwnerKey)+len(owner)+len(Delimiter))

	copy(bz, VoucherByOwnerKey)
	copy(bz[len(VoucherByOwnerKey):], owner)
	copy(bz[len(VoucherByOwnerKey)+len(owner):], Delimiter)

	return bz
}
