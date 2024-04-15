package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// CalculateHalvingEra calculates HalvingEra
func (k Keeper) CalculateHalvingEra(ctx sdk.Context) sdk.Dec {
	// get daily issuance
	dailyIssuance := k.mintKeeper.GetDailyIssuance(ctx)
	// get current run days
	currentRunDays := k.GetCurrentRunDays(ctx)
	currentRunDaysInt := sdk.NewInt(int64(currentRunDays))

	// get captain node sequence
	captainNodeSequence := k.captainNodeKeeper.GetNodeSequence(ctx) - 1
	captainNodeSequenceInt := sdk.NewDec(int64(captainNodeSequence))
	// calculate current issuance amount
	currentIssuanceAmount := dailyIssuance.MulInt(currentRunDaysInt)

	// Compute lvalue of comparison
	coefficientLvalue := sdk.NewDecWithPrec(25, 2) // constant 0.25
	lvalue := currentIssuanceAmount.Mul(coefficientLvalue).Add(captainNodeSequenceInt)

	claimsInitRewardRate := sdk.NewDecWithPrec(43, 2) // constant 0.43

	// get minter
	minter := k.mintKeeper.GetMinter(ctx)
	inflationBase := sdk.NewDecFromInt(minter.InflationBase)

	// Compute rvalue of comparison
	coefficientRvalue1 := sdk.NewDecWithPrec(50, 2) // constant 0.5
	rvalue1 := inflationBase.Mul(claimsInitRewardRate).Mul(coefficientRvalue1)

	coefficientRvalue2 := sdk.NewDecWithPrec(75, 2) // constant 0.75
	rvalue2 := inflationBase.Mul(claimsInitRewardRate).Mul(coefficientRvalue2)

	// Compare lvalue and rvalue1 and rvalue2
	// if lvalue > rvalue1, return 0.5
	// if lvalue > rvalue2, return 0.25
	// else return 1
	if lvalue.GT(rvalue1) {
		if lvalue.GT(rvalue2) {
			// return 0.25
			return sdk.NewDecWithPrec(25, 2)
		}
		// return 0.5
		return sdk.NewDecWithPrec(5, 1)
	}

	return sdk.NewDec(1)
}

func (k Keeper) GetCurrentRunDays(ctx sdk.Context) uint64 {
	// todo
	return 0
}
