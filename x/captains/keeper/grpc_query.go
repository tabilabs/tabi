package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

type Querier struct {
	*Keeper
}

// NewQuerierImpl returns an implementation of the captains QueryServer interface.
func NewQuerierImpl(k *Keeper) types.QueryServer {
	return &Querier{k}
}

var _ types.QueryServer = Querier{}

// Params queries the captains module parameters
func (q Querier) Params(
	goCtx context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Node queries a Node.
func (q Querier) Node(
	goCtx context.Context,
	request *types.QueryNodeRequest,
) (*types.QueryNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	node, found := q.GetNode(ctx, request.NodeId)
	if !found {
		return nil, types.ErrNodeNotExists.Wrapf("not found node: %s", request.NodeId)
	}

	return &types.QueryNodeResponse{
		Node: &node,
	}, nil
}

// Nodes queries all node of a given owner
func (q Querier) Nodes(
	goCtx context.Context,
	request *types.QueryNodesRequest,
) (*types.QueryNodesResponse, error) {
	var err error
	var owner sdk.AccAddress
	if len(request.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(request.Owner)
		if err != nil {
			return nil, err
		}
	}

	var nodes []types.Node
	var pageRes *query.PageResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	switch {
	case len(request.Owner) > 0:
		if pageRes, err = query.Paginate(q.getNodeByOwnerPrefixStore(ctx, owner), request.Pagination,
			func(key []byte, value []byte) error {
				node, found := q.GetNode(ctx, string(key))
				if found {
					nodes = append(nodes, node)
				}
				return nil
			}); err != nil {
			return nil, err
		}
	default:
		nodeStore := q.getNodesPrefixStore(ctx)
		if pageRes, err = query.Paginate(nodeStore, request.Pagination,
			func(_ []byte, value []byte) error {
				var node types.Node
				if err := q.cdc.Unmarshal(value, &node); err != nil {
					return err
				}
				nodes = append(nodes, node)
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

// NodeLastEpochInfo queries the last epoch info of a node
func (q Querier) NodeLastEpochInfo(
	goCtx context.Context,
	request *types.QueryNodeLastEpochInfoRequest,
) (*types.QueryNodeLastEpochInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	node, found := q.GetNode(ctx, request.NodeId)
	if !found {
		return nil, types.ErrNodeNotExists.Wrapf("not found node: %s", request.NodeId)
	}

	prevEpoch := q.Keeper.GetCurrentEpoch(ctx) - 1
	emission := q.Keeper.CalcNodeEmissionOnEpoch(ctx, prevEpoch, node.Id).String()
	historical := q.Keeper.CalcNodeHistoricalEmissionOnEpoch(ctx, prevEpoch, node.Id).String()
	ratio := q.Keeper.CalcNodePledgeRatioOnEpoch(ctx, prevEpoch, node.Id).String()

	return &types.QueryNodeLastEpochInfoResponse{
		Epoch:              prevEpoch,
		LastEpochEmission:  emission,
		HistoricalEmission: historical,
		PledgeRatio:        ratio,
	}, nil
}

// Division queries an node division by its ID
func (q Querier) Division(
	goCtx context.Context,
	request *types.QueryDivisionRequest,
) (*types.QueryDivisionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	division, found := q.GetDivision(ctx, request.DivisionId)
	if !found {
		return nil, types.ErrDivisionNotExists.Wrapf("division not found: %s", request.DivisionId)

	}
	return &types.QueryDivisionResponse{Division: &division}, nil
}

// Divisions queries all Node divisions
func (q Querier) Divisions(
	goCtx context.Context,
	request *types.QueryDivisionsRequest,
) (*types.QueryDivisionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(q.storeKey)

	divisionStore := prefix.NewStore(store, types.DivisionKey)

	var divisions []types.Division
	pageRes, err := query.Paginate(divisionStore, request.Pagination, func(key []byte, value []byte) error {
		var division types.Division
		if err := q.cdc.Unmarshal(value, &division); err != nil {
			return err
		}
		divisions = append(divisions, division)
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

// Supply queries the number of Node from the given division
func (q Querier) Supply(
	goCtx context.Context,
	request *types.QuerySupplyRequest,
) (*types.QuerySupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	division, found := q.GetDivision(ctx, request.DivisionId)
	if !found {
		return nil, types.ErrDivisionNotExists.Wrapf("division not found: %s", request.DivisionId)
	}
	return &types.QuerySupplyResponse{Amount: division.SoldCount}, nil
}

// SaleLevel queries the current sale level
func (q Querier) SaleLevel(
	goCtx context.Context,
	_ *types.QuerySaleLevelRequest,
) (*types.QuerySaleLevelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	level := q.GetSaleLevel(ctx)
	return &types.QuerySaleLevelResponse{SaleLevel: level}, nil
}

// AuthorizedMembers queries the list of authorized members
func (q Querier) AuthorizedMembers(
	goCtx context.Context,
	_ *types.QueryAuthorizedMembersRequest,
) (*types.QueryAuthorizedMembersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	members := q.GetAuthorizedMembers(ctx)
	return &types.QueryAuthorizedMembersResponse{Members: members}, nil
}

// CurrentEpoch queries current epoch and its block height.
func (q Querier) CurrentEpoch(
	goCtx context.Context,
	_ *types.QueryCurrentEpochRequest,
) (*types.QueryCurrentEpochResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	epoch := q.GetCurrentEpoch(ctx)
	height := ctx.BlockHeight()

	return &types.QueryCurrentEpochResponse{
		Epoch:  epoch,
		Height: uint64(height),
	}, nil
}
