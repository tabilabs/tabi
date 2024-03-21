package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthCaller asserts whether a Callers is already registered
func (k Keeper) AuthCaller(ctx sdk.Context, callerAddr sdk.AccAddress) bool {
	params := k.GetParams(ctx)

	allowCall := false
	for _, caller := range params.Callers {
		if sdk.MustAccAddressFromBech32(caller).Equals(callerAddr) {
			allowCall = true
		}
	}
	return allowCall
}
