package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captain-node/types"
)

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

func (k Keeper) GetDivisions(ctx sdk.Context) (divisions []*types.Division) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DivisionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var division types.Division
		k.cdc.MustUnmarshal(iterator.Value(), &division)
		divisions = append(divisions, &division)
	}
	return
}

// GetDivisionTotalSupply returns the number of all sold nodes by the specified division ID
func (k Keeper) GetDivisionTotalSupply(ctx sdk.Context, divisionID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DivisionTotalSupplyKey(divisionID))
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) HasDivision(ctx sdk.Context, divisionID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.DivisionStoreKey(divisionID))
}

func (k Keeper) IsDivisionSoldOut(ctx sdk.Context, divisionID string) bool {
	division, found := k.GetDivision(ctx, divisionID)
	if !found {
		return true
	}
	if k.GetDivisionTotalSupply(ctx, divisionID) >= division.TotalCount {
		return true
	}
	return false
}

func (k Keeper) incrDivisionTotalSupply(ctx sdk.Context, divisionID string) {
	supply := k.GetDivisionTotalSupply(ctx, divisionID) + 1
	k.updateDivisionTotalSupply(ctx, divisionID, supply)
}

func (k Keeper) updateDivisionTotalSupply(ctx sdk.Context, divisionID string, supply uint64) {
	store := ctx.KVStore(k.storeKey)
	supplyKey := types.DivisionTotalSupplyKey(divisionID)
	store.Set(supplyKey, sdk.Uint64ToBigEndian(supply))
}
