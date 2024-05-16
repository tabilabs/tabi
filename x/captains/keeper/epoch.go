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

// GetEpochBase returns the base information of an epoch.
func (k Keeper) GetEpochBase(ctx sdk.Context, epochID uint64) types.EpochBase {
	emission := k.GetEpochEmission(ctx, epochID)
	computingPower := k.GetGlobalComputingPowerOnEpoch(ctx, epochID)
	pledgeAmount := k.GetGlobalPledge(ctx, epochID)

	return types.EpochBase{
		EmissionSum:       emission,
		ComputingPowerSum: computingPower,
		PledgeAmountSum:   pledgeAmount,
	}
}

// GetEpochsState returns the state of the epochs.
func (k Keeper) GetEpochsState(ctx sdk.Context) types.EpochState {
	epochId := k.GetCurrentEpoch(ctx)
	isEnd := k.HasEndEpoch(ctx, epochId)
	batches := k.GetReportBatches(ctx, epochId)
	digest, _ := k.GetReportDigest(ctx, epochId)
	curr := k.GetEpochBase(ctx, epochId)
	prev := k.GetEpochBase(ctx, epochId-1)
	emissionSum := k.GetGlobalClaimedEmission(ctx)

	return types.EpochState{
		CurrEpoch:          epochId,
		IsEnd:              isEnd,
		Digest:             digest,
		Batches:            batches,
		Current:            curr,
		Previous:           prev,
		EmissionClaimedSum: emissionSum,
	}
}

// setEpochBase sets epoch base info.
func (k Keeper) setEpochBase(ctx sdk.Context, epochID uint64, base types.EpochBase) {
	if epochID < 1 {
		return
	}

	if !base.ComputingPowerSum.IsZero() {
		k.setGlobalComputingPowerOnEpoch(ctx, epochID, base.ComputingPowerSum)
	}
	if !base.PledgeAmountSum.IsZero() {
		k.SetGlobalPledge(ctx, epochID, base.PledgeAmountSum)
	}
	if !base.EmissionSum.IsZero() {
		k.setEpochEmission(ctx, epochID, base.EmissionSum)
	}
}
