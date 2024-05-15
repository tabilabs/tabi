package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

var (
	globalPledgeRatioUpperBound = sdk.OneDec()
	globalPledgeRatioLowerBound = sdk.NewDecWithPrec(3, 1)
	ownerPledgeRatioUpperBound  = sdk.NewDecWithPrec(3, 1)
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
	return ratio
}

// CalcNodePledgeRatioOnEpoch calculates the pledge rate of the node on the epoch t but also prune pledge 2 epochs before.
func (k Keeper) CalcNodePledgeRatioOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	owner := k.GetNodeOwner(ctx, nodeID)

	claimed := k.GetOwnerHistoricalEmissionOnLastClaim(ctx, owner)
	if claimed.Equal(sdk.ZeroDec()) {
		return ownerPledgeRatioUpperBound
	}

	ownerPledge := k.GetOwnerPledge(ctx, owner, epochID)
	nodePledgeRatio := ownerPledge.Quo(claimed)

	// no more than 0.3
	if nodePledgeRatio.GT(ownerPledgeRatioUpperBound) {
		nodePledgeRatio = ownerPledgeRatioUpperBound
	}

	return nodePledgeRatio
}

// Owner pledge amount

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
func (k Keeper) GetOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// HasOwnerPledge checks if the owner has the sampled pledge amount on the epoch.
func (k Keeper) HasOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerPledgeOnEpochStoreKey(owner, epochID)
	return store.Has(key)
}

// SetOwnerPledge sets the sampled pledge amount of the owner on the epoch.
func (k Keeper) SetOwnerPledge(ctx sdk.Context, owner sdk.AccAddress, epochID uint64, pledge sdk.Dec) {
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

// Pledge Sum

// GetPledgeSum returns the total pledge amount of one epoch end.
func (k Keeper) GetPledgeSum(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// IncrPledgeSum increments the total pledge amount of captains' owners on the epoch end.
func (k Keeper) IncrPledgeSum(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	sum := k.GetPledgeSum(ctx, epochID)
	sum = sum.Add(amount)
	k.SetPledgeSum(ctx, epochID, sum)
}

// SetPledgeSum sets the total pledge amount of captains' owners on the epoch end.
func (k Keeper) SetPledgeSum(ctx sdk.Context, epochID uint64, sum sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(sum.String()))
}

// DelPledgeSum deletes the total pledge amount of captains' owners on the epoch end.
func (k Keeper) DelPledgeSum(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.PledgeSumOnEpochStoreKey(epochID)
	store.Delete(key)
}
