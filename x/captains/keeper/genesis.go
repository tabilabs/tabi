package keeper

import (
	"fmt"

	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores the NFT genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	if err := data.Validate(); err != nil {
		panic(fmt.Errorf("failed to initialize mint genesis state: %s", err.Error()))
	}

	// set params
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(fmt.Errorf("failed to set mint genesis state: %s", err.Error()))
	}

	// set base state
	k.SetBaseState(ctx, &data.BaseState)

	// set divisions
	for _, division := range data.Divisions {
		if err := k.SaveDivision(ctx, division); err != nil {
			panic(fmt.Errorf("failed to save division: %s", err.Error()))
		}
	}

	// set nodes
	for _, node := range data.Nodes {
		if err := k.setNode(ctx, node); err != nil {
			panic(fmt.Errorf("failed to set node: %s", err.Error()))
		}
	}

	// set emission
	for _, ee := range data.EpochesEmission {
		k.setEpochEmission(ctx, ee.EpochId, ee.Emission)
	}
	for _, nce := range data.NodesClaimedEmission {
		k.SetNodeClaimedEmission(ctx, nce.NodeId, nce.Emission)
	}
	for _, nce := range data.NodesCumulativeEmission {
		k.setNodeCumulativeEmissionByEpoch(ctx, nce.EpochId, nce.NodeId, nce.Emission)
	}

	// set pledge
	for _, gp := range data.GlobalsPledge {
		k.SetGlobalPledge(ctx, gp.EpochId, gp.Amount)
	}
	for _, op := range data.OwnersPledge {
		k.SetOwnerPledge(ctx, sdk.MustAccAddressFromBech32(op.Owner), op.EpochId, op.Amount)
	}

	// set powers
	for _, cp := range data.OwnersClaimableComputingPower {
		k.setClaimableComputingPower(ctx, cp.Amount, sdk.MustAccAddressFromBech32(cp.Owner))
	}
	for _, gcp := range data.GlobalsComputingPower {
		k.setGlobalComputingPowerOnEpoch(ctx, gcp.EpochId, gcp.Amount)
	}
	for _, ncp := range data.NodesComputingPower {
		k.setNodeComputingPowerOnEpoch(ctx, ncp.EpochId, ncp.NodeId, ncp.Amount)
	}

	// set batches
	for _, batch := range data.Batches {
		k.setReportBatch(ctx, data.BaseState.EpochId, batch.BatchId, batch.Count)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:                        k.GetParams(ctx),
		BaseState:                     k.GetBaseState(ctx),
		Divisions:                     k.GetDivisions(ctx),
		Nodes:                         k.GetNodes(ctx),
		EpochesEmission:               k.GetEpochesEmission(ctx),
		NodesClaimedEmission:          k.GetNodesClaimedEmission(ctx),
		NodesCumulativeEmission:       k.GetNodesCumulativeEmission(ctx),
		GlobalsPledge:                 k.GetGlobalsPledge(ctx),
		OwnersPledge:                  k.GetOwnersPledge(ctx),
		OwnersClaimableComputingPower: k.GetClaimableComputingPowers(ctx),
		GlobalsComputingPower:         k.GetGlobalsComputingPower(ctx),
		NodesComputingPower:           k.GetNodesComputingPower(ctx),
		Batches:                       k.GetReportBatches(ctx, k.GetCurrentEpoch(ctx)),
	}
}

// GetBaseState returns the base state of the current epoch.
func (k Keeper) GetBaseState(ctx sdk.Context) types.BaseState {
	var baseState types.BaseState
	epochId := k.GetCurrentEpoch(ctx)

	baseState.EpochId = epochId
	baseState.NextNodeSequence = k.GetNodeSequence(ctx)
	baseState.GlobalClaimedEmission = k.GetGlobalClaimedEmission(ctx)

	if k.HasEndEpoch(ctx, epochId) {
		baseState.IsEpochEnd = true
	}

	if digest, found := k.GetReportDigest(ctx, epochId); found {
		baseState.ReportDigest = digest
	}

	if k.IsStandByPhase(ctx) {
		baseState.IsStandBy = true
	}

	return baseState
}

// SetBaseState sets the base state of the current epoch.
func (k Keeper) SetBaseState(ctx sdk.Context, bs *types.BaseState) {
	k.setEpoch(ctx, bs.EpochId)
	k.SetNodeSequence(ctx, bs.NextNodeSequence)
	if bs.IsEpochEnd {
		k.setEndOnEpoch(ctx, bs.EpochId)
	}
	if !bs.GlobalClaimedEmission.IsZero() {
		k.SetGlobalClaimedEmission(ctx, bs.GlobalClaimedEmission)
	}
	if bs.ReportDigest != nil {
		k.setReportDigest(ctx, bs.EpochId, bs.ReportDigest)
	}
	if !bs.IsStandBy {
		k.setStandByOverFlag(ctx)
	}
}
