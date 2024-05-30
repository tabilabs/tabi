package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"
)

// ValidateBasic validate a report digest
func (digest *ReportDigest) ValidateBasic() error {
	if digest.EpochId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch id is zero")
	}
	if digest.TotalBatchCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total batch is zero")
	}
	if digest.TotalNodeCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total node is zero")
	}
	if digest.MaximumNodeCountPerBatch == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max node per batch id is zero")
	}
	if !digest.GlobalOnOperationRatio.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "operation ratio must be greater than zero")
	}
	return nil
}

// ValidateBasic validate a report batch
func (batch *ReportBatch) ValidateBasic() error {
	if batch.EpochId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch id is zero")
	}
	if batch.BatchId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "batch id is zero")
	}
	if batch.NodeCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node count is zero")
	}
	if uint64(len(batch.Nodes)) != batch.NodeCount {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node count node ids length unmatched")
	}

	seenMap := make(map[string]bool)
	for _, node := range batch.Nodes {
		if node.NodeId == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node id is empty")
		}
		if seenMap[node.NodeId] {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate node id")
		}
		if node.OnOperationRatio.IsNegative() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "operation ratio must be positive")
		}
		seenMap[node.NodeId] = true
	}

	return nil
}
