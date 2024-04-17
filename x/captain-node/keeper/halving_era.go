package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) SetHalvingEra(ctx sdk.Context, era uint64) {
	//todo
}

func (k Keeper) GetHalvingEra(ctx sdk.Context) sdk.Dec {
	//todo
	return sdk.NewDec(0)
}
