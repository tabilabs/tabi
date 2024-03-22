package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captain-node/types"
)

func (k Keeper) UpdateSaleLevel(ctx sdk.Context, saleLevel uint64) (sdk.Event, error) {
	params := k.GetParams(ctx)
	oldLevel := params.CurrentLevelForSale
	if oldLevel > saleLevel {
		return sdk.Event{}, types.ErrInvalidSaleLevel
	}

	params.CurrentLevelForSale = saleLevel
	event := sdk.NewEvent(
		types.EventTypeAddCaller,
		sdk.NewAttribute(types.AttributeKeyOldSaleLevel, fmt.Sprintf("%d", oldLevel)),
		sdk.NewAttribute(types.AttributeKeyNewSaleLevel, fmt.Sprintf("%d", saleLevel)),
	)
	return event, k.SetParams(ctx, params)
}

// GetParams sets the mint module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the mint module parameters.
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

func (k Keeper) GetSaleLevel(ctx sdk.Context) uint64 {
	return k.GetParams(ctx).CurrentLevelForSale
}
