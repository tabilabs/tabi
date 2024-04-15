package keeper

import (
	"context"

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

func (k Keeper) NodeTotalRewards(ctx context.Context, request *types.QueryNodeTotalRewardsRequest) (*types.QueryNodeTotalRewardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) HolderUnclaimedRewards(ctx context.Context, request *types.QueryHolderUnclaimedRewardsRequest) (*types.QueryHolderUnclaimedRewardsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) HolderTotalRewards(ctx context.Context, request *types.QueryHolderTotalRewardsRequest) (*types.QueryHolderTotalRewardsResponse, error) {
	//TODO implement me
	panic("implement me")
}
