package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) HasNode(ctx sdk.Context, sender sdk.AccAddress) bool {
	if k.captainsKeeper.GetUserHoldingQuantity(ctx, sender) > 0 {
		return true
	}
	return false
}
