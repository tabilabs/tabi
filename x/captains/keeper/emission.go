package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// Emission handles all emission calculations.

// CalcEpochEmissionByNode returns the emission reward for a node in an epoch.
func (k Keeper) CalcEpochEmissionByNode(
	ctx sdk.Context,
	epochID uint64,
	nodeID string,
	powerOnRatio sdk.Dec,
	emissionSum sdk.Dec,
) (sdk.Dec, error) {
	cpr, err := k.CalcNodeComputingPowerRatioOnEpoch(ctx, epochID, nodeID, powerOnRatio)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return emissionSum.Mul(cpr), nil
}

// CalcEpochEmissionSum returns the total emission reward for an epoch.
func (k Keeper) CalcEpochEmissionSum(ctx sdk.Context, epochID uint64, onOperationRatio sdk.Dec) (sdk.Dec, error) {
	base := k.CalcBaseEpochEmission(ctx)
	pledgeRatio, err := k.CalcGlobalPledgeRatio(ctx, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return base.Mul(pledgeRatio).Mul(onOperationRatio), nil
}

// CalcBaseEpochEmission returns the base emission reward for an epoch.
func (k Keeper) CalcBaseEpochEmission(ctx sdk.Context) sdk.Dec {
	tech := k.CalculateTechProgressCoefficient(ctx)
	halving := k.GetHalvingEraCoefficient(ctx)
	cc := k.GetCaptainsConstant(ctx)
	return tech.Mul(halving).Mul(cc)
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

// GetOwnerHistoricalEmissionOnEpoch returns the historical emission for an owner at the end of an epoch.
func (k Keeper) GetOwnerHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, owner sdk.AccAddress) (sdk.Dec, error) {
	nodes := k.GetNodesByOwner(ctx, owner)
	total := sdk.ZeroDec()

	for _, node := range nodes {
		amount := k.GetNodeHistoricalEmissionOnEpoch(ctx, epochID, node.Id)
		total = total.Add(amount)
	}
	return total, nil
}

// GetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	// TODO: add epoch safe check here in case we call at epoch(t) for epoch(t-1) data.
	// TODO: return error as well.
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// HasWithdrawableRewardsOnNode check if the node has rewards to be withdrawn.
func (k Keeper) HasWithdrawableRewardsOnNode(ctx sdk.Context, nodeID string) bool {
	epoch := k.GetCurrentEpoch(ctx)
	lastClaim := k.GetNodeHistoricalEmissionOnLastClaim(ctx, nodeID)
	epochBefore := k.GetNodeHistoricalEmissionOnEpoch(ctx, epoch-2, nodeID)
	return lastClaim.LT(epochBefore)
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

// UpdateNodeHistoricalEmissionOnEpoch updates node_historical_emission_on_epoch after the user claim.
func (k Keeper) setNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeHistoricalEmissionOnLastClaimStoreKey(nodeID)
	store.Set(key, []byte(amount.String()))
}
