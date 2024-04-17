package keeper

import (
	"context"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

type Querier struct {
	*Keeper
}

// NewQuerierImpl returns an implementation of the token-convert QueryServer interface.
func NewQuerierImpl(k *Keeper) types.QueryServer {
	return &Querier{k}
}

var _ types.QueryServer = Querier{}

func (q Querier) Strategy(goCtx context.Context, req *types.QueryStrategyRequest) (*types.QueryStrategyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) Strategies(goCtx context.Context, req *types.QueryStrategiesRequest) (*types.QueryStrategiesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) Voucher(goCtx context.Context, req *types.QueryVoucherRequest) (*types.QueryVoucherResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) Vouchers(goCtx context.Context, req *types.QueryVouchersRequest) (*types.QueryVouchersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q Querier) UnlockStatus(goCtx context.Context, req *types.QueryUnlockStatusRequest) (*types.QueryUnlockStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}
