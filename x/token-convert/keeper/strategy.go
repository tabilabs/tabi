package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

// createStrategy sets a strategy with the given name, period, and conversion rate.
func (k Keeper) createStrategy(ctx sdk.Context, name string, period int64, conversionRate sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)

	if k.HasStrategy(ctx, name) {
		return sdkerrors.Wrapf(types.ErrInvalidStrategy, "strategy-%s already exists", name)
	}

	strategy := types.Strategy{
		Name:           name,
		Period:         period,
		ConversionRate: conversionRate,
	}
	bz := k.cdc.MustMarshal(&strategy)
	store.Set(types.StrategyStoreKey([]byte(name)), bz)

	return nil
}

// HasStrategy checks if a strategy exists
func (k Keeper) HasStrategy(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.StrategyStoreKey([]byte(name)))
}

// GetStrategy gets a strategy by strategy id
func (k Keeper) GetStrategy(ctx sdk.Context, name string) (types.Strategy, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StrategyStoreKey([]byte(name)))
	if len(bz) == 0 {
		return types.Strategy{}, false
	}

	var strategy types.Strategy
	k.cdc.MustUnmarshal(bz, &strategy)
	return strategy, true
}

// GetStrategies gets all strategies
func (k Keeper) GetStrategies(ctx sdk.Context) (strategies []types.Strategy) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StrategyKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var strategy types.Strategy
		k.cdc.MustUnmarshal(iterator.Value(), &strategy)
		strategies = append(strategies, strategy)
	}
	return
}
