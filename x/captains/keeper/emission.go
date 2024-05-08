package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// Emission handles all emission calculations.

// GetBaseEpochEmission returns the base emission reward for an epoch.
func (k Keeper) GetBaseEpochEmission(ctx sdk.Context) sdk.Dec {
	tech := k.GetTechProgressCoefficient(ctx)
	halving := k.GetHalvingEraCoefficient(ctx)
	cc := k.GetCaptainsConstant(ctx)
	return tech.Mul(halving).Mul(cc)
}

// GetEpochEmission returns the emission reward for an epoch.
func (k Keeper) GetEpochEmission(ctx sdk.Context, epochID uint64) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec(), false
	}
	res, _ := sdk.NewDecFromStr(string(bz))
	return res, true
}

// calcEpochEmission returns the total emission reward for an epoch.
func (k Keeper) calcEpochEmission(ctx sdk.Context, epochID uint64, globalOperationRatio sdk.Dec) (sdk.Dec, error) {
	base := k.GetBaseEpochEmission(ctx)
	pledgeRatio := k.CalcGlobalPledgeRatio(ctx, epochID)
	sum := base.Mul(pledgeRatio).Mul(globalOperationRatio)

	k.setEpochEmission(ctx, epochID, sum)

	return sum, nil
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
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// setEmissionClaimedSum sets the emission claimed sum.
func (k Keeper) setEmissionClaimedSum(ctx sdk.Context, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.EmissionClaimedSumKey
	store.Set(key, []byte(amount.String()))
}

// incrEmissionClaimedSum increases the emission claimed sum by amount.
func (k Keeper) incrEmissionClaimedSum(ctx sdk.Context, amount sdk.Dec) {
	emission := k.GetEmissionClaimedSum(ctx)
	emission.Add(amount)
	k.setEmissionClaimedSum(ctx, emission)
}

// calcOwnerHistoricalEmissionSum returns the historical emission sum for an owner at the end of an epoch.
func (k Keeper) calcOwnerHistoricalEmissionSum(ctx sdk.Context, epochID uint64, owner sdk.AccAddress) sdk.Dec {
	nodes := k.GetNodesByOwner(ctx, owner)
	total := sdk.ZeroDec()

	for _, node := range nodes {
		amount := k.calNodeHistoricalEmissionOnEpoch(ctx, epochID, node.Id)
		total = total.Add(amount)
	}
	return total
}

// calNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
// NOTE: this function set the historical emission by the end of epoch(t) and removes that of epoch(t-1).
func (k Keeper) calNodeHistoricalEmissionOnEpoch(
	ctx sdk.Context,
	epochID uint64,
	nodeID string,
) sdk.Dec {
	if epochID == 0 {
		return sdk.ZeroDec()
	}

	historyEmission := k.GetNodeHistoricalEmissionOnEpoch(ctx, epochID, nodeID)
	if !historyEmission.Equal(sdk.ZeroDec()) {
		return historyEmission
	}

	prevHistoryEmission := k.GetNodeHistoricalEmissionOnEpoch(ctx, epochID-1, nodeID)
	emission, _ := k.GetEpochEmission(ctx, epochID)
	power := k.GetNodeComputingPowerOnEpoch(ctx, epochID, nodeID)
	powerSum := k.GetComputingPowerSumOnEpoch(ctx, epochID)

	res := emission.Mul(power).Quo(powerSum).Add(prevHistoryEmission)

	k.setNodeHistoricalEmissionOnEpoch(ctx, epochID, nodeID, emission)
	k.delNodeHistoricalEmissionOnEpoch(ctx, epochID-1, nodeID)
	k.delNodeComputingPowerOnEpoch(ctx, epochID-1, nodeID)

	return res
}

// CalAndGetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) CalAndGetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	return k.calNodeHistoricalEmissionOnEpoch(ctx, epochID, nodeID)
}

// GetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// setNodeHistoricalEmissionOnEpoch sets the historical emission for a node at the end of an epoch.
func (k Keeper) setNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	store.Set(key, []byte(amount.String()))
}

// delNodeHistoricalEmissionOnEpoch deletes the historical emission for a node at the end of an epoch.
func (k Keeper) delNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	store.Delete(key)
}

// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (k Keeper) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	bz := store.Get(key)
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// UpdateNodeHistoricalEmissionOnLastClaim updates node_historical_emission_on_last_claim after the user claim.
// NOTE: call this function only after claiming the rewards.
func (k Keeper) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	epoch := k.GetCurrentEpoch(ctx)

	amount := k.GetNodeHistoricalEmissionOnEpoch(ctx, epoch-1, nodeID)
	k.setNodeHistoricalEmissionOnLastClaim(ctx, nodeID, amount)

	k.incrEmissionClaimedSum(ctx, amount)

	return nil
}

// setNodeHistoricalEmissionOnLastClaim sets the historical emission the last time user claimed.ß
func (k Keeper) setNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	store.Set(key, []byte(amount.String()))
}

// GetOwnerHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (k Keeper) GetOwnerHistoricalEmissionOnLastClaim(ctx sdk.Context, owner sdk.AccAddress) sdk.Dec {
	nodes := k.GetNodes(ctx)
	total := sdk.ZeroDec()
	for _, node := range nodes {
		total = total.Add(k.GetNodeHistoricalEmissionOnLastClaim(ctx, node.Id))
	}
	return total
}
