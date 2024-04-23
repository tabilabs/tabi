package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

// SaveDivision saves a division to the store
func (k Keeper) SaveDivision(ctx sdk.Context, division types.Division) error {
	if k.HasDivision(ctx, division.Id) {
		return errorsmod.Wrap(types.ErrDivisionExists, division.Id)
	}
	bz, err := k.cdc.Marshal(&division)
	if err != nil {
		return errorsmod.Wrap(err, "Marshal division failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DivisionStoreKey(division.Id), bz)
	return nil
}

// UpdateDivision updates a division in the store
func (k Keeper) UpdateDivision(ctx sdk.Context, division types.Division) error {
	if !k.HasDivision(ctx, division.Id) {
		return errorsmod.Wrap(types.ErrDivisionNotExists, division.Id)
	}
	bz, err := k.cdc.Marshal(&division)
	if err != nil {
		return errorsmod.Wrap(err, "Marshal division failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DivisionStoreKey(division.Id), bz)
	return nil
}

// HasDivision checks if the division exists in the store
func (k Keeper) HasDivision(ctx sdk.Context, divisionID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.DivisionStoreKey(divisionID))
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

func (k Keeper) incrDivisionSoldCount(ctx sdk.Context, divisionID string) error {
	panic("impl me")
}

func (k Keeper) incrDivisionTotalCount() {
	panic("impl me")
}

func (k Keeper) decrDivisionTotalCount() {
	panic("impl me")
}

// GetDivisionSoldCount returns the number of all sold nodes by the specified division ID
func (k Keeper) GetDivisionSoldCount(ctx sdk.Context, divisionID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DivisionStoreKey(divisionID))

	var division types.Division
	if len(bz) == 0 {
		return 0
	}
	k.cdc.MustUnmarshal(bz, &division)
	return division.SoldCount
}

// IsDivisionSoldOut checks if the division is sold out
func (k Keeper) IsDivisionSoldOut(ctx sdk.Context, divisionID string) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DivisionStoreKey(divisionID))

	var division types.Division
	if len(bz) == 0 {
		return false
	}
	k.cdc.MustUnmarshal(bz, &division)
	return division.SoldCount == division.InitialSupply
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
