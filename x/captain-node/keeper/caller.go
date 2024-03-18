package keeper

import (
	"github.com/tabilabs/tabi/x/captain-node/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterCallers saves the callers
func (k Keeper) RegisterCallers(ctx sdk.Context, callers []string) {
	store := ctx.KVStore(k.storeKey)
	ic := &types.IdentifiedCallers{Callers: callers}
	store.Set([]byte(types.KeyCallers), k.cdc.MustMarshal(ic))
}

func (k Keeper) GetCallers(ctx sdk.Context) (callers []string) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyCallers))

	var ic = &types.IdentifiedCallers{}
	k.cdc.MustUnmarshal(bz, ic)
	return ic.Callers
}

// GetAllCallers returns all registered Callers addresses
func (k Keeper) GetAllCallers(ctx sdk.Context) (callers []types.IdentifiedCallers) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.KeyCallers))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ir = &types.IdentifiedCallers{}
		k.cdc.MustUnmarshal(iterator.Value(), ir)
		callers = append(callers, *ir)
	}
	return
}

// AuthCaller asserts whether a Callers is already registered
func (k Keeper) AuthCaller(ctx sdk.Context, caller string) bool {
	for _, r := range k.GetCallers(ctx) {
		if r == caller {
			return true
		}
	}
	return false
}
