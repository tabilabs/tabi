package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

// Strategy queries the strategy of a given name
func (q Querier) Strategy(goCtx context.Context, req *types.QueryStrategyRequest) (*types.QueryStrategyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "strategy name cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	strategy, found := q.GetStrategy(ctx, req.Name)
	if !found {
		return nil, status.Error(codes.NotFound, "strategy not found")
	}

	return &types.QueryStrategyResponse{
		Name:           strategy.Name,
		Period:         strategy.Period,
		ConversionRate: strategy.ConversionRate.String(),
	}, nil
}

// Strategies queries all strategies
func (q Querier) Strategies(goCtx context.Context, req *types.QueryStrategiesRequest) (*types.QueryStrategiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var strategies []types.Strategy
	strategyStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.StrategyKey)
	pageRes, err := query.Paginate(strategyStore, req.Pagination, func(key []byte, value []byte) error {
		var strategy types.Strategy
		q.cdc.MustUnmarshal(value, &strategy)
		strategies = append(strategies, strategy)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryStrategiesResponse{
		Strategies: strategies,
		Pagination: pageRes,
	}, nil
}

// Voucher queries a voucher by its id
func (q Querier) Voucher(goCtx context.Context, req *types.QueryVoucherRequest) (*types.QueryVoucherResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(req.VoucherId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "voucher id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	voucher, found := q.GetVoucher(ctx, req.VoucherId)
	if !found {
		return nil, status.Error(codes.NotFound, "voucher not found")
	}

	return &types.QueryVoucherResponse{
		Id:          voucher.Id,
		Owner:       voucher.Owner,
		Strategy:    voucher.Strategy,
		CreatedTime: voucher.CreatedTime,
	}, nil
}

// Vouchers queries all vouchers
func (q Querier) Vouchers(goCtx context.Context, req *types.QueryVouchersRequest) (*types.QueryVouchersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(req.Owner) != 0 {
		owner, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner address %s", req.Owner)
		}
		return q.vouchersByOwner(goCtx, owner, req.Pagination)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var vouchers []types.Voucher
	voucherStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.VoucherKey)
	pageRes, err := query.Paginate(voucherStore, req.Pagination, func(key []byte, value []byte) error {
		var voucher types.Voucher
		q.cdc.MustUnmarshal(value, &voucher)
		vouchers = append(vouchers, voucher)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryVouchersResponse{
		Vouchers:   vouchers,
		Pagination: pageRes,
	}, nil
}

// vouchersByOwner queries vouchers by owner
func (q Querier) vouchersByOwner(goCtx context.Context, owner sdk.AccAddress, pageReq *query.PageRequest) (*types.QueryVouchersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var vouchers []types.Voucher
	voucherByOwnerStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.VoucherByOwnerStorePrefixKey(owner))
	pageRes, err := query.Paginate(voucherByOwnerStore, pageReq, func(key []byte, _ []byte) error {
		voucher, _ := q.GetVoucher(ctx, string(key))
		vouchers = append(vouchers, voucher)
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryVouchersResponse{
		Vouchers:   vouchers,
		Pagination: pageRes,
	}, nil
}

// VoucherStatus queries the lock status of a given voucher
func (q Querier) VoucherStatus(goCtx context.Context, req *types.QueryVoucherStatusRequest) (*types.QueryVoucherStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if len(req.VoucherId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "voucher id cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	voucher, found := q.GetVoucher(ctx, req.VoucherId)
	if !found {
		return nil, status.Error(codes.NotFound, "voucher not found")
	}

	strategy, found := q.GetStrategy(ctx, voucher.Strategy)
	if !found {
		return nil, status.Error(codes.NotFound, "strategy not found")
	}

	withdrawableTabi, _, returnableVetabi := q.calVoucher(ctx, voucher, strategy)
	return &types.QueryVoucherStatusResponse{
		CurrentTime:      ctx.BlockTime().String(),
		TabiWithdrawable: withdrawableTabi,
		VetabiReturnable: returnableVetabi,
	}, nil
}
