package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/tabilabs/tabi/x/captains/types"
)

// SaveDivision saves a division to the store
func (k Keeper) SaveDivision(ctx sdk.Context, division types.Division) error {
	if k.HasDivision(ctx, division.Id) {
		return errorsmod.Wrap(types.ErrDivisionExists, division.Id)
	}
	return k.setDivision(ctx, division)
}

// setDivision sets a division to the store
func (k Keeper) setDivision(ctx sdk.Context, division types.Division) error {
	bz, err := k.cdc.Marshal(&division)
	if err != nil {
		return errorsmod.Wrap(err, "Marshal division failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DivisionStoreKey(division.Id), bz)
	return nil
}

// GetDivision returns the division by the specified division ID
func (k Keeper) GetDivision(ctx sdk.Context, divisionID string) (types.Division, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DivisionStoreKey(divisionID))

	var division types.Division
	if len(bz) == 0 {
		return division, false
	}
	k.cdc.MustUnmarshal(bz, &division)
	return division, true
}

// GetDivisions returns all divisions
func (k Keeper) GetDivisions(ctx sdk.Context) (divisions []types.Division) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DivisionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var division types.Division
		k.cdc.MustUnmarshal(iterator.Value(), &division)
		divisions = append(divisions, division)
	}
	return
}

// HasDivision checks if the division exists in the store
func (k Keeper) HasDivision(ctx sdk.Context, divisionID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.DivisionStoreKey(divisionID))
}

// DecideDivision decides the division of the node as per power.
func (k Keeper) DecideDivision(ctx sdk.Context, power uint64) types.Division {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DivisionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var division types.Division
		k.cdc.MustUnmarshal(iterator.Value(), &division)
		if power <= division.ComputingPowerUpperBound && power >= division.ComputingPowerLowerBound {
			return division
		}
	}
	return types.Division{}
}

// incrDivisionTotalCount increments the sold count of the division
func (k Keeper) incrDivisionTotalCount(ctx sdk.Context, division types.Division) {
	division.TotalCount++
	k.setDivision(ctx, division)
}

// decrDivisionTotalCount decrements the sold count of the division
func (k Keeper) decrDivisionTotalCount(ctx sdk.Context, division types.Division) {
	division.TotalCount--
	k.setDivision(ctx, division)
}
