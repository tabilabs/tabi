package keeper

import (
	"fmt"

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
func (k Keeper) CommitComputingPower(
	ctx sdk.Context,
	claimableComputingPowers []types.ClaimableComputingPower,
) []sdk.Event {
	events := make([]sdk.Event, 0)
	for _, claimableComputingPower := range claimableComputingPowers {
		receiver := sdk.MustAccAddressFromBech32(claimableComputingPower.Owner)
		before := k.GetComputingPowerClaimable(ctx, receiver)
		k.incrExtractableComputingPower(ctx, receiver, claimableComputingPower.Amount)
		after := k.GetComputingPowerClaimable(ctx, receiver)
		events = append(
			events,
			sdk.NewEvent(
				types.EventCommitComputingPower,
				sdk.NewAttribute(types.AttributeKeyOwner, claimableComputingPower.Owner),
				sdk.NewAttribute(types.AttributeKeyComputingPowerBefore, fmt.Sprintf("%d", before)),
				sdk.NewAttribute(types.AttributeKeyComputingPowerAfter, fmt.Sprintf("%d", after)),
			),
		)
	}

	return nil
}

func (k Keeper) incrExtractableComputingPower(
	ctx sdk.Context,
	owner sdk.AccAddress,
	amount uint64,
) {
	store := ctx.KVStore(k.storeKey)
	result := k.GetComputingPowerClaimable(ctx, owner) + amount
	store.Set(types.ComputingPowerClaimableStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) decrComputingPowerClaimable(
	ctx sdk.Context,
	owner sdk.AccAddress,
	amount uint64,
) {
	store := ctx.KVStore(k.storeKey)
	result := k.GetComputingPowerClaimable(ctx, owner) - amount
	store.Set(types.ComputingPowerClaimableStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) GetComputingPowerClaimable(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ComputingPowerClaimableStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
