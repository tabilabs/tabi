package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// epochs handle stage transitions in an epoch.

// PruneEpochs prunes useless state on epochs.
func (k Keeper) PruneEpochs(ctx sdk.Context) {
	panic("not implemented")
}

func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	panic("implement me")
}
