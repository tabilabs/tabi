package keeper

import (
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// CalcNodeComputingPowerRatioOnEpoch returns the computing power ratio of a node as per its node info.
func (k Keeper) CalcNodeComputingPowerRatioOnEpoch(
	ctx sdk.Context,
	epochID uint64,
	nodeID string,
	powerOnRatio sdk.Dec,
) (sdk.Dec, error) {
	nodePower, err := k.CalcNodeComputingPowerOnEpoch(ctx, epochID, nodeID, powerOnRatio)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	sum := sdk.NewDec(int64(k.GetComputingPowerSumOnEpoch(ctx, epochID)))
	if sum.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec(), errorsmod.Wrap(types.ErrInvalidCalculation, "computing sum is zero")
	}

	return nodePower.Quo(sum), nil
}

// GetComputingPowerSumOnEpoch gets the sum of computing power of all nodes.
func (k Keeper) GetComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.ComputingPowerSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	return sdk.BigEndianToUint64(bz)
}

// incrComputingPowerSumOnEpoch increases the sum of computing power of all nodes.
// NOTE: call only after computing a node power so that by end of epoch we have the power sum of all nodes.
func (k Keeper) incrComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64, amount uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.ComputingPowerSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	sum := sdk.BigEndianToUint64(bz)
	sum += amount
	store.Set(key, sdk.Uint64ToBigEndian(sum))
}

// CalcNodeComputingPowerOnEpoch returns the computing power of a node as per its node info.
func (k Keeper) CalcNodeComputingPowerOnEpoch(
	ctx sdk.Context,
	epochID uint64,
	nodeID string,
	powerOnRatio sdk.Dec,
) (sdk.Dec, error) {
	if powerOnRatio.Equal(sdk.ZeroDec()) {
		powerOnRatio = sdk.OneDec()
	}

	basePower := sdk.NewDec(int64(k.GetNodeBaseComputingPower(ctx, nodeID)))
	pledgeRatio, err := k.CalcNodePledgeRatioOnEpoch(ctx, epochID, nodeID)
	if err != nil {
		return sdk.ZeroDec(), err

	}

	// exponent = pledge_ratio / 0.5
	exponentiation, err := pledgeRatio.Mul(sdk.NewDec(2)).Float64()
	if err != nil {
		return sdk.ZeroDec(), err
	}
	exponentiated, err := sdk.NewDecFromStr(fmt.Sprintf("%f", math.Exp(exponentiation)))
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return basePower.Mul(exponentiated).Mul(powerOnRatio), nil
}

// setComputingPowerByNode returns the computing power of a node as per its node info.
func (k Keeper) setNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string, power sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	store.Set(key, []byte(power.String()))
}

// GetNodeComputingPowerOnEpoch returns the computing power of a node as per its node info.
func (k Keeper) GetNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	return sdk.NewDecFromStr(string(bz))
}

// GetNodeBaseComputingPower returns the base computing power of a node as per its node info.
func (k Keeper) GetNodeBaseComputingPower(ctx sdk.Context, nodeID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NodeStoreKey(nodeID))
	if bz == nil {
		return 0
	}

	var node types.Node
	k.cdc.MustUnmarshal(bz, &node)
	return node.ComputingPower
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
	store.Set(types.NodeClaimableComputingPowerStoreKey(owner), sdk.Uint64ToBigEndian(amount))
}

// GetComputingPowerClaimable returns the claimable computing power of an owner.
func (k Keeper) GetComputingPowerClaimable(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NodeClaimableComputingPowerStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
