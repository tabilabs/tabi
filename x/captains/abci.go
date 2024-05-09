package captains

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/keeper"
)

// BeginBlocker runs at the start of each block
func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	epoch := k.GetCurrentEpoch(ctx)

	if k.HasEndEpoch(ctx, epoch) {
		k.DelDigest(ctx, epoch)
		k.DelEndEpoch(ctx, epoch)
		k.DelReportBatches(ctx, epoch)
		// Let's enter new epoch!
		k.EnterNewEpoch(ctx)

		// TODO: add log
	}
}
