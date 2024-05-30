package keeper

import (
	"fmt"
	"math"

	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetGlobalComputingPowerOnEpoch returns the sum of computing power of all nodes.
func (k Keeper) GetGlobalComputingPowerOnEpoch(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalComputingPowerOnEpochStoreKey(epochID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// setGlobalComputingPowerOnEpoch sets the sum of computing power of all nodes.
func (k Keeper) setGlobalComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalComputingPowerOnEpochStoreKey(epochID)
	store.Set(key, []byte(amount.String()))
}

// incrGlobalComputingPowerOnEpoch increases the sum of computing power of all nodes.
func (k Keeper) incrGlobalComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	sum := k.GetGlobalComputingPowerOnEpoch(ctx, epochID)
	sum = sum.Add(amount)
	k.setGlobalComputingPowerOnEpoch(ctx, epochID, sum)
}

// delGlobalComputingPowerOnEpoch deletes the sum of computing power of all nodes.
func (k Keeper) delGlobalComputingPowerOnEpoch(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalComputingPowerOnEpochStoreKey(epochID)
	store.Delete(key)
}

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

// HasNodeComputingPowerOnEpoch checks if the node has computing power on the epoch.
func (k Keeper) HasNodeComputingPowerOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeComputingPowerOnEpochStoreKey(epochID, nodeID)
	return store.Has(key)
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

// CommitComputingPower commits the pending computing power.
func (k Keeper) CommitComputingPower(ctx sdk.Context, amount uint64, owner sdk.AccAddress) (uint64, uint64) {
	before := k.GetClaimableComputingPower(ctx, owner)
	after := before + amount
	k.setClaimableComputingPower(ctx, after, owner)
	return before, after
}

// decrClaimableComputingPower decrements the claimable computing power of an owner.
func (k Keeper) decrClaimableComputingPower(ctx sdk.Context, amount uint64, owner sdk.AccAddress) {
	power := k.GetClaimableComputingPower(ctx, owner)
	power -= amount
	k.setClaimableComputingPower(ctx, power, owner)
}

// setClaimableComputingPower sets the claimable computing power of an owner.
func (k Keeper) setClaimableComputingPower(ctx sdk.Context, amount uint64, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ClaimableComputingPowerStoreKey(owner), sdk.Uint64ToBigEndian(amount))
}

// GetClaimableComputingPower returns the claimable computing power of an owner.
func (k Keeper) GetClaimableComputingPower(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ClaimableComputingPowerStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// Genesis Export/Import Helpers

// GetClaimableComputingPowers returns all claimable computing powers.
func (k Keeper) GetClaimableComputingPowers(ctx sdk.Context) []types.ClaimableComputingPower {
	var powers []types.ClaimableComputingPower
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClaimableComputingPowerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var power types.ClaimableComputingPower
		accAddr := types.SplitStrFromStoreKey(types.ClaimableComputingPowerKey, iterator.Key())
		power.Owner = sdk.AccAddress(accAddr).String()
		power.Amount = sdk.BigEndianToUint64(iterator.Value())
		powers = append(powers, power)
	}

	return powers
}

// GetGlobalsComputingPower returns all global computing power.
func (k Keeper) GetGlobalsComputingPower(ctx sdk.Context) []types.GlobalComputingPower {
	var powers []types.GlobalComputingPower
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GlobalComputingPowerOnEpochKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var power types.GlobalComputingPower
		power.EpochId = types.SplitEpochFromStoreKey(types.GlobalComputingPowerOnEpochKey, iterator.Key())
		power.Amount = k.GetGlobalComputingPowerOnEpoch(ctx, power.EpochId)
		powers = append(powers, power)
	}
	return powers
}

// GetNodesComputingPower returns all nodes computing power.
func (k Keeper) GetNodesComputingPower(ctx sdk.Context) []types.NodesComputingPower {
	var powers []types.NodesComputingPower
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NodeComputingPowerOnEpochKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var power types.NodesComputingPower
		nodeId, epochId := types.SplitNodeAndEpochFromStoreKey(types.NodeComputingPowerOnEpochKey, iterator.Key())
		power.EpochId = epochId
		power.NodeId = nodeId
		power.Amount = k.GetNodeComputingPowerOnEpoch(ctx, power.EpochId, nodeId)
		powers = append(powers, power)
	}
	return powers
}
