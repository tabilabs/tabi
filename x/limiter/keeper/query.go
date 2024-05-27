package keeper

import (
	"context"

	"github.com/tabilabs/tabi/x/limiter/types"
)

type Querier struct {
	*Keeper
}

func NewQuerierImpl(keeper *Keeper) Querier {
	return Querier{keeper}
}

// Params queries the params of limiter.
// NOTE: Use x/params instead until before sdk v0.47.
func (q Querier) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}
