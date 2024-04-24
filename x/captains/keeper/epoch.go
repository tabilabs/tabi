package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// epochs handle stage transitions in an epoch.

// PruneEpochs prunes useless state on epochs.
func (k Keeper) PruneEpochs(ctx sdk.Context) {
	panic("not implemented")
}

// GetCurrentEpoch returns the current epoch.
func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochKey)
	return sdk.BigEndianToUint64(bz)
}

// setEpoch sets the epoch id.
func (k Keeper) setEpoch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(k.GetCurrentEpoch(ctx) + 1)
	store.Set(types.EpochKey, bz)
}
