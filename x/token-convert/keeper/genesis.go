package keeper

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/token-convert/types"
)

func (k Keeper) InitGenesis(ctx sdktypes.Context, state types.GenesisState) {
	panic("implement me")
}

func (k Keeper) ExportGenesis(ctx sdktypes.Context) types.GenesisState {
	panic("implement me")
}
