package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/limiter/types"
)

// Keeper of the limiter store
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	authority sdk.AccAddress
}

// NewKeeper creates a new limiter Keeper instance
func NewKeeper(cdc codec.Codec, key storetypes.StoreKey, authority sdk.AccAddress) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		authority: authority,
	}
}

// SetParams sets the parameters of the limiter module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetParams returns the total set of limiter parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.Params{}
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// IsEnabled returns the enabled status of the limiter module
func (k Keeper) IsEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.Enabled
}

// IsAuthorized checks if the addr is in white list.
func (k Keeper) IsAuthorized(ctx sdk.Context, addr sdk.AccAddress) bool {
	params := k.GetParams(ctx)
	for _, member := range params.AllowList {
		if member == addr.String() {
			return true
		}
	}
	return false
}
