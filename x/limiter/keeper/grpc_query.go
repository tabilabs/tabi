package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/limiter/types"
)

type Querier struct {
	k *Keeper
}

func NewQuerierImpl(keeper *Keeper) Querier {
	return Querier{keeper}
}

// Params queries the params of limiter.
func (q Querier) Params(goCtx context.Context, msg *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
