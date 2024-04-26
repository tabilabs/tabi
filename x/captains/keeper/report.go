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
	epoch := report.EpochId

	// TODO:
	if epoch == 1 {
		// the first report commit results epoch getting into 2 from 1.
		// from now on its epoch 2
		k.setEpoch(ctx) // 2
	}

	epochEmissionSum, err := k.calcEpochEmission(ctx, epoch, report.GlobalOnOperationRatio)
	if err != nil {
		return err
	}

	// TODO: package this in a function
	historicalEmission := sdk.ZeroDec()
	if epoch > 1 {
		historicalEmission, err = k.GetHistoricalEmissionSum(ctx, epoch-1)
		if err != nil {
			return err
		}
	}
	historicalEmission.Add(epochEmissionSum)
	k.setHistoricalEmissionSum(ctx, epoch, historicalEmission)

	// save package digest
	k.setDigest(ctx, epoch, &report)

	return nil
}

// handleReportBatch processes a report batch
func (k Keeper) handleReportBatch(ctx sdk.Context, report types.ReportBatch) error {
	epoch := report.Epoch

	for _, nodeId := range report.NodeIds {
		owner := k.GetNodeOwner(ctx, nodeId)

		pledgeRatio, err := k.calcNodePledgeRatioOnEpoch(ctx, epoch, nodeId)
		if err != nil {
			return err
		}

		power, err := k.calcNodeComputingPowerOnEpoch(ctx, epoch, nodeId, pledgeRatio)
		if err != nil {
			return err
		}

		// accumulate computing power sum
		k.incrComputingPowerSumOnEpoch(ctx, epoch, power)

		// sample pledge by owner
		pledge, _ := k.GetOwnerPledge(ctx, owner, epoch)
		if pledge.Equal(sdk.ZeroDec()) {
			pledge, err = k.SampleOwnerPledge(ctx, owner)
			if err != nil {
				return err
			}
			k.setOwnerPledge(ctx, owner, epoch+1, pledge)

			sumPledge, _ := k.GetPledgeSum(ctx, epoch)
			sumPledge.Add(pledge)
			k.setPledgeSum(ctx, epoch, sumPledge)
		}
	}

	k.setBatchCount(ctx, epoch, report.BatchId, report.NodeCount)

	return nil
}

// handleReportEnd processes a report end
func (k Keeper) handleReportEnd(ctx sdk.Context, report types.ReportEnd) error {
	epoch := report.Epoch

	// validate calculation finished.
	// TODO: package this into func.
	nodeCount := uint64(0)
	batchCount := uint64(0)
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.BatchCountOnEpochPrefixKey(epoch))
	iterator := prefixStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		batchCount++
		nodeCount += sdk.BigEndianToUint64(iterator.Value())
	}

	digest := k.GetDigest(ctx, epoch)
	if digest.TotalBatchCount != batchCount || digest.TotalNodeCount != nodeCount {
		return errorsmod.Wrapf(types.ErrEpochUnfinished, "commit later")
	}

	// marks we can get into the next epoch.
	k.setEndEpoch(ctx, epoch+1)

	k.delHistoricalEmissionSum(ctx, epoch-1)
	k.delEpochEmission(ctx, epoch)
	k.delComputingPowerSumOnEpoch(ctx, epoch)

	return nil
}
