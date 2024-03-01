package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/claims/types"
)

// GetParamSet returns inflation params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParamSet set inflation params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
