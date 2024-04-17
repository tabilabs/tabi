package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) GetDailyIssuance(ctx sdk.Context) sdk.Dec {
	// todo
	return sdk.ZeroDec()
}

func (k Keeper) SetDailyIssuance(ctx sdk.Context, dailyTabiIssuance sdk.Dec) {
	// todo
}

func (k Keeper) GetDailyBaseEmission(ctx sdk.Context) sdk.Dec {
	// todo
	return sdk.ZeroDec()
}

func (k Keeper) SetDailyBaseEmission(ctx sdk.Context, dailyBaseEmission sdk.Dec) {
	// todo
}
