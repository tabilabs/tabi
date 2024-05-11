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

// GetDigest returns the digest.
func (k Keeper) GetDigest(ctx sdk.Context, epochID uint64) (*types.ReportDigest, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DigestOnEpochStoreKey(epochID))
	if len(bz) == 0 {
		return nil, false
	}

	var digest types.ReportDigest
	k.cdc.Unmarshal(bz, &digest)
	return &digest, true
}

// setDigest sets the digest.
func (k Keeper) setDigest(ctx sdk.Context, epochID uint64, digest *types.ReportDigest) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := k.cdc.Marshal(digest)
	key := types.DigestOnEpochStoreKey(epochID)
	store.Set(key, bz)
}

// delDigest deletes the digest.
func (k Keeper) delDigest(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DigestOnEpochStoreKey(epochID))
}

// HasReportBatch checks if the batch exists.
func (k Keeper) HasReportBatch(ctx sdk.Context, epochID, batchID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ReportBatchOnEpochStoreKey(epochID, batchID))
}

// SetReportBatch sets the batch count.
func (k Keeper) SetReportBatch(ctx sdk.Context, epochID, batchID, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	key := types.ReportBatchOnEpochStoreKey(epochID, batchID)
	store.Set(key, bz)
}

// GetReportBatches returns the batch count.
func (k Keeper) GetReportBatches(ctx sdk.Context, epochID uint64) []types.BatchBase {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReportBatchOnEpochPrefixKey(epochID))
	defer iterator.Close()

	var batches []types.BatchBase
	for ; iterator.Valid(); iterator.Next() {
		batchID := sdk.BigEndianToUint64(iterator.Key())
		count := sdk.BigEndianToUint64(iterator.Value())
		batches = append(batches, types.BatchBase{
			Id:    batchID,
			Count: count,
		})
	}

	return batches
}

// delReportBatches deletes the batch count.
func (k Keeper) delReportBatches(ctx sdk.Context, epochID uint64) {
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

// delEndEpoch deletes the end epoch.
func (k Keeper) delEndEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.EndOnEpochStoreKey(epochID))
}

// GetEpochBase returns the base information of an epoch.
func (k Keeper) GetEpochBase(ctx sdk.Context, epochID uint64) types.EpochBase {
	emission := k.GetEpochEmission(ctx, epochID)
	computingPower := k.GetComputingPowerSumOnEpoch(ctx, epochID)
	pledgeAmount := k.GetPledgeSum(ctx, epochID)

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
	digest, _ := k.GetDigest(ctx, epochId)
	curr := k.GetEpochBase(ctx, epochId)
	prev := k.GetEpochBase(ctx, epochId-1)
	emissionSum := k.GetEmissionClaimedSum(ctx)

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

// SetEpochBase sets epoch base info.
func (k Keeper) SetEpochBase(ctx sdk.Context, epochID uint64, base types.EpochBase) {
	if epochID < 1 {
		return
	}

	if !base.ComputingPowerSum.IsZero() {
		k.setComputingPowerSumOnEpoch(ctx, epochID, base.ComputingPowerSum)
	}
	if !base.PledgeAmountSum.IsZero() {
		k.setPledgeSum(ctx, epochID, base.PledgeAmountSum)
	}
	if !base.EmissionSum.IsZero() {
		k.setEpochEmission(ctx, epochID, base.EmissionSum)
	}
}
