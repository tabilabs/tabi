package keeper

import (
	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams sets the captains module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the captains module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

// UpdateSaleLevel updates the sale level in params.
func (k Keeper) UpdateSaleLevel(ctx sdk.Context, newSaleLevel uint64) (uint64, error) {
	params := k.GetParams(ctx)
	oldSaleLevel := params.CurrentSaleLevel
	if oldSaleLevel > newSaleLevel {
		return oldSaleLevel, types.ErrInvalidSaleLevel
	}

	params.CurrentSaleLevel = newSaleLevel
	if err := k.SetParams(ctx, params); err != nil {
		return oldSaleLevel, err
	}

	return oldSaleLevel, nil
}

// GetSaleLevel returns the current sale level.
func (k Keeper) GetSaleLevel(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).CurrentSaleLevel
}
