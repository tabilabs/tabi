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

	pledgeRatio, err := k.calcGlobalPledgeRatio(ctx, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}

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

// GetHistoricalEmissionSum returns the historical emission sum at the end of a epoch.
func (k Keeper) GetHistoricalEmissionSum(ctx sdk.Context, epochID uint64) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalEmissionSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	res, err := sdk.NewDecFromStr(string(bz))
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return res, nil
}

// setHistoricalEmissionSum sets the historical emission sum for an epoch.
func (k Keeper) setHistoricalEmissionSum(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalEmissionSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(amount.String()))
}

// delHistoricalEmissionSum deletes the historical emission sum for an epoch.
func (k Keeper) delHistoricalEmissionSum(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalEmissionSumOnEpochStoreKey(epochID)
	store.Delete(key)
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
	historyEmission := k.GetNodeHistoricalEmissionOnEpoch(ctx, epochID, nodeID)
	// emission already exists
	if !historyEmission.Equal(sdk.ZeroDec()) {
		return historyEmission
	}

	prevHistoryEmission := sdk.ZeroDec()
	if epochID >= 2 {
		// avoid overflow when epochID is 0, but we shouldn't be worried about a max uint64.
		prevHistoryEmission = k.GetNodeHistoricalEmissionOnEpoch(ctx, epochID-1, nodeID)
	}

	emission, _ := k.GetEpochEmission(ctx, epochID)
	emission.Add(prevHistoryEmission)

	// set and del
	k.setNodeHistoricalEmissionOnEpoch(ctx, epochID, nodeID, emission)
	k.delNodeHistoricalEmissionOnEpoch(ctx, epochID-1, nodeID)

	return emission
}

// GetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	// TODO: add epoch safe check here in case we call at epoch(t) for epoch(t-1) data.
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
// That says, if the user claims on epoch(T+1), the user will get rewards accrued at the end of epoch(T).
// The next time the user claims, let's say on epoch(T+k+1), we can easily calc the rewards by subtracting
// node_historical_emission_on_last_claim from node_historical_emission(T+k)
func (k Keeper) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	bz := store.Get(key)
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// UpdateNodeHistoricalEmissionOnLastClaim updates node_historical_emission_on_last_claim after the user claim.
// NOTE: call this function after the user claims the rewards.
func (k Keeper) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	epoch := k.GetCurrentEpoch(ctx)
	// NOTE: if we are in epoch(t), we are calculating the rewards for epoch(t-1) so we
	// can only get historical emssion by the end of epoch(t-2).
	amount := k.GetNodeHistoricalEmissionOnEpoch(ctx, epoch-2, nodeID)

	k.setNodeHistoricalEmissionOnLastClaim(ctx, nodeID, amount)
	return nil
}

// setNodeHistoricalEmissionOnLastClaim sets the historical emission the last time user claimed.ß
func (k Keeper) setNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	store.Set(key, []byte(amount.String()))
}
