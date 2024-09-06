package keeper

import (
	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	pledgeRatioUpperBound                   = sdk.OneDec()
	globalPledgeRatioLowerBound             = sdk.NewDecWithPrec(3, 1)
	ownerPledgeRatioUpperBoundWhenCalcPower = sdk.NewDecWithPrec(3, 1)
)

// CalcGlobalPledgeRatio calculates the pledge rate of the global on the epoch t.
func (k Keeper) CalcGlobalPledgeRatio(ctx sdk.Context, epochID uint64) sdk.Dec {
	sum := k.GetGlobalClaimedEmission(ctx)
	if sum.IsZero() {
		return sdk.OneDec()
	}

	pledgeSum := k.GetGlobalPledge(ctx, epochID)
	ratio := pledgeSum.Quo(sum)
	// no more than 1.0
	if ratio.GT(pledgeRatioUpperBound) {
		ratio = pledgeRatioUpperBound
	}
	// no less than 0.3
	if ratio.LT(globalPledgeRatioLowerBound) {
		ratio = globalPledgeRatioLowerBound
	}
	return ratio
}

// CalcNodePledgeRatioOnEpoch calculates the pledge rate of the node on the epoch t but also prune pledge 2 epochs before.
func (k Keeper) CalcNodePledgeRatioOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	owner, _ := k.GetNodeOwner(ctx, nodeID)

	claimed := k.GetOwnerClaimedEmission(ctx, owner)
	if claimed.Equal(sdk.ZeroDec()) {
		return pledgeRatioUpperBound
	}

	ownerPledge := k.GetOwnerPledge(ctx, owner, epochID)
	nodePledgeRatio := ownerPledge.Quo(claimed)

	// no more than 1
	if nodePledgeRatio.GT(pledgeRatioUpperBound) {
		nodePledgeRatio = pledgeRatioUpperBound
	}

	return nodePledgeRatio
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

func (k Keeper) HasGlobalPledge(ctx sdk.Context, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalPledgeOnEpochStoreKey(epochID)
	return store.Has(key)
}

// GetGlobalPledge returns the total pledge amount of one epoch end.
func (k Keeper) GetGlobalPledge(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalPledgeOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// IncrGlobalPledge increments the total pledge amount of captains' owners on the epoch end.
func (k Keeper) IncrGlobalPledge(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	sum := k.GetGlobalPledge(ctx, epochID)
	sum = sum.Add(amount)
	k.SetGlobalPledge(ctx, epochID, sum)
}

// SetGlobalPledge sets the total pledge amount of captains' owners on the epoch end.
func (k Keeper) SetGlobalPledge(ctx sdk.Context, epochID uint64, sum sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalPledgeOnEpochStoreKey(epochID)
	store.Set(key, []byte(sum.String()))
}

// DelGlobalPledge deletes the total pledge amount of captains' owners on the epoch end.
func (k Keeper) DelGlobalPledge(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalPledgeOnEpochStoreKey(epochID)
	store.Delete(key)
}

// Genesis State Export/Import Helpers

// GetGlobalsPledge returns all global pledge.
func (k Keeper) GetGlobalsPledge(ctx sdk.Context) []types.GlobalPledge {
	var globalsPledge []types.GlobalPledge
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GlobalPledgeOnEpochKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var globalPledge types.GlobalPledge
		globalPledge.EpochId = types.SplitEpochFromStoreKey(types.GlobalPledgeOnEpochKey, iterator.Key())
		globalPledge.Amount = k.GetGlobalPledge(ctx, globalPledge.EpochId)
		globalsPledge = append(globalsPledge, globalPledge)
	}
	return globalsPledge
}

// GetOwnersPledge returns all owner pledge.
func (k Keeper) GetOwnersPledge(ctx sdk.Context) []types.OwnerPledge {
	var ownersPledge []types.OwnerPledge
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OwnerPledgeOnEpochKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ownerPledge types.OwnerPledge
		epochId, owner := types.SplitEpochAndStrFromStoreKey(types.OwnerPledgeOnEpochKey, iterator.Key())
		ownerPledge.EpochId = epochId
		ownerPledge.Owner = sdk.AccAddress(owner).String()
		ownerPledge.Amount = k.GetOwnerPledge(ctx, sdk.AccAddress(owner), epochId)
		ownersPledge = append(ownersPledge, ownerPledge)
	}
	return ownersPledge
}
