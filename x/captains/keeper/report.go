package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// CommitReport processes a report
func (k Keeper) CommitReport(ctx sdk.Context, reportType types.ReportType, report []byte) error {
	switch reportType {
	case types.ReportType_REPORT_TYPE_DIGEST:
		var digest types.ReportDigest
		if err := k.cdc.Unmarshal(report, &digest); err != nil {
			return err
		}
		return k.handleReportDigest(ctx, digest)
	case types.ReportType_REPORT_TYPE_BATCH:
		var batch types.ReportBatch
		if err := k.cdc.Unmarshal(report, &batch); err != nil {
			return err
		}
		return k.handleReportBatch(ctx, batch)
	case types.ReportType_REPORT_TYPE_END:
		var end types.ReportEnd
		if err := k.cdc.Unmarshal(report, &end); err != nil {
			return err
		}
		return k.handleReportEnd(ctx, end)
	}
	return errorsmod.Wrapf(types.ErrInvalidReportType, "report type: %s", reportType)
}

// handleReportDigest processes a report digest
func (k Keeper) handleReportDigest(ctx sdk.Context, report types.ReportDigest) error {
	epochId := report.EpochId

	epochEmissionSum, err := k.calcEpochEmission(ctx, epochId, report.GlobalOnOperationRatio)
	if err != nil {
		return err
	}

	k.incrHistoricalEmissionSum(ctx, epochId, epochEmissionSum)
	k.setDigest(ctx, epochId, &report)

	return nil
}

// handleReportBatch processes a report batch
func (k Keeper) handleReportBatch(ctx sdk.Context, report types.ReportBatch) error {
	epochId := report.EpochId

	for _, nodeId := range report.NodeIds {
		owner := k.GetNodeOwner(ctx, nodeId)

		pledgeRatio, err := k.calcNodePledgeRatioOnEpoch(ctx, epochId, nodeId)
		if err != nil {
			return err
		}

		power, err := k.calcNodeComputingPowerOnEpoch(ctx, epochId, nodeId, pledgeRatio)
		if err != nil {
			return err
		}

		// accumulate computing power sum
		k.incrComputingPowerSumOnEpoch(ctx, epochId, power)

		// sample owner pledge once for next epoch
		pledge, found := k.GetOwnerPledge(ctx, owner, epochId+1)
		if !found {
			pledge, err = k.SampleOwnerPledge(ctx, owner)
			if err != nil {
				return err
			}

			k.setOwnerPledge(ctx, owner, epochId+1, pledge)
			k.incrPledgeSum(ctx, epochId+1, pledge)
		}
	}

	// mark we have handle this batch.
	k.setReportBatch(ctx, epochId, report.BatchId, report.NodeCount)

	return nil
}

// handleReportEnd processes a report end
func (k Keeper) handleReportEnd(ctx sdk.Context, report types.ReportEnd) error {
	epochId := report.Epoch

	// validate calculation finished.
	if err := k.validateReportCompleted(ctx, epochId); err != nil {
		return err
	}

	// prune useless epoch data
	k.delHistoricalEmissionSum(ctx, epochId-1)
	k.delEpochEmission(ctx, epochId-1)
	k.delComputingPowerSumOnEpoch(ctx, epochId-1)

	// marks we are ready for the next epoch.
	k.setEndEpoch(ctx, epochId)

	return nil
}

// validateReportCompleted checks if the report is completed
func (k Keeper) validateReportCompleted(ctx sdk.Context, epochId uint64) error {
	nodeCount := uint64(0)
	batchCount := uint64(0)

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ReportBatchOnEpochPrefixKey(epochId))
	iterator := prefixStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		batchCount++
		nodeCount += sdk.BigEndianToUint64(iterator.Value())
	}

	digest := k.GetDigest(ctx, epochId)
	if digest.TotalBatchCount != batchCount || digest.TotalNodeCount != nodeCount {
		return errorsmod.Wrapf(types.ErrEpochUnfinished, "commit later")
	}

	return nil
}
