package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/limiter/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	panic("implement me")
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	panic("implement me")
}
