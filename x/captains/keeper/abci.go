package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// BeginBlocker called every block, process the epoch if it's ended.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	epoch := k.GetCurrentEpoch(ctx)

	if k.HasEndEpoch(ctx, epoch) {
		// prune useless epoch data
		k.delEpochEmission(ctx, epoch-1)
		k.delGlobalComputingPowerOnEpoch(ctx, epoch-1)

		// current epoch's
		k.delReportDigest(ctx, epoch)
		k.delEndEpoch(ctx, epoch)
		k.delReportBatches(ctx, epoch)

		// Let's enter new epoch!
		k.incrEpoch(ctx)
		// and we are in stand-by phrase again.
		k.delStandByFlag(ctx)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeBeginBlock,
				sdk.NewAttribute(types.AttributeKeyEpochID, fmt.Sprintf("%d", epoch)),
				sdk.NewAttribute(types.EventTypeEpochPhase, "into_stand_by"),
			),
		})

		// TODO: add telemetry
	}
}

// EndBlocker called every block, process the report digest if exists.
func (k Keeper) EndBlocker(ctx sdk.Context) {
	epoch := k.GetCurrentEpoch(ctx)

	if k.HasReportDigest(ctx, epoch) {
		// NOTE: there's a very scenario where reporter commits digest report but
		// also creates new nodes in the same block. In this case, there's a mismatch
		// between digest node count and actual node count. So we will check it again
		// and ask reporter to resubmit it later.
		digest, _ := k.GetReportDigest(ctx, epoch)
		if digest.TotalNodeCount != k.GetNodesCount(ctx) {
			k.delReportDigest(ctx, epoch)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeEndBlock,
					sdk.NewAttribute(types.AttributeKeyEpochID, fmt.Sprintf("%d", epoch)),
					sdk.NewAttribute(types.EventTypeEpochPhase, "fail_into_busy"),
				),
			})
		}

		// TODO: once we enter in busy phrase, we won't go back until
		// report ends. Considering if we need way to go back manually.
		k.execReportDigestEndBlock(ctx, digest)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeEndBlock,
				sdk.NewAttribute(types.AttributeKeyEpochID, fmt.Sprintf("%d", epoch)),
				sdk.NewAttribute(types.EventTypeEpochPhase, "into_busy"),
			),
		})
	}
}
