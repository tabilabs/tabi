package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/claims/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the staking parameters
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) NodeTotalRewards(goCtx context.Context, request *types.QueryNodeTotalRewardsRequest) (*types.QueryNodeTotalRewardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) HolderTotalRewards(goCtx context.Context, request *types.QueryHolderTotalRewardsRequest) (*types.QueryHolderTotalRewardsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if request.Owner == "" {
		return nil, status.Error(codes.InvalidArgument, "empty holder address")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := sdk.ValAddressFromBech32(request.Owner)
	if err != nil {
		return nil, err
	}

	rewards, err := k.CalculateRewardsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	return &types.QueryHolderTotalRewardsResponse{Rewards: sdk.NewDecCoinsFromCoins(rewards...)}, nil
}
