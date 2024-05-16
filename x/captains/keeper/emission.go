package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// Emission handles all emission calculations.

// CalcBaseEpochEmission returns the base emission reward for an epoch.
func (k Keeper) CalcBaseEpochEmission(ctx sdk.Context) sdk.Dec {
	tech := k.CalcTechProgressCoefficient(ctx)
	halving := k.GetHalvingEraCoefficient(ctx)
	cc := k.GetCaptainsConstant(ctx)
	return tech.Mul(halving).Mul(cc)
}

// GetEpochEmission returns the emission reward for an epoch.
func (k Keeper) GetEpochEmission(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// CalcEpochEmission returns the total emission reward for an epoch.
func (k Keeper) CalcEpochEmission(ctx sdk.Context, epochID uint64, globalOperationRatio sdk.Dec) sdk.Dec {
	base := k.CalcBaseEpochEmission(ctx)
	pledgeRatio := k.CalcGlobalPledgeRatio(ctx, epochID)
	return base.Mul(pledgeRatio).Mul(globalOperationRatio)
}

// setEpochEmission sets the emission sum for an epoch.
func (k Keeper) setEpochEmission(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(amount.String()))
}

// delEpochEmission deletes the emission sum for an epoch.
func (k Keeper) delEpochEmission(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionSumOnEpochStoreKey(epochID)
	store.Delete(key)
}

// GetEmissionClaimedSum returns the emission claimed sum.
func (k Keeper) GetEmissionClaimedSum(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionClaimedSumKey
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// SetEmissionClaimedSum sets the emission claimed sum.
func (k Keeper) SetEmissionClaimedSum(ctx sdk.Context, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionClaimedSumKey
	store.Set(key, []byte(amount.String()))
}

// incrEmissionClaimedSum increases the emission claimed sum by amount.
func (k Keeper) incrEmissionClaimedSum(ctx sdk.Context, amount sdk.Dec) {
	emission := k.GetEmissionClaimedSum(ctx)
	emission.Add(amount)
	k.SetEmissionClaimedSum(ctx, emission)
}

// CalcAndSetNodeHistoricalEmissionByEpoch returns the historical emission for a node at the end of an epoch.
// NOTE: this function set the historical emission by the end of epoch(t) and removes that of epoch(t-1).
func (k Keeper) CalcAndSetNodeHistoricalEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	res := k.CalcNodeHistoricalEmissionByEpoch(ctx, epochID, nodeID)

	k.setNodeHistoricalEmissionByEpoch(ctx, epochID, nodeID, res)
	k.delNodeHistoricalEmissionByEpoch(ctx, epochID-1, nodeID)
	k.delNodeComputingPowerOnEpoch(ctx, epochID-1, nodeID)

	return res
}

// CalcNodeHistoricalEmissionByEpoch returns the historical emission for a node at the end of an epoch.
// NOTE: this func is only used to query and it won't set or prune state data.
func (k Keeper) CalcNodeHistoricalEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	if epochID == 0 {
		return sdk.ZeroDec()
	}

	// NOTE: we may have already set this value when user claimed the rewards before report handles it.
	historyEmission := k.GetNodeHistoricalEmissionByEpoch(ctx, epochID, nodeID)
	if !historyEmission.Equal(sdk.ZeroDec()) {
		return historyEmission
	}

	prevHistoryEmission := k.GetNodeHistoricalEmissionByEpoch(ctx, epochID-1, nodeID)
	epochEmission := k.CalcNodeEmissionOnEpoch(ctx, epochID, nodeID)

	return epochEmission.Add(prevHistoryEmission)
}

// CalcNodeEmissionOnEpoch returns the emission for a node at the end of an epoch.
func (k Keeper) CalcNodeEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	emission := k.GetEpochEmission(ctx, epochID)
	power := k.GetNodeComputingPowerOnEpoch(ctx, epochID, nodeID)
	powerSum := k.GetComputingPowerSumOnEpoch(ctx, epochID)
	return emission.Mul(power).Quo(powerSum)
}

// CalcAndGetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) CalcAndGetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	return k.CalcNodeHistoricalEmissionByEpoch(ctx, epochID, nodeID)
}

// GetNodeHistoricalEmissionByEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeHistoricalEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// setNodeHistoricalEmissionByEpoch sets the historical emission for a node at the end of an epoch.
func (k Keeper) setNodeHistoricalEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	store.Set(key, []byte(amount.String()))
}

// delNodeHistoricalEmissionByEpoch deletes the historical emission for a node at the end of an epoch.
func (k Keeper) delNodeHistoricalEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	store.Delete(key)
}

// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (k Keeper) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// UpdateNodeHistoricalEmissionOnLastClaim updates node_historical_emission_on_last_claim after the user claim.
// NOTE: call this function only after claiming the rewards.
func (k Keeper) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	epoch := k.GetCurrentEpoch(ctx)
	amount := k.CalcNodeHistoricalEmissionByEpoch(ctx, epoch-1, nodeID)

	k.SetNodeHistoricalEmissionOnLastClaim(ctx, nodeID, amount)
	k.incrEmissionClaimedSum(ctx, amount)

	// TODO: we don't need to return error? we should return sum right?
	return nil
}

// SetNodeHistoricalEmissionOnLastClaim sets the historical emission the last time user claimed.ß
func (k Keeper) SetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	store.Set(key, []byte(amount.String()))
}

// GetOwnerHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (k Keeper) GetOwnerHistoricalEmissionOnLastClaim(ctx sdk.Context, owner sdk.AccAddress) sdk.Dec {
	nodes := k.GetNodesByOwner(ctx, owner)
	total := sdk.ZeroDec()
	for _, node := range nodes {
		total = total.Add(k.GetNodeHistoricalEmissionOnLastClaim(ctx, node.Id))
	}
	return total
}
