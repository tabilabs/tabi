package keeper

import (
	"context"

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
	//todo
	return &types.QueryOwnerResponse{}, nil
}

// Supply queries the number of Node from the given division
func (k Keeper) Supply(goCtx context.Context, request *types.QuerySupplyRequest) (*types.QuerySupplyResponse, error) {
	//todo
	return &types.QuerySupplyResponse{}, nil
}

// Division queries an node division by its ID
func (k Keeper) Division(goCtx context.Context, request *types.QueryDivisionRequest) (*types.QueryDivisionResponse, error) {
	//todo
	return &types.QueryDivisionResponse{}, nil
}

// Divisions queries all Node divisions
func (k Keeper) Divisions(goCtx context.Context, request *types.QueryDivisionsRequest) (*types.QueryDivisionsResponse, error) {
	//todo
	return &types.QueryDivisionsResponse{}, nil
}

// Node queries an Node based on its id.
func (k Keeper) Node(goCtx context.Context, request *types.QueryNodeRequest) (*types.QueryNodeResponse, error) {
	//todo
	return &types.QueryNodeResponse{}, nil
}

// Nodes queries all node of a given owner
func (k Keeper) Nodes(goCtx context.Context, request *types.QueryNodesRequest) (*types.QueryNodesResponse, error) {
	//todo
	return &types.QueryNodesResponse{}, nil
}
