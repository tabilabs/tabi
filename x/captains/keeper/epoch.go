package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// epochs handle stage transitions in an epoch.

// GetCurrentEpoch returns the current epoch.
func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochKey)
	return sdk.BigEndianToUint64(bz)
}

// setEpoch sets the epoch id.
func (k Keeper) setEpoch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(k.GetCurrentEpoch(ctx) + 1)
	store.Set(types.EpochKey, bz)
}

func (k Keeper) GetDigest(ctx sdk.Context, epochID uint64) types.ReportDigest {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DigestOnEpochStoreKey(epochID))
	var digest types.ReportDigest
	k.cdc.Unmarshal(bz, &digest)
	return digest
}

func (k Keeper) setDigest(ctx sdk.Context, epochID uint64, digest *types.ReportDigest) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := k.cdc.Marshal(digest)
	key := types.DigestOnEpochStoreKey(epochID)
	store.Set(key, bz)
}

func (k Keeper) setBatchCount(ctx sdk.Context, epochID, batchID, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	key := types.BatchCountOnEpochStoreKey(epochID, batchID)
	store.Set(key, bz)
}

func (k Keeper) setEndEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.EndOnEpochStoreKey(epochID)
	store.Set(key, []byte{0x01})
}
