package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/mint/types"
)

// ______________________________________________________________________

// GetMinter returns the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("Stored minter should not have been nil")
	}
	k.cdc.MustUnmarshal(b, &minter)
	return
}

// SetMinter set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		k.feeCollectorName,
		coins,
	)
}
