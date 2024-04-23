package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// Emission handles all emission calculations.

// CalcEpochEmissionByNode returns the emission reward for a node in an epoch.
func (k Keeper) CalcEpochEmissionByNode(ctx sdk.Context, epochID uint64, nodeID string) {
	panic("not implemented")
}

// CalcEpochEmissionSum returns the total emission reward for an epoch.
func (k Keeper) CalcEpochEmissionSum(ctx sdk.Context, onOperationRatio sdk.Dec) sdk.Dec {
	panic("not implemented")
}

// CalcEpochEmissionBase returns the base emission reward for an epoch.
func (k Keeper) CalcEpochEmissionBase(ctx sdk.Context) sdk.Dec {
	panic("not implemented")
}

// GetHistoricalEmissionSum returns the historical emission sum at the end of a epoch.
func (k Keeper) GetHistoricalEmissionSum(ctx sdk.Context, epochID uint64) sdk.Dec {
	panic("not implemented")
}

// GetNodeHistoricalEmissionOnEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	panic("not implemented")
}

// HasWithdrawableRewardsOnNode check if the node has rewards to be withdrawn.
func (k Keeper) HasWithdrawableRewardsOnNode(ctx sdk.Context, nodeID string) bool {
	epoch := k.GetCurrentEpoch(ctx)
	lastClaim := k.GetNodeHistoricalEmissionOnLastClaim(ctx, nodeID)
	epochBefore := k.GetNodeHistoricalEmissionOnEpoch(ctx, epoch-1, nodeID)
	return lastClaim.Equal(epochBefore)
}

// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
// That says, if the user claims on epoch(T+1), the user will get rewards accrued at the end of epoch(T).
// The next time the user claims, let's say on epoch(T+k+1), we can easily calc the rewards by subtracting
// node_historical_emission_on_last_claim from node_historical_emission(T+k)
func (k Keeper) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	panic("not implemented")
}

// UpdateNodeHistoricalEmissionOnLastClaim updates node_historical_emission_on_last_claim after the user claim.
func (k Keeper) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	panic("not implemented")
}
