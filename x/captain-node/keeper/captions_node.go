package keeper

import (
	"crypto/sha256"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

const nodeIdPrefix = "node-%d"

// CreateNode defines a method for create a new node
func (k Keeper) CreateNode(ctx sdk.Context, node types.Node, receiver sdk.AccAddress) error {
	// Check Division exists
	if !k.HasDivision(ctx, node.DivisionId) {
		return errorsmod.Wrap(types.ErrDivisionNotExists, node.DivisionId)
	}

	// Check if the Division is sold out
	if k.IsDivisionSoldOut(ctx, node.DivisionId) {
		return errorsmod.Wrap(types.ErrDivisionSoldOut, node.DivisionId)
	}

	// Check if the node already exists
	// node.id is unique in the global scope
	if k.HasNode(ctx, node.Id) {
		return errorsmod.Wrap(types.ErrNodeExists, node.Id)
	}

	// Check user-holding quantity
	if k.isUserHoldingQuantityExceeded(ctx, receiver) {
		return errorsmod.Wrap(types.ErrUserHoldingQuantityExceeded, receiver.String())
	}

	// Set the node
	k.setNode(ctx, node)
	// Set the owner
	k.setOwner(ctx, node.Id, receiver)
	// Set user-holding quantity
	k.incrUserHoldingQuantity(ctx, receiver)
	// Set the division total supply
	k.incrDivisionTotalSupply(ctx, node.DivisionId)

	bz, err := k.cdc.Marshal(&node)
	if err != nil {
		return errorsmod.Wrap(err, "Marshal node failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NodeStoreKey(node.Id), bz)
	return nil
}

// GenerateNodeID defines a method for generating a new node id
func (k Keeper) GenerateNodeID(ctx sdk.Context) string {
	sequence := k.GetNodeSequence(ctx)
	nodeID := fmt.Sprintf(nodeIdPrefix, sequence)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(nodeID)))
	k.SetNodeSequence(ctx, sequence+1)
	return hash
}

// HasNode defines a method for checking the existence of a node
func (k Keeper) HasNode(ctx sdk.Context, nodeID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.NodeStoreKey(nodeID))
}

// GetNode defines a method for returning the node information of the specified id
func (k Keeper) GetNode(ctx sdk.Context, nodeID string) (types.Node, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeStoreKey(nodeID)

	bz := store.Get(key)

	var node types.Node
	if len(bz) == 0 {
		return node, false
	}
	k.cdc.MustUnmarshal(bz, &node)
	return node, true
}

// GetNodesByOwner return all nodes owned by the specified owner
func (k Keeper) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node) {
	ownerStore := k.getNodeStoreByOwner(ctx, owner)
	iterator := ownerStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		node, has := k.GetNode(ctx, string(iterator.Key()))
		if has {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// GetOwner returns the owner of the specified node
func (k Keeper) GetOwner(ctx sdk.Context, nodeID string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.OwnerStoreKey(nodeID)
	bz := store.Get(key)
	return sdk.AccAddress(bz)
}

// GetUserHoldingQuantity returns the amount of nodes owned by the specified owner
func (k Keeper) GetUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	return k.getUserHoldingQuantity(ctx, owner)
}

func (k Keeper) setNode(ctx sdk.Context, node types.Node) {
	bz := k.cdc.MustMarshal(&node)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NodeStoreKey(node.Id), bz)
}

// SetOwner defines a method for setting the owner of the specified node
// and setting the owner of the specified node
func (k Keeper) setOwner(ctx sdk.Context, nodeID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OwnerStoreKey(nodeID), owner.Bytes())

	ownerStore := k.getNodeStoreByOwner(ctx, owner)
	ownerStore.Set([]byte(nodeID), types.Placeholder)
}

func (k Keeper) isUserHoldingQuantityExceeded(ctx sdk.Context, owner sdk.AccAddress) bool {
	params := k.GetParams(ctx)
	maximumNumberOfHoldings := params.MaximumNumberOfHoldings
	if k.getUserHoldingQuantity(ctx, owner) >= maximumNumberOfHoldings {
		return false
	}
	return true
}

func (k Keeper) getUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OwnerHoldingTotalSupplyStoreKey(owner))
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) incrUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress) {
	supply := k.getUserHoldingQuantity(ctx, owner) + 1
	k.updateUserHoldingQuantity(ctx, owner, supply)
}

func (k Keeper) updateUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress, supply uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OwnerHoldingTotalSupplyStoreKey(owner), sdk.Uint64ToBigEndian(supply))
}

func (k Keeper) getNodeStoreByOwner(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.NodeByOwnerStoreKey(owner))
}

// GetNodeSequence gets the next Node sequence from the store.
func (k Keeper) GetNodeSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextNodeSequence))
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNodeSequence sets the next Node sequence to the store.
func (k Keeper) SetNodeSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextNodeSequence), bz)
}
