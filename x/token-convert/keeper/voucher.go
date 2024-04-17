package keeper

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

// genVoucherID generates a unique voucher id with hash algo.
func (k Keeper) genVoucherID(ctx sdk.Context) string {
	seq := k.GetVoucherSeq(ctx)
	voucherID := fmt.Sprintf("voucher-%d", seq)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(voucherID)))
	k.setVoucherSeq(ctx, seq+1)
	return hash
}

// createVoucher creates and sets a voucher.
func (k Keeper) createVoucher(ctx sdk.Context, owner sdk.AccAddress, strategy string, amount sdk.Coin) string {
	store := ctx.KVStore(k.storeKey)

	voucher := types.Voucher{
		Id:          k.genVoucherID(ctx),
		Owner:       owner.String(),
		Amount:      amount,
		CreatedTime: ctx.BlockTime().Unix(),
		Strategy:    strategy,
	}
	bz := k.cdc.MustMarshal(&voucher)

	store.Set(types.VoucherStoreKey(voucher.Id), bz)

	return voucher.Id
}

func (k Keeper) deleteVoucher(ctx sdk.Context, voucherID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VoucherStoreKey(voucherID))
}

// GetVoucher gets a voucher by voucher id
func (k Keeper) GetVoucher(ctx sdk.Context, voucherID string) (types.Voucher, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VoucherStoreKey(voucherID))
	if len(bz) == 0 {
		return types.Voucher{}, false
	}

	var voucher types.Voucher
	k.cdc.MustUnmarshal(bz, &voucher)
	return voucher, true
}

// setVoucherByOwner sets the VoucherByOwnerStore.
func (k Keeper) setVoucherByOwner(ctx sdk.Context, owner sdk.AccAddress, voucherID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.VoucherByOwnerStoreKey(owner.Bytes(), voucherID), types.PlaceHolder)
}

// deleteVoucherByOwner deletes the VoucherByOwnerStore.
func (k Keeper) deleteVoucherByOwner(ctx sdk.Context, owner sdk.AccAddress, voucherID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.VoucherByOwnerStoreKey(owner.Bytes(), voucherID))
}

// GetVoucherSeq returns the next usable voucher seq
func (k Keeper) GetVoucherSeq(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VoucherSeqKey)

	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// setVoucherSeq sets the next usable voucher seq.
func (k Keeper) setVoucherSeq(ctx sdk.Context, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(seq)
	store.Set(types.VoucherSeqKey, bz)
}
