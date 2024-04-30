package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// CommitReport processes a report
func (k Keeper) CommitReport(ctx sdk.Context, report any) error {
	switch report := report.(type) {
	case *types.ReportDigest:
		return k.handleReportDigest(ctx, report)
	case *types.ReportBatch:
		return k.handleReportBatch(ctx, report)
	case *types.ReportEnd:
		return k.handleReportEnd(ctx, report)
	}
	return errorsmod.Wrapf(types.ErrInvalidReport, "invalid report type")
}

// handleReportDigest processes a report digest
func (k Keeper) handleReportDigest(ctx sdk.Context, report *types.ReportDigest) error {
	epochId := report.EpochId

	_, err := k.calcEpochEmission(ctx, epochId, report.GlobalOnOperationRatio)
	if err != nil {
		return err
	}

	k.setDigest(ctx, epochId, report)

	return nil
}

// handleReportBatch processes a report batch
func (k Keeper) handleReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
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
func (k Keeper) handleReportEnd(ctx sdk.Context, report *types.ReportEnd) error {
	epochId := report.Epoch

	// validate calculation finished.
	if err := k.isReportCompleted(ctx, epochId); err != nil {
		return err
	}

	// prune useless epoch data
	k.delEpochEmission(ctx, epochId-1)
	k.delComputingPowerSumOnEpoch(ctx, epochId-1)

	// marks we are ready for the next epoch.
	k.setEndEpoch(ctx, epochId)

	return nil
}

// isReportCompleted checks if the report is completed
func (k Keeper) isReportCompleted(ctx sdk.Context, epochId uint64) error {
	nodeCount := uint64(0)
	batchCount := uint64(0)

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ReportBatchOnEpochPrefixKey(epochId))
	iterator := prefixStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		batchCount++
		nodeCount += sdk.BigEndianToUint64(iterator.Value())
	}

	digest, _ := k.GetDigest(ctx, epochId)
	if digest.TotalBatchCount != batchCount || digest.TotalNodeCount != nodeCount {
		return errorsmod.Wrapf(types.ErrInvalidReport, "commit end report too early!")
	}

	return nil
}

// validateReport checks if the report is valid
func (k Keeper) validateReport(ctx sdk.Context, reportType types.ReportType, report []byte) (any, error) {
	switch reportType {
	case types.ReportType_REPORT_TYPE_DIGEST:
		var digest types.ReportDigest
		if err := k.cdc.Unmarshal(report, &digest); err != nil {
			return nil, err
		}
		if err := k.validateReportEpoch(ctx, digest.EpochId); err != nil {
			return nil, err
		}
		return &digest, nil

	case types.ReportType_REPORT_TYPE_BATCH:
		var batch types.ReportBatch
		if err := k.cdc.Unmarshal(report, &batch); err != nil {
			return nil, err
		}
		if err := k.validateReportEpoch(ctx, batch.EpochId); err != nil {
			return nil, err
		}
		if err := k.validateReportBatch(ctx, &batch); err != nil {
			return nil, err
		}

		return &batch, nil
	case types.ReportType_REPORT_TYPE_END:
		var end types.ReportEnd
		if err := k.cdc.Unmarshal(report, &end); err != nil {
			return nil, err
		}
		if err := k.validateReportEpoch(ctx, end.Epoch); err != nil {
			return nil, err
		}
		return &end, nil
	}
	return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report type")
}

// validateReportDigest checks if the report digest is valid
func (k Keeper) validateReportEpoch(ctx sdk.Context, epochID uint64) error {
	if epochID != k.GetCurrentEpoch(ctx) {
		return errorsmod.Wrapf(types.ErrInvalidReport, "invalid epoch")
	}
	return nil
}

// validateReportBatch checks if the report batch is valid
func (k Keeper) validateReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
	digest, found := k.GetDigest(ctx, report.EpochId)
	if !found {
		return errorsmod.Wrapf(types.ErrInvalidReport, "digest not found")
	}

	if report.NodeCount > digest.MaximumNodeCountPerBatch {
		return errorsmod.Wrapf(types.ErrInvalidReport, "node count exceeded")
	}

	for _, nodeId := range report.NodeIds {
		if !k.HasNode(ctx, nodeId) {
			return errorsmod.Wrapf(types.ErrNodeNotExists, "node-%s not exists", nodeId)
		}
	}

	return nil
}
