package keeper

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/token-convert/types"
)

// InitGenesis sets the initial state of the module
func (k Keeper) InitGenesis(ctx sdktypes.Context, state types.GenesisState) {
	if err := types.ValidateGenesis(state); err != nil {
		panic(err)
	}

	for _, s := range state.Strategies {
		k.createStrategy(ctx, s.Name, s.Period, s.ConversionRate)
	}

	for _, v := range state.Vouchers {
		k.setVoucher(ctx, v.Id, v.Owner, v.Strategy, v.Amount)
		k.setVoucherByOwner(ctx, v.Owner, v.Id)
	}

	k.setVoucherSeq(ctx, state.VoucherSequence)
}

// ExportGenesis returns a GenesisState
func (k Keeper) ExportGenesis(ctx sdktypes.Context) *types.GenesisState {
	strategies := k.GetStrategies(ctx)
	vouchers := k.GetVouchers(ctx)
	return types.NewGenesisState(k.GetVoucherSeq(ctx), strategies, vouchers)
}
