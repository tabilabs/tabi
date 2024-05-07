package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkcdc "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/tabilabs/tabi/x/captains/types"
)

// CommitReport processes a report
func (k Keeper) CommitReport(ctx sdk.Context, report any) error {
	switch report := report.(type) {
	case *types.ReportDigest:
		return k.HandleReportDigest(ctx, report)
	case *types.ReportBatch:
		return k.HandleReportBatch(ctx, report)
	case *types.ReportEnd:
		return k.HandleReportEnd(ctx, report)
	}
	return errorsmod.Wrapf(types.ErrInvalidReport, "invalid report type")
}

// HandleReportDigest processes a report digest
func (k Keeper) HandleReportDigest(ctx sdk.Context, report *types.ReportDigest) error {
	epochId := report.EpochId

	_, err := k.calcEpochEmission(ctx, epochId, report.GlobalOnOperationRatio)
	if err != nil {
		return err
	}

	k.SetDigest(ctx, epochId, report)

	return nil
}

// HandleReportBatch processes a report batch
func (k Keeper) HandleReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
	epochId := report.EpochId

	for _, nodeId := range report.NodeIds {

		owner := k.GetNodeOwner(ctx, nodeId)

		pledgeRatio, err := k.CalcNodePledgeRatioOnEpoch(ctx, epochId, nodeId)
		if err != nil {
			return err
		}

		power, err := k.CalcNodeComputingPowerOnEpoch(ctx, epochId, nodeId, pledgeRatio)
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
	k.SetReportBatch(ctx, epochId, report.BatchId, report.NodeCount)

	return nil
}

// HandleReportEnd processes a report end
func (k Keeper) HandleReportEnd(ctx sdk.Context, report *types.ReportEnd) error {
	epochId := report.Epoch

	// validate calculation finished.
	if err := k.IsReportCompleted(ctx, epochId); err != nil {
		return err
	}

	// prune useless epoch data
	k.delEpochEmission(ctx, epochId-1)
	k.delComputingPowerSumOnEpoch(ctx, epochId-1)

	// marks we are ready for the next epoch.
	k.SetEndEpoch(ctx, epochId)

	return nil
}

// IsReportCompleted checks if the report is completed
func (k Keeper) IsReportCompleted(ctx sdk.Context, epochId uint64) error {
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

// ValidateReport checks if the report is valid
func (k Keeper) ValidateReport(ctx sdk.Context, reportType types.ReportType, report *sdkcdc.Any) (any, error) {
	var message proto.Message
	if report == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidReport, "report is nil")
	}
	if err := k.cdc.UnpackAny(report, &message); err != nil {
		return "", err
	}

	switch reportType {
	case types.ReportType_REPORT_TYPE_DIGEST:
		digest, ok := message.(*types.ReportDigest)
		if !ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report")
		}

		if err := k.ValidateReportEpoch(ctx, digest.EpochId); err != nil {
			return nil, err
		}

		// TODO: maybe validate digest fields as well?

		return &digest, nil
	case types.ReportType_REPORT_TYPE_BATCH:
		batch, ok := message.(*types.ReportBatch)
		if !ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report")
		}

		if err := k.ValidateReportEpoch(ctx, batch.EpochId); err != nil {
			return nil, err
		}
		if err := k.ValidateReportBatch(ctx, batch); err != nil {
			return nil, err
		}

		return &batch, nil
	case types.ReportType_REPORT_TYPE_END:
		end, ok := message.(*types.ReportEnd)
		if !ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report")
		}
		if err := k.ValidateReportEpoch(ctx, end.Epoch); err != nil {
			return nil, err
		}
		return &end, nil
	}
	return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report type")
}

// ValidateReportEpoch checks if the report is valid
func (k Keeper) ValidateReportEpoch(ctx sdk.Context, epochID uint64) error {
	if epochID != k.GetCurrentEpoch(ctx) {
		return errorsmod.Wrapf(types.ErrInvalidReport, "invalid epoch")
	}
	return nil
}

// ValidateReportBatch checks if the report batch is valid
func (k Keeper) ValidateReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
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
