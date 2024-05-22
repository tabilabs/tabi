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

// CreateNode defines a method for create a new node for owner.
func (k Keeper) CreateNode(
	ctx sdk.Context,
	divisionID string,
	owner sdk.AccAddress,
) (string, error) {
	division, found := k.GetDivision(ctx, divisionID)
	if !found {
		return "", errorsmod.Wrap(types.ErrDivisionNotExists, divisionID)
	}
	if division.SoldCount == division.InitialSupply {
		return "", errorsmod.Wrap(types.ErrDivisionSoldOut, divisionID)
	}

	nodeID := k.GenerateNodeID(ctx)
	if k.HasNode(ctx, nodeID) {
		return "", errorsmod.Wrap(types.ErrNodeExists, nodeID)
	}

	node := types.Node{
		Id:             nodeID,
		DivisionId:     divisionID,
		Owner:          owner.String(),
		ComputingPower: division.ComputingPowerLowerBound,
	}
	if err := k.setNode(ctx, node); err != nil {
		return "", err
	}
	k.setNodeByOwner(ctx, nodeID, owner)

	division.TotalCount += 1
	division.SoldCount += 1
	if err := k.setDivision(ctx, division); err != nil {
		return "", err
	}

	return nodeID, nil
}

// UpdateNode defines a method for updating the computing power of the specified node
func (k Keeper) UpdateNode(
	ctx sdk.Context,
	nodeID string,
	amount uint64,
	owner sdk.AccAddress,
) error {
	node, found := k.GetNode(ctx, nodeID)
	if !found {
		return errorsmod.Wrap(types.ErrNodeNotExists, nodeID)
	}

	if err := k.AuthorizeNode(ctx, nodeID, owner); err != nil {
		return err
	}

	claimable := k.GetClaimableComputingPower(ctx, owner)
	if claimable < amount {
		return errorsmod.Wrap(types.ErrInsufficientComputingPower, nodeID)
	}

	after := node.ComputingPower + amount
	currDivision, _ := k.GetDivision(ctx, node.DivisionId)
	if after > currDivision.ComputingPowerUpperBound {
		// check if we need to improve node division
		nextDivision := k.DecideDivision(ctx, after)
		node.DivisionId = nextDivision.Id
		k.incrDivisionTotalCount(ctx, nextDivision)
		k.decrDivisionTotalCount(ctx, currDivision)
	}

	// set node info
	node.ComputingPower = after
	if err := k.setNode(ctx, node); err != nil {
		return err
	}

	// set claimable power
	k.decrClaimableComputingPower(ctx, amount, owner)

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
	addr, found := k.GetNodeOwner(ctx, nodeID)
	if !found {
		return errorsmod.Wrap(types.ErrNodeNotExists, nodeID)
	}

	if !owner.Equals(addr) {
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

// GetNodesCount returns the total number of nodes
func (k Keeper) GetNodesCount(ctx sdk.Context) uint64 {
	var count uint64
	divisions := k.GetDivisions(ctx)
	for _, division := range divisions {
		count += division.TotalCount
	}
	return count
}

// GetNodesByOwner return all nodes owned by the specified owner
func (k Keeper) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node) {
	store := k.getNodeByOwnerPrefixStore(ctx, owner)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		node, has := k.GetNode(ctx, string(iterator.Key()))
		if has {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetNodeOwner returns the owner of the specified node
func (k Keeper) GetNodeOwner(ctx sdk.Context, nodeID string) (sdk.AccAddress, bool) {
	node, found := k.GetNode(ctx, nodeID)
	if !found {
		return nil, false
	}
	owner := sdk.MustAccAddressFromBech32(node.Owner)
	return owner, true
}

// setNode defines a method for setting the node
func (k Keeper) setNode(ctx sdk.Context, node types.Node) error {
	bz, err := k.cdc.Marshal(&node)
	if err != nil {
		return errorsmod.Wrap(err, "Marshal node failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.NodeStoreKey(node.Id), bz)
	return nil
}

// SetOwner defines a method for setting the owner of the specified node
func (k Keeper) setNodeByOwner(ctx sdk.Context, nodeID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeByOwnerStoreKey(owner, nodeID)
	store.Set(key, types.PlaceHolder)
}

// getNodesStoreByOwner returns the store for the nodes owned by the specified owner
func (k Keeper) getNodeByOwnerPrefixStore(ctx sdk.Context, owner sdk.AccAddress) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	key := types.NodeByOwnerPrefixStoreKey(owner)
	return prefix.NewStore(store, key)
}

// getNodesPrefixStore returns the store for the nodes
func (k Keeper) getNodesPrefixStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.NodeKey)
}

// GetNodeSequence gets the next Node sequence from the store.
func (k Keeper) GetNodeSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NodeNextSequenceKey)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

// SetNodeSequence sets the next Node sequence to the store.
func (k Keeper) SetNodeSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set(types.NodeNextSequenceKey, bz)
}
