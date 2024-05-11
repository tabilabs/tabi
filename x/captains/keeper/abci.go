package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	epoch := k.GetCurrentEpoch(ctx)

	if k.HasEndEpoch(ctx, epoch) {
		// prune useless epoch data
		k.delEpochEmission(ctx, epoch-1)
		k.delComputingPowerSumOnEpoch(ctx, epoch-1)

		// current epoch's
		k.delDigest(ctx, epoch)
		k.delEndEpoch(ctx, epoch)
		k.delReportBatches(ctx, epoch)

		// Let's enter new epoch!
		k.incrEpoch(ctx)

		// TODO: add log
	}
}
