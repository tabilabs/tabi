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

// incrEpoch enters a new epoch.
func (k Keeper) incrEpoch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(k.GetCurrentEpoch(ctx) + 1)
	store.Set(types.CurrEpochKey, bz)
}

// setEpoch sets the epoch id.
func (k Keeper) setEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(epochID)
	store.Set(types.CurrEpochKey, bz)
}

// setEndOnEpoch sets the end epoch.
func (k Keeper) setEndOnEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.EndOnEpochStoreKey(epochID)
	store.Set(key, types.PlaceHolder)
}

// HasEndEpoch checks if the end epoch exists.
func (k Keeper) HasEndEpoch(ctx sdk.Context, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.EndOnEpochStoreKey(epochID))
}

// delEndEpoch deletes the end epoch.
func (k Keeper) delEndEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.EndOnEpochStoreKey(epochID))
}

// IsStandByPhase checks if the stand by phrase is active.
func (k Keeper) IsStandByPhase(ctx sdk.Context) bool {
	return !k.HasStandByOverFlag(ctx)
}

// HasStandByOverFlag checks if the stand by flag exists.
func (k Keeper) HasStandByOverFlag(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.StandByOverKey)
}

// setStandBy sets the stand by flag.
// NOTE: if set, the stand-by phrase is over in current epoch.
func (k Keeper) setStandByOverFlag(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.StandByOverKey, types.PlaceHolder)
}

// delStandByOverFlag deletes the stand by flag.
func (k Keeper) delStandByOverFlag(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.StandByOverKey)
}
