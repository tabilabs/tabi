package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// calcGlobalPledgeRatio calculates the pledge rate of the global on the epoch t.
func (k Keeper) calcGlobalPledgeRatio(ctx sdk.Context, epochID uint64) (sdk.Dec, error) {
	if epochID == 1 {
		return sdk.NewDecWithPrec(3, 1), nil
	}

	sum, err := k.GetHistoricalEmissionSum(ctx, epochID-1)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	pledgeSum, err := k.GetPledgeSum(ctx, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	k.delPledgeSum(ctx, epochID)

	return pledgeSum.Quo(sum), nil
}

// calcNodePledgeRatioOnEpoch calculates the pledge rate of the node on the epoch t.
// node_pledge_ratio(t) = owner_pledge(t) / owner_historical_emission_sum(t-1)
func (k Keeper) calcNodePledgeRatioOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) (sdk.Dec, error) {
	if epochID == 0 {
		return sdk.OneDec(), nil
	}

	owner := k.GetNodeOwner(ctx, nodeID)
	ownerPledge, err := k.GetOwnerPledge(ctx, owner, epochID)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	if epochID > 1 {
		k.delOwnerPledge(ctx, owner, epochID-1)
	}

	ownerHistoricalEmissionSum := k.calcOwnerHistoricalEmissionSum(ctx, epochID-1, owner)

	nodePledgeRatio := ownerPledge.Quo(ownerHistoricalEmissionSum)
	if nodePledgeRatio.GTE(sdk.OneDec()) {
		nodePledgeRatio = sdk.OneDec()
	}

	return nodePledgeRatio, nil
}

// SampleOwnerPledge sample pledge amount of the owner on the epoch.
func (k Keeper) SampleOwnerPledge(ctx sdk.Context, owner sdk.AccAddress) (sdk.Dec, error) {
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

	return totalAmount, nil
}

// GetOwnerPledge returns the sampled pledge amount of the owner on the epoch.
func (k Keeper) GetOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	bz := store.Get(key)
	return sdk.NewDecFromStr(string(bz))
}

// setOwnerPledge sets the sampled pledge amount of the owner on the epoch.
func (k Keeper) setOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64, pledge sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	store.Set(key, []byte(pledge.String()))
}

// delOwnerPledge deletes the sampled pledge amount of the owner on the epoch.
func (k Keeper) delOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	store.Delete(key)
}

// GetPledgeSum returns the total pledge amount of captains' owners on the epoch end.
func (k Keeper) GetPledgeSum(ctx sdk.Context, epochID uint64) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	res, err := sdk.NewDecFromStr(string(bz))
	if err != nil {
		return res, err
	}
	return res, nil
}

// setPledgeSum sets the total pledge amount of captains' owners on the epoch end.
func (k Keeper) setPledgeSum(ctx sdk.Context, epochID uint64, sum sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(sum.String()))
}

// delPledgeSum deletes the total pledge amount of captains' owners on the epoch end.
func (k Keeper) delPledgeSum(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	store.Delete(key)
}
