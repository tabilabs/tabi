package keeper

import (
	"crypto/sha256"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
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

	return nil
}

func (k Keeper) UpdateNode(
	ctx sdk.Context,
	nodeID string,
	computingPower uint64,
	owner sdk.AccAddress,
) error {
	// Check if the node exists
	if !k.HasNode(ctx, nodeID) {
		return errorsmod.Wrap(types.ErrNodeNotExists, nodeID)
	}

	// Check if owner of the node is the sender
	if err := k.AuthorizeNode(ctx, nodeID, owner); err != nil {
		return errorsmod.Wrap(types.ErrUnauthorized, owner.String())
	}

	// Check if the node has enough extractable computing power
	if k.getExtractableComputingPower(ctx, owner) < computingPower {
		return errorsmod.Wrap(types.ErrInsufficientExperience, nodeID)
	}

	// Update the node
	node, _ := k.GetNode(ctx, nodeID)
	currentDivision, _ := k.GetDivision(ctx, node.DivisionId)

	node.ComputingPower += computingPower

	// Check if the node has enough experience to be promoted to the next division
	if node.ComputingPower > currentDivision.ComputingPowerUpperBound {
		divisions := k.GetDivisions(ctx)
		for _, division := range divisions {
			// the node should be promoted to the next division
			if node.ComputingPower <= division.ComputingPowerUpperBound && node.ComputingPower >= division.ComputingPowerLowerBound {
				node.DivisionId = division.Id
				break
			}
		}
	}

	// Set the node
	k.setNode(ctx, node)
	// Set the experience
	k.decrExtractableComputingPower(ctx, owner, computingPower)
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

// GetLastNodeID defines a method for returning the last node id
func (k Keeper) GetLastNodeID(ctx sdk.Context) string {
	sequence := k.GetNodeSequence(ctx)
	return fmt.Sprintf(nodeIdPrefix, sequence-1)
}

// HasNode defines a method for checking the existence of a node
func (k Keeper) HasNode(ctx sdk.Context, nodeID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.NodeStoreKey(nodeID))
}

// AuthorizeNode defines a method for checking if the sender is the owner of the given node
func (k Keeper) AuthorizeNode(ctx sdk.Context, nodeID string, owner sdk.AccAddress) error {
	if !owner.Equals(k.GetOwner(ctx, nodeID)) {
		return errorsmod.Wrap(types.ErrUnauthorized, owner.String())
	}
	return nil
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

// GetNodes returns all nodes
func (k Keeper) GetNodes(ctx sdk.Context) (nodes []types.Node) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NodeKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var node types.Node
		k.cdc.MustUnmarshal(iterator.Value(), &node)
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNodesByOwner return all nodes owned by the specified owner
func (k Keeper) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node) {
	ownerStore := k.getNodesStoreByOwner(ctx, owner)
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

// GetOwners returns all owners
func (k Keeper) GetOwners(ctx sdk.Context) (owners []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.NodeByOwnerKey)
	defer iterator.Close()
	// Need to filter duplicates from the iterator
	ownerMap := make(map[string]struct{})
	for ; iterator.Valid(); iterator.Next() {
		// iterator.Key() = <owner><Delimiter><nodeID>
		// strings.Split(string(iterator.Key()), types.Delimiter)[0] = <owner>
		owner, _ := types.ParseNodeByOwnerStoreKey(iterator.Key())
		ownerMap[owner.String()] = struct{}{}
	}
	for owner := range ownerMap {
		accAddr, _ := sdk.AccAddressFromBech32(owner)
		owners = append(owners, accAddr)
	}
	return owners

}

// GetUserHoldingQuantity returns the amount of nodes owned by the specified owner
func (k Keeper) GetUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	return k.getUserHoldingQuantity(ctx, owner)
}

// setNode defines a method for setting the node
// only can have one owner
func (k Keeper) setNode(ctx sdk.Context, node types.Node) {
	bz := k.cdc.MustMarshal(&node)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NodeStoreKey(node.Id), bz)
	fmt.Println("store.Get(types.NodeStoreKey(node.Id)): ", store.Get(types.NodeStoreKey(node.Id)))
}

// SetOwner defines a method for setting the owner of the specified node
// and setting the owner of the specified node
// user can have multiple nodes
func (k Keeper) setOwner(ctx sdk.Context, nodeID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// nodeID -> owner
	key := types.OwnerStoreKey(nodeID)
	store.Set(key, owner.Bytes())

	// Set the node by owner
	// 0x02<owner><Delimiter>
	ownerStore := k.getNodesStoreByOwner(ctx, owner)

	// 0x02<owner><Delimiter><nodeID> -> Placeholder
	ownerStore.Set([]byte(nodeID), types.Placeholder)
}

func (k Keeper) getNodesStoreByOwner(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeByOwnerStoreKey(owner) // 0x02<owner><Delimiter>
	return prefix.NewStore(store, key)
}

// isUserHoldingQuantityExceeded checks if the user holding quantity exceeded
// if exceeded, return true
func (k Keeper) isUserHoldingQuantityExceeded(ctx sdk.Context, owner sdk.AccAddress) bool {
	params := k.GetParams(ctx)
	maximumHoldingAmount := params.MaximumHoldingAmount
	if k.getUserHoldingQuantity(ctx, owner) >= maximumHoldingAmount {
		return true
	}
	return false
}

func (k Keeper) getUserHoldingQuantity(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OwnerHoldingTotalSupplyStoreKey(owner))
	if bz == nil {
		return 0
	}
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

func (k Keeper) getNodesStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.NodeKey)
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
