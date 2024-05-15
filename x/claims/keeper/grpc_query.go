package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/claims/types"
)

var _ types.QueryServer = Keeper{}

type Querier struct {
	*Keeper
}

// NewQuerierImpl returns an implementation of the captains QueryServer interface.
func NewQuerierImpl(k *Keeper) types.QueryServer {
	return &Querier{k}
}

var _ types.QueryServer = Querier{}

// Params queries the staking parameters
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) NodeTotalRewards(goCtx context.Context, request *types.QueryNodeTotalRewardsRequest) (*types.QueryNodeTotalRewardsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if request.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty holder address")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	rewards, err := k.CalculateRewardsByNodeId(ctx, request.NodeId)
	if err != nil {
		return nil, err
	}

	return &types.QueryNodeTotalRewardsResponse{Rewards: rewards}, nil
}

func (k Keeper) HolderTotalRewards(goCtx context.Context, request *types.QueryHolderTotalRewardsRequest) (*types.QueryHolderTotalRewardsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if request.Owner == "" {
		return nil, status.Error(codes.InvalidArgument, "empty holder address")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return nil, err
	}

	nodes := k.captainsKeeper.GetNodesByOwner(ctx, owner.Bytes())
	if len(nodes) == 0 {
		return nil, types.ErrHolderNotFound
	}

	rewards, err := k.CalculateRewards(ctx, nodes)
	if err != nil {
		return nil, err
	}

	return &types.QueryHolderTotalRewardsResponse{Rewards: rewards}, nil
}
