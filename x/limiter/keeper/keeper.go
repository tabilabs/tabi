package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tabilabs/tabi/x/limiter/types"
)

// Keeper of the limiter store
type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	authority sdk.AccAddress
}

// NewKeeper creates a new limiter Keeper instance
func NewKeeper(cdc codec.Codec, key storetypes.StoreKey, paramSpace paramtypes.Subspace, authority sdk.AccAddress) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   key,
		paramSpace: paramSpace,
		authority:  authority,
	}
}

// SetParams sets the parameters of the limiter module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams returns the total set of limiter parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
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
	for _, member := range params.WhiteList {
		if member == addr.String() {
			return true
		}
	}
	return false
}

// SetParamsInModule sets the parameters of the limiter module
// NOTE: use x/params instead before sdk v0.47.
func (k Keeper) SetParamsInModule(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)
	return nil
}
