package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// Emission handles all emission calculations.

// EpochEmissionByNode returns the emission reward for a node in an epoch.
func (k Keeper) EpochEmissionByNode(ctx sdk.Context, epochID, nodeID string) {
	panic("not implemented")
}

// EpochEmissionSum returns the total emission reward for an epoch.
func (k Keeper) EpochEmissionSum(ctx sdk.Context, onOperationRatio sdk.Dec) sdk.Dec {
	panic("not implemented")
}

// EpochEmissionBase returns the base emission reward for an epoch.
func (k Keeper) EpochEmissionBase(ctx sdk.Context) sdk.Dec {
	panic("not implemented")
}

func (k Keeper) HistoricalEmissionSum(ctx sdk.Context, epochID string) sdk.Dec {
	panic("not implemented")
}

func (k Keeper) HistoricalEmissionByOwner(ctx sdk.Context, owner sdk.AccAddress) sdk.Dec {
	panic("not implemented")
}
