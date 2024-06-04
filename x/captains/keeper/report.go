package keeper

import (
	"github.com/tabilabs/tabi/x/captains/types"

	sdkcdc "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
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
//
// NOTE: Part of handle logic are delayed to the end blocker. This is because we may encounter
// multiple params update in the same block which affect the calculation.
func (k Keeper) HandleReportDigest(ctx sdk.Context, report *types.ReportDigest) error {
	epochId := report.EpochId
	k.setReportDigest(ctx, epochId, report)

	return nil
}

// execReportDigestEndBlock processes a report digest at the end block
func (k Keeper) execReportDigestEndBlock(ctx sdk.Context, digest *types.ReportDigest) {
	epochId := k.GetCurrentEpoch(ctx)
	sum := k.CalcEpochEmission(ctx, epochId, digest.GlobalOnOperationRatio)

	k.DelGlobalPledge(ctx, epochId)
	k.setEpochEmission(ctx, epochId, sum)
	// we will enter report calculation in the next block.
	k.setStandByOverFlag(ctx)
}

// HandleReportBatch processes a report batch
func (k Keeper) HandleReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
	epochId := report.EpochId

	for _, node := range report.Nodes {
		owner, found := k.GetNodeOwner(ctx, node.NodeId)
		if !found {
			return errorsmod.Wrapf(types.ErrNodeNotExists, "node-%s not exists", node.NodeId)
		}

		// try to calculate historical emission
		k.CalcAndSetNodeCumulativeEmissionByEpoch(ctx, epochId-1, node.NodeId)
		power := k.CalcNodeComputingPowerOnEpoch(ctx, epochId, node.NodeId, node.OnOperationRatio)

		k.setNodeComputingPowerOnEpoch(ctx, epochId, node.NodeId, power)
		k.delOwnerPledge(ctx, owner, epochId-2) // it's fine to delete epoch(-1) which doesn't exist at all.
		k.incrGlobalComputingPowerOnEpoch(ctx, epochId, power)

		// sample owner pledge once for next epoch
		if !k.HasOwnerPledge(ctx, owner, epochId+1) {
			pledge, err := k.SampleOwnerPledge(ctx, owner)
			if err != nil {
				return err
			}

			k.SetOwnerPledge(ctx, owner, epochId+1, pledge)
			k.IncrGlobalPledge(ctx, epochId+1, pledge)
		}
	}

	// mark we have handle this batch.
	k.setReportBatch(ctx, epochId, report.BatchId, report.NodeCount)

	return nil
}

// HandleReportEnd processes a report end
func (k Keeper) HandleReportEnd(ctx sdk.Context, report *types.ReportEnd) error {
	epochId := report.EpochId

	// validate calculation finished.
	if err := k.IsReportCompleted(ctx, epochId); err != nil {
		return err
	}

	// marks we are ready for the next epoch.
	k.setEndOnEpoch(ctx, epochId)

	return nil
}

// IsReportCompleted checks if the report is completed
func (k Keeper) IsReportCompleted(ctx sdk.Context, epochId uint64) error {
	nodeCount := uint64(0)
	batchCount := uint64(0)

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ReportBatchOnEpochPrefixStoreKey(epochId))
	iterator := prefixStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		batchCount++
		nodeCount += sdk.BigEndianToUint64(iterator.Value())
	}

	digest, _ := k.GetReportDigest(ctx, epochId)
	if digest.TotalBatchCount != batchCount || digest.TotalNodeCount != nodeCount {
		return errorsmod.Wrapf(types.ErrInvalidReport, "commit end report too early")
	}

	return nil
}

// ValidateReport checks if the report is valid
func (k Keeper) ValidateReport(ctx sdk.Context, reportType types.ReportType, report *sdkcdc.Any) (any, error) {
	var message types.ReportContent
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
		if err := k.ValidateReportDigest(ctx, digest); err != nil {
			return nil, err
		}

		return digest, nil
	case types.ReportType_REPORT_TYPE_BATCH:
		batch, ok := message.(*types.ReportBatch)
		if !ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report")
		}
		if err := k.ValidateReportBatch(ctx, batch); err != nil {
			return nil, err
		}

		return batch, nil
	case types.ReportType_REPORT_TYPE_END:
		end, ok := message.(*types.ReportEnd)
		if !ok {
			return nil, errorsmod.Wrapf(types.ErrInvalidReport, "invalid report")
		}
		if err := k.ValidateReportEnd(ctx, end); err != nil {
			return nil, err
		}
		return end, nil
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

// ValidateReportDigest checks if the report digest is valid
func (k Keeper) ValidateReportDigest(ctx sdk.Context, report *types.ReportDigest) error {
	if err := k.ValidateReportEpoch(ctx, report.EpochId); err != nil {
		return err
	}

	_, found := k.GetReportDigest(ctx, report.EpochId)
	if found {
		return errorsmod.Wrapf(types.ErrInvalidReport, "digest already exists")
	}

	// NOTE:  assure all nodes created on chain submitted, otherwise emission calc will be incorrect.
	if report.TotalNodeCount != k.GetNodesCount(ctx) {
		return errorsmod.Wrapf(types.ErrInvalidReport, "node count mismatch")
	}

	return nil
}

// ValidateReportBatch checks if the report batch is valid
func (k Keeper) ValidateReportBatch(ctx sdk.Context, report *types.ReportBatch) error {
	if err := k.ValidateReportEpoch(ctx, report.EpochId); err != nil {
		return err
	}

	digest, found := k.GetReportDigest(ctx, report.EpochId)
	if !found {
		return errorsmod.Wrapf(types.ErrInvalidReport, "digest not found")
	}

	if report.NodeCount > digest.MaximumNodeCountPerBatch {
		return errorsmod.Wrapf(types.ErrInvalidReport, "node count exceeded")
	}

	if k.HasReportBatch(ctx, report.EpochId, report.BatchId) {
		return errorsmod.Wrapf(types.ErrInvalidReport, "batch already precessed")
	}

	for _, node := range report.Nodes {
		if !k.HasNode(ctx, node.NodeId) {
			return errorsmod.Wrapf(types.ErrNodeNotExists, "node-%s not exists", node.NodeId)
		}
	}

	return nil
}

// ValidateReportEnd checks if the report end is valid
func (k Keeper) ValidateReportEnd(ctx sdk.Context, report *types.ReportEnd) error {
	if err := k.ValidateReportEpoch(ctx, report.EpochId); err != nil {
		return err
	}

	if k.HasEndEpoch(ctx, report.EpochId) {
		return errorsmod.Wrapf(types.ErrInvalidReport, "eport end epoch already ended")
	}
	return nil
}

// HasReportDigest checks if the digest exists.
func (k Keeper) HasReportDigest(ctx sdk.Context, epochID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ReportDigestOnEpochStoreKey(epochID))
}

// GetReportDigest returns the digest.
func (k Keeper) GetReportDigest(ctx sdk.Context, epochID uint64) (*types.ReportDigest, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ReportDigestOnEpochStoreKey(epochID))
	if len(bz) == 0 {
		return nil, false
	}

	var digest types.ReportDigest
	k.cdc.Unmarshal(bz, &digest)
	return &digest, true
}

// setReportDigest sets the digest.
func (k Keeper) setReportDigest(ctx sdk.Context, epochID uint64, digest *types.ReportDigest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(digest)
	key := types.ReportDigestOnEpochStoreKey(epochID)
	store.Set(key, bz)
}

// delReportDigest deletes the digest.
func (k Keeper) delReportDigest(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ReportDigestOnEpochStoreKey(epochID))
}

// HasReportBatch checks if the batch exists.
func (k Keeper) HasReportBatch(ctx sdk.Context, epochID, batchID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ReportBatchOnEpochStoreKey(epochID, batchID))
}

// setReportBatch sets the batch count.
func (k Keeper) setReportBatch(ctx sdk.Context, epochID, batchID, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(count)
	key := types.ReportBatchOnEpochStoreKey(epochID, batchID)
	store.Set(key, bz)
}

// GetReportBatches returns the batch count.
func (k Keeper) GetReportBatches(ctx sdk.Context, epochID uint64) []types.BatchBase {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReportBatchOnEpochPrefixStoreKey(epochID))
	defer iterator.Close()

	var batches []types.BatchBase
	for ; iterator.Valid(); iterator.Next() {
		batchID := sdk.BigEndianToUint64(iterator.Key())
		count := sdk.BigEndianToUint64(iterator.Value())
		batches = append(batches, types.BatchBase{
			BatchId: batchID,
			Count:   count,
		})
	}

	return batches
}

// delReportBatches deletes the batch count.
func (k Keeper) delReportBatches(ctx sdk.Context, epochID uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReportBatchOnEpochPrefixStoreKey(epochID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
