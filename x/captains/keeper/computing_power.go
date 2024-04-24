package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

func (k Keeper) ComputingPowerRatioByNode(ctx sdk.Context, epochID, nodeID string) {
	panic("implement me")
}

func (k Keeper) ComputingPowerSum(ctx sdk.Context, epochID string) {
	panic("implement me")
}

func (k Keeper) ComputingPowerByNode(ctx sdk.Context, epochID, nodeID string) {
	panic("implement me")
}

// ComputingPowerBaseByNode returns the base computing power of a node as per its node info.
func (k Keeper) ComputingPowerBaseByNode(ctx sdk.Context, nodeID string) {
	panic("implement me")
}

// CommitComputingPower commits the pending computing power.
func (k Keeper) CommitComputingPower(ctx sdk.Context, amount uint64, owner sdk.AccAddress) (uint64, uint64, error) {
	before := k.GetComputingPowerClaimable(ctx, owner)
	after := before + amount
	k.setComputingPowerClaimable(ctx, after, owner)
	return before, after, nil
}

// incrComputingPowerClaimable decrements the claimable computing power of an owner.
func (k Keeper) incrComputingPowerClaimable(ctx sdk.Context, amount uint64, owner sdk.AccAddress) {
	before := k.GetComputingPowerClaimable(ctx, owner)
	after := before + amount
	k.setComputingPowerClaimable(ctx, after, owner)
}

// decrComputingPowerClaimable decrements the claimable computing power of an owner.
func (k Keeper) decrComputingPowerClaimable(ctx sdk.Context, amount uint64, owner sdk.AccAddress) {
	before := k.GetComputingPowerClaimable(ctx, owner)
	after := before - amount
	k.setComputingPowerClaimable(ctx, after, owner)
}

// setComputingPowerClaimable sets the claimable computing power of an owner.
func (k Keeper) setComputingPowerClaimable(ctx sdk.Context, amount uint64, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ComputingPowerClaimableStoreKey(owner), sdk.Uint64ToBigEndian(amount))
}

// GetComputingPowerClaimable returns the claimable computing power of an owner.
func (k Keeper) GetComputingPowerClaimable(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ComputingPowerClaimableStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
