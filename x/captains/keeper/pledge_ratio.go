package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// CalcGlobalPledgeRatio calculates the pledge rate of the global on the epoch t.
// global_pledge_ratio(t) = pledge_sum(t) / historical_emission_sum(t-1)
func (k Keeper) CalcGlobalPledgeRatio(ctx sdk.Context, epochID uint64) (sdk.Dec, error) {
	sum, err := k.GetHistoricalEmissionSum(ctx, epochID-1)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	pledgeSum, err := k.GetPledgeSum(ctx, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return pledgeSum.Quo(sum), nil
}

// CalcNodePledgeRatioOnEpoch calculates the pledge rate of the node on the epoch t.
// node_pledge_ratio(t) = owner_pledge(t) / owner_historical_emission_sum(t-1)
func (k Keeper) CalcNodePledgeRatioOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) (sdk.Dec, error) {
	owner := k.GetNodeOwner(ctx, nodeID)
	ownerPledge, err := k.CalcOwnerPledge(ctx, owner, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	ownerHistoricalEmissionSum, err := k.GetOwnerHistoricalEmissionOnEpoch(ctx, epochID-1, owner)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	nodePledgeRatio := ownerPledge.Quo(ownerHistoricalEmissionSum)
	if nodePledgeRatio.GTE(sdk.OneDec()) {
		nodePledgeRatio = sdk.OneDec()
	}

	return nodePledgeRatio, nil
}

// CalcOwnerPledge returns the pledge amount of the owner on the epoch end.
// FIXME: we want to sample pledge on one block during the target epoch.
func (k Keeper) CalcOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) (sdk.Dec, error) {
	stakingParams := k.stakingKeeper.GetParams(ctx)
	maxRetrieve := stakingParams.GetMaxValidators()

	// FIXME: it seems that we don't have a good way to get owner delegations at the end of one epoch.
	delegations := k.stakingKeeper.GetDelegatorDelegations(ctx, owner, uint16(maxRetrieve))
	totalAmount := sdk.ZeroDec()

	for _, delegation := range delegations {
		val, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
		if !found {
			continue
		}
		totalAmount = totalAmount.Add(val.TokensFromShares(delegation.GetShares()))
	}
	// NOTE: there's no need to do truncation here
	return totalAmount, nil
}

// GetPledgeSum returns the total pledge amount of captains' owners on the epoch end.
func (k Keeper) GetPledgeSum(ctx sdk.Context, epochID uint64) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeAmountSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	res, err := sdk.NewDecFromStr(string(bz))
	if err != nil {
		return res, err
	}
	return res, nil
}

// SetPledgeSum sets the total pledge amount of captains' owners on the epoch end.
func (k Keeper) SetPledgeSum(ctx sdk.Context, epochID uint64, sum sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeAmountSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(sum.String()))
}
