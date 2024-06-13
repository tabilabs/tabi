package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// Emission handles all emission calculations.

// CalcBaseEpochEmission returns the base emission reward for an epoch.
func (k Keeper) CalcBaseEpochEmission(ctx sdk.Context) sdk.Dec {
	tech := k.CalcTechProgressCoefficient(ctx)
	halving := k.GetHalvingEraCoefficient(ctx)
	cc := k.GetCaptainsConstant(ctx)
	return tech.Mul(halving).Mul(cc).Mul(sdk.NewDec(1e18))
}

// HasEpochEmission returns if the emission reward for an epoch exists.
func (k Keeper) HasEpochEmission(ctx sdk.Context, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEmissionStoreKey(epochID)
	return store.Has(key)
}

// GetEpochEmission returns the emission reward for an epoch.
func (k Keeper) GetEpochEmission(ctx sdk.Context, epochID uint64) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEmissionStoreKey(epochID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// CalcEpochEmission returns the total emission reward for an epoch.
func (k Keeper) CalcEpochEmission(ctx sdk.Context, epochID uint64, globalOperationRatio sdk.Dec) sdk.Dec {
	base := k.CalcBaseEpochEmission(ctx)
	pledgeRatio := k.CalcGlobalPledgeRatio(ctx, epochID)
	return base.Mul(pledgeRatio).Mul(globalOperationRatio)
}

// setEpochEmission sets the emission sum for an epoch.
func (k Keeper) setEpochEmission(ctx sdk.Context, epochID uint64, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEmissionStoreKey(epochID)
	store.Set(key, []byte(amount.String()))
}

// delEpochEmission deletes the emission sum for an epoch.
func (k Keeper) delEpochEmission(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEmissionStoreKey(epochID)
	store.Delete(key)
}

// GetGlobalClaimedEmission returns global claimed emission.
func (k Keeper) GetGlobalClaimedEmission(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalClaimedEmissionKey
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// SetGlobalClaimedEmission sets global claimed emission.
func (k Keeper) SetGlobalClaimedEmission(ctx sdk.Context, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.GlobalClaimedEmissionKey
	store.Set(key, []byte(amount.String()))
}

// incrGlobalClaimedEmission increases global claimed emission by amount.
func (k Keeper) incrGlobalClaimedEmission(ctx sdk.Context, amount sdk.Dec) sdk.Dec {
	emission := k.GetGlobalClaimedEmission(ctx)
	emission = emission.Add(amount)
	k.SetGlobalClaimedEmission(ctx, emission)
	return emission
}

// CalcAndSetNodeCumulativeEmissionByEpoch returns the historical emission for a node at the end of an epoch.
// NOTE: this function set the historical emission by the end of epoch(t) and removes that of epoch(t-1).
func (k Keeper) CalcAndSetNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	res := k.CalcNodeCumulativeEmissionByEpoch(ctx, epochID, nodeID)

	// no need to set the emission for the first epoch
	if epochID == 0 {
		return res
	}

	k.setNodeCumulativeEmissionByEpoch(ctx, epochID, nodeID, res)
	k.delNodeCumulativeEmissionByEpoch(ctx, epochID-1, nodeID)
	k.delNodeComputingPowerOnEpoch(ctx, epochID-1, nodeID)

	return res
}

// CalcNodeCumulativeEmissionByEpoch returns the historical emission for a node at the end of an epoch.
// NOTE: this func is only used to query and it won't set or prune state data.
func (k Keeper) CalcNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	if epochID == 0 {
		return sdk.ZeroDec()
	}

	// NOTE: we may have already set this value when user claimed the rewards before report handles it.
	historyEmission := k.GetNodeCumulativeEmissionByEpoch(ctx, epochID, nodeID)
	if !historyEmission.Equal(sdk.ZeroDec()) {
		return historyEmission
	}

	prevHistoryEmission := k.GetNodeCumulativeEmissionByEpoch(ctx, epochID-1, nodeID)
	epochEmission := k.CalcNodeEmissionOnEpoch(ctx, epochID, nodeID)

	return epochEmission.Add(prevHistoryEmission)
}

// CalcNodeEmissionOnEpoch returns the emission for a node at the end of an epoch.
func (k Keeper) CalcNodeEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	emission := k.GetEpochEmission(ctx, epochID)
	power := k.GetNodeComputingPowerOnEpoch(ctx, epochID, nodeID)
	powerSum := k.GetGlobalComputingPowerOnEpoch(ctx, epochID)
	// NOTE: zero power sum is never expected; it panics if it happens.
	return emission.Mul(power).Quo(powerSum)
}

// HasNodeCumulativeEmissionByEpoch returns if the historical emission for a node at the end of an epoch exists.
func (k Keeper) HasNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeCumulativeEmissionByEpochStoreKey(epochID, nodeID)
	return store.Has(key)
}

// GetNodeCumulativeEmissionByEpoch returns the historical emission for a node at the end of an epoch.
func (k Keeper) GetNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeCumulativeEmissionByEpochStoreKey(epochID, nodeID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// setNodeCumulativeEmissionByEpoch sets the historical emission for a node at the end of an epoch.
func (k Keeper) setNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeCumulativeEmissionByEpochStoreKey(epochID, nodeID)
	store.Set(key, []byte(amount.String()))
}

// delNodeCumulativeEmissionByEpoch deletes the historical emission for a node at the end of an epoch.
func (k Keeper) delNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeCumulativeEmissionByEpochStoreKey(epochID, nodeID)
	store.Delete(key)
}

// UpdateGlobalAndNodeClaimedEmission updates node's claimed emission and incr global claimed emission.
// NOTE: call this function only after claiming the rewards.
func (k Keeper) UpdateGlobalAndNodeClaimedEmission(ctx sdk.Context, nodeID string) error {
	epoch := k.GetCurrentEpoch(ctx)
	before := k.GetNodeClaimedEmission(ctx, nodeID)
	after := k.CalcNodeCumulativeEmissionByEpoch(ctx, epoch-1, nodeID)
	k.SetNodeClaimedEmission(ctx, nodeID, after)
	k.incrGlobalClaimedEmission(ctx, after.Sub(before))
	return nil
}

// SetNodeClaimedEmission sets the historical emission the last time user claimed.
func (k Keeper) SetNodeClaimedEmission(ctx sdk.Context, nodeID string, amount sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeClaimedEmissionStoreKey(nodeID)
	store.Set(key, []byte(amount.String()))
}

// GetNodeClaimedEmission returns the claimed emission of a node.
func (k Keeper) GetNodeClaimedEmission(ctx sdk.Context, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeClaimedEmissionStoreKey(nodeID)
	bz := store.Get(key)
	if len(bz) == 0 {
		return sdk.ZeroDec()
	}
	return sdk.MustNewDecFromStr(string(bz))
}

// GetOwnerClaimedEmission returns the historical emission the last time user claimed.
func (k Keeper) GetOwnerClaimedEmission(ctx sdk.Context, owner sdk.AccAddress) sdk.Dec {
	nodes := k.GetNodesByOwner(ctx, owner)
	total := sdk.ZeroDec()
	for _, node := range nodes {
		total = total.Add(k.GetNodeClaimedEmission(ctx, node.Id))
	}
	return total
}

// Genesis State Export/Import Helpers

// GetEpochesEmission sets the emission reward for an epoch.
func (k Keeper) GetEpochesEmission(ctx sdk.Context) []types.EpochEmission {
	var epochesEmission []types.EpochEmission
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EpochEmissionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var epochEmission types.EpochEmission
		epochEmission.EpochId = types.SplitEpochFromStoreKey(types.EpochEmissionKey, iterator.Key())
		epochEmission.Emission = k.GetEpochEmission(ctx, epochEmission.EpochId)
		epochesEmission = append(epochesEmission, epochEmission)
	}
	return epochesEmission
}

// GetNodesClaimedEmission returns the claimed emission of all nodes.
func (k Keeper) GetNodesClaimedEmission(ctx sdk.Context) []types.NodeClaimedEmission {
	var nodesClaimedEmission []types.NodeClaimedEmission
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NodeClaimedEmissionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var nodeClaimedEmission types.NodeClaimedEmission
		nodeClaimedEmission.NodeId = types.SplitStrFromStoreKey(types.NodeClaimedEmissionKey, iterator.Key())
		nodeClaimedEmission.Emission = k.GetNodeClaimedEmission(ctx, nodeClaimedEmission.NodeId)
		nodesClaimedEmission = append(nodesClaimedEmission, nodeClaimedEmission)
	}
	return nodesClaimedEmission
}

// GetNodesCumulativeEmission returns the cumulative emission of all nodes.
func (k Keeper) GetNodesCumulativeEmission(ctx sdk.Context) []types.NodeCumulativeEmission {
	var nodesCumulativeEmission []types.NodeCumulativeEmission
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NodeCumulativeEmissionByEpochKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var nodeCumulativeEmission types.NodeCumulativeEmission
		epochId, nodeId := types.SplitEpochAndStrFromStoreKey(types.NodeCumulativeEmissionByEpochKey, iterator.Key())
		nodeCumulativeEmission.EpochId = epochId
		nodeCumulativeEmission.NodeId = nodeId
		nodeCumulativeEmission.Emission = sdk.MustNewDecFromStr(string(iterator.Value()))
		nodesCumulativeEmission = append(nodesCumulativeEmission, nodeCumulativeEmission)
	}
	return nodesCumulativeEmission
}
