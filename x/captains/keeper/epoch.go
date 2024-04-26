package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// epochs handle stage transitions in an epoch.

// GetCurrentEpoch returns the current epoch.
func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrEpochKey)
	return sdk.BigEndianToUint64(bz)
}

// EnterNewEpoch enters a new epoch.
func (k Keeper) EnterNewEpoch(ctx sdk.Context) {
	k.setEpoch(ctx)
}

// setEpoch sets the epoch id.
func (k Keeper) setEpoch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(k.GetCurrentEpoch(ctx) + 1)
	store.Set(types.CurrEpochKey, bz)
}

// GetDigest returns the digest.
func (k Keeper) GetDigest(ctx sdk.Context, epochID uint64) types.ReportDigest {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DigestOnEpochStoreKey(epochID))
	var digest types.ReportDigest
	k.cdc.Unmarshal(bz, &digest)
	return digest
}

// setDigest sets the digest.
func (k Keeper) setDigest(ctx sdk.Context, epochID uint64, digest *types.ReportDigest) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := k.cdc.Marshal(digest)
	key := types.DigestOnEpochStoreKey(epochID)
	store.Set(key, bz)
}

// setReportBatch sets the batch count.
func (k Keeper) setReportBatch(ctx sdk.Context, epochID, batchID, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	key := types.ReportBatchOnEpochStoreKey(epochID, batchID)
	store.Set(key, bz)
}

// DelReportBatches deletes the batch count.
func (k Keeper) DelReportBatches(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReportBatchOnEpochPrefixKey(epochID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// setEndEpoch sets the end epoch.
func (k Keeper) setEndEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.EndOnEpochStoreKey(epochID)
	store.Set(key, types.PlaceHolder)
}

// HasEndEpoch checks if the end epoch exists.
func (k Keeper) HasEndEpoch(ctx sdk.Context, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.EndOnEpochStoreKey(epochID))
}

// DelEndEpoch deletes the end epoch.
func (k Keeper) DelEndEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.EndOnEpochStoreKey(epochID))
}
