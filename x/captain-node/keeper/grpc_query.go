package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the staking parameters
func (k Keeper) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// Owner queries the owner of the node
func (k Keeper) Owner(goCtx context.Context, request *types.QueryOwnerRequest) (*types.QueryOwnerResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := k.GetOwner(ctx, request.Id)

	return &types.QueryOwnerResponse{Owner: owner.String()}, nil
}

// Supply queries the number of Node from the given division
func (k Keeper) Supply(goCtx context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	supply := k.GetDivisionTotalSupply(ctx, request.DivisionId)
	return &types.QuerySupplyResponse{Amount: supply}, nil
}

// Division queries an node division by its ID
func (k Keeper) Division(goCtx context.Context, request *types.QueryDivisionRequest) (*types.QueryDivisionResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	division, found := k.GetDivision(ctx, request.DivisionId)
	if !found {
		return nil, types.ErrDivisionNotExists.Wrapf("division not found: %s", request.DivisionId)

	}
	return &types.QueryDivisionResponse{Division: &division}, nil
}

// Divisions queries all Node divisions
func (k Keeper) Divisions(goCtx context.Context, request *types.QueryDivisionsRequest) (*types.QueryDivisionsResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)

	divisionStore := prefix.NewStore(store, types.DivisionKey)

	var divisions []*types.Division
	pageRes, err := query.Paginate(divisionStore, request.Pagination, func(key []byte, value []byte) error {
		var division types.Division
		if err := k.cdc.Unmarshal(value, &division); err != nil {
			return err
		}
		divisions = append(divisions, &division)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryDivisionsResponse{
		Divisions:  divisions,
		Pagination: pageRes,
	}, nil
}

// Node queries an Node based on its id.
func (k Keeper) Node(goCtx context.Context, request *types.QueryNodeRequest) (*types.QueryNodeResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	node, found := k.GetNode(ctx, request.Id)
	if !found {
		return nil, types.ErrNodeNotExists.Wrapf("not found node: %s", request.Id)
	}
	return &types.QueryNodeResponse{
		Node: &node,
	}, nil
}

// Nodes queries all node of a given owner
func (k Keeper) Nodes(goCtx context.Context, request *types.QueryNodesRequest) (*types.QueryNodesResponse, error) {
	if request == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	var err error
	var owner sdk.AccAddress
	if len(request.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, err
		}
	}

	var nodes []*types.Node
	var pageRes *query.PageResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	switch {
	case len(request.Owner) > 0:

		if pageRes, err = query.Paginate(k.getNodesStoreByOwner(ctx, owner), request.Pagination, func(key []byte, value []byte) error {
			node, has := k.GetNode(ctx, string(key))
			if has {
				nodes = append(nodes, &node)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	default:
		// return all nodes
		nodeStore := k.getNodesStore(ctx)
		if pageRes, err = query.Paginate(nodeStore, request.Pagination, func(_ []byte, value []byte) error {
			var node types.Node
			if err := k.cdc.Unmarshal(value, &node); err != nil {
				return err
			}
			nodes = append(nodes, &node)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &types.QueryNodesResponse{
		Pagination: pageRes,
		Nodes:      nodes,
	}, nil
}

func (k Keeper) SaleLevel(goCtx context.Context, _ *types.QuerySaleLevelRequest) (*types.QuerySaleLevelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	level := k.GetSaleLevel(ctx)
	return &types.QuerySaleLevelResponse{SaleLevel: level}, nil
}

func (k Keeper) Callers(goCtx context.Context, _ *types.QueryCallersRequest) (*types.QueryCallersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	callers := k.GetCallers(ctx)
	return &types.QueryCallersResponse{Callers: callers}, nil
}
