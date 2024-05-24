package captains

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/keeper"
)

// BeginBlocker runs at the start of each block
func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	k.BeginBlocker(ctx)
}

// EndBlocker runs at the end of each block
func EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {
	k.EndBlocker(ctx)
	return []abci.ValidatorUpdate{}
}
