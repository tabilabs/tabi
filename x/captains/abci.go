package captains

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	epoch := k.GetCurrentEpoch(ctx)
	if epoch == 1 {
		return
	}

	if len(k.GetEndEpoch(ctx, epoch)) == 0 {
		return
	}

	k.DelBatchCount(ctx, epoch)

	// Let's go into new epoch!
	k.EnterNewEpoch(ctx)
}
