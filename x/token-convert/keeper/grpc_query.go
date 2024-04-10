package keeper

import (
	"context"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Strategy(goCtx context.Context, req *types.QueryStrategyRequest) (*types.QueryStrategyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Strategies(goCtx context.Context, req *types.QueryStrategiesRequest) (*types.QueryStrategiesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Voucher(goCtx context.Context, req *types.QueryVoucherRequest) (*types.QueryVoucherResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Vouchers(goCtx context.Context, req *types.QueryVouchersRequest) (*types.QueryVouchersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) UnlockStatus(goCtx context.Context, req *types.QueryUnlockStatusRequest) (*types.QueryUnlockStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}
