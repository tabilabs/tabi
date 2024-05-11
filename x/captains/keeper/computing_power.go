package keeper

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// ComputingPowerSumOnEpoch functions

// GetComputingPowerSumOnEpoch returns the sum of computing power of all nodes.
func (k Keeper) GetComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.ComputingPowerSumOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}

	sum, _ := sdk.NewDecFromStr(string(bz))
	return sum
}

// setComputingPowerSumOnEpoch sets the sum of computing power of all nodes.
func (k Keeper) setComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.ComputingPowerSumOnEpochStoreKey(epochID)
	store.Set(key, []byte(amount.String()))
}

// incrComputingPowerSumOnEpoch increases the sum of computing power of all nodes.
// NOTE: call only after computing a node power so that by end of epoch we have the power sum of all nodes.
func (k Keeper) incrComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	sum := k.GetComputingPowerSumOnEpoch(ctx, epochID)
	sum = sum.Add(amount)
	k.setComputingPowerSumOnEpoch(ctx, epochID, sum)
}

// DelComputingPowerSumOnEpoch deletes the sum of computing power of all nodes.
func (k Keeper) DelComputingPowerSumOnEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.ComputingPowerSumOnEpochStoreKey(epochID)
	store.Delete(key)
}

// NodeComputingPowerOnEpoch functions

// CalcNodeComputingPowerOnEpoch returns the computing power of a node as per its node info.
func (k Keeper) CalcNodeComputingPowerOnEpoch(
	ctx sdk.Context,
	epochID uint64,
	nodeID string,
	powerOnRatio sdk.Dec,
) sdk.Dec {
	basePower := sdk.NewDec(int64(k.GetNodeBaseComputingPower(ctx, nodeID)))
	pledgeRatio := k.CalcNodePledgeRatioOnEpoch(ctx, epochID, nodeID)
	exponentiation := pledgeRatio.Mul(sdk.NewDec(20)).Quo(sdk.NewDec(3)).MustFloat64()
	exponentiated := sdk.MustNewDecFromStr(fmt.Sprintf("%f", math.Exp(exponentiation)))
	return basePower.Mul(exponentiated).Mul(powerOnRatio)
}

// setComputingPowerByNode returns the computing power of a node as per its node info.
func (k Keeper) setNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string, power sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	store.Set(key, []byte(power.String()))
}

// delNodeComputingPowerOnEpoch deletes the computing power of a node as per its node info.
func (k Keeper) delNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	store.Delete(key)
}

// GetNodeComputingPowerOnEpoch returns the computing power of a node as per its node info.
func (k Keeper) GetNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
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

// ComputingPowerClaimable functions

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

// GetComputingPowersClaimable returns all claimable computing powers.
func (k Keeper) GetComputingPowersClaimable(ctx sdk.Context) []types.ClaimableComputingPower {
	var claimableComputingPowers []types.ClaimableComputingPower

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ComputingPowerClaimableKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.BigEndianToUint64(iterator.Value())
		owner := string(iterator.Key())

		power := types.ClaimableComputingPower{
			Amount: amount,
			Owner:  owner,
		}
		claimableComputingPowers = append(claimableComputingPowers, power)
	}

	return claimableComputingPowers
}
