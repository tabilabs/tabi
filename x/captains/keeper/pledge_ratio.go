package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

var (
	globalPledgeRatioUpperBound = sdk.OneDec()
	globalPledgeRatioLowerBound = sdk.NewDecWithPrec(3, 1)
	ownerPledgeRatioUpperBound  = sdk.OneDec()
	ownerPledgeRatioLowerBound  = sdk.NewDecWithPrec(3, 1)
)

// CalcGlobalPledgeRatio calculates the pledge rate of the global on the epoch t.
func (k Keeper) CalcGlobalPledgeRatio(ctx sdk.Context, epochID uint64) sdk.Dec {
	sum := k.GetEmissionClaimedSum(ctx)
	if sum.IsZero() {
		return sdk.OneDec()
	}

	pledgeSum := k.GetPledgeSum(ctx, epochID)
	ratio := pledgeSum.Quo(sum)
	// no more than 1.0
	if ratio.GT(globalPledgeRatioUpperBound) {
		ratio = globalPledgeRatioUpperBound
	}
	// no less than 0.3
	if ratio.LT(globalPledgeRatioLowerBound) {
		ratio = globalPledgeRatioLowerBound
	}

	k.delPledgeSum(ctx, epochID)

	return ratio
}

// CalcNodePledgeRatioOnEpoch calculates the pledge rate of the node on the epoch t.
func (k Keeper) CalcNodePledgeRatioOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) (sdk.Dec, error) {
	owner := k.GetNodeOwner(ctx, nodeID)

	claimed := k.GetOwnerHistoricalEmissionOnLastClaim(ctx, owner)
	if claimed.Equal(sdk.ZeroDec()) {
		return ownerPledgeRatioUpperBound, nil
	}

	ownerPledge, _ := k.GetOwnerPledge(ctx, owner, epochID)
	nodePledgeRatio := ownerPledge.Quo(claimed)

	if nodePledgeRatio.GTE(ownerPledgeRatioLowerBound) {
		nodePledgeRatio = ownerPledgeRatioLowerBound
	}

	// prune previous pledge
	k.delOwnerPledge(ctx, owner, epochID-1)

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
func (k Keeper) GetOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec(), false
	}

	res, _ := sdk.NewDecFromStr(string(bz))
	return res, true
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
func (k Keeper) GetPledgeSum(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	res, _ := sdk.NewDecFromStr(string(bz))
	return res
}

// incrPledgeSum increments the total pledge amount of captains' owners on the epoch end.
func (k Keeper) incrPledgeSum(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	sum := k.GetPledgeSum(ctx, epochID)
	sum = sum.Add(amount)
	k.setPledgeSum(ctx, epochID, sum)
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
