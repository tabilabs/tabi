package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// InitGenesis stores the NFT genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := data.Validate(); err != nil {
		panic(fmt.Errorf("failed to initialize mint genesis state: %s", err.Error()))
	}

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(fmt.Errorf("failed to set mint genesis state: %s", err.Error()))
	}
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
	// set nodes extra info
	for _, extra := range data.NodesExtraInfo {
		k.SetNodeHistoricalEmissionOnLastClaim(ctx, extra.Id, extra.LastClaimHistoricalEmission)
		for _, epoch := range extra.Epochs {
			k.setNodeComputingPowerOnEpoch(ctx, epoch.EpochId, extra.Id, epoch.ComputingPower)
			k.setNodeHistoricalEmissionByEpoch(ctx, epoch.EpochId, extra.Id, epoch.HistoricalEmission)
		}
	}
	// set epoch state
	currEpochID := data.EpochState.CurrEpoch
	k.setEpoch(ctx, currEpochID)
	if data.EpochState.IsEnd {
		k.setEndOnEpoch(ctx, currEpochID)
	}
	if data.EpochState.Digest != nil {
		k.setReportDigest(ctx, currEpochID, data.EpochState.Digest)
	}
	k.setEpochBase(ctx, currEpochID, data.EpochState.Current)
	k.setEpochBase(ctx, currEpochID-1, data.EpochState.Previous)
	k.SetGlobalClaimedEmission(ctx, data.EpochState.EmissionClaimedSum)

	// set computing power
	for _, power := range data.ClaimableComputingPowers {
		owner, err := sdk.AccAddressFromBech32(power.Owner)
		if err != nil {
			panic(fmt.Errorf("failed to convert addr from bech32: %s", err.Error()))
		}
		k.setClaimableComputingPower(ctx, power.Amount, owner)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:                   k.GetParams(ctx),
		EpochState:               k.GetEpochsState(ctx),
		Divisions:                k.GetDivisions(ctx),
		Nodes:                    k.GetNodes(ctx),
		NodesExtraInfo:           k.GetNodesExtraInfo(ctx),
		ClaimableComputingPowers: k.GetComputingPowersClaimable(ctx),
	}
}
