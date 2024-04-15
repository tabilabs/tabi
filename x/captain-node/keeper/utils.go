package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// calculatePowerOnPeriod calculate the proportion of online nodes
// result = proportion / maximumPowerOnPeriod
func calculatePowerOnPeriod(proportion, maximumPowerOnPeriod uint64) sdk.Dec {
	proportionBigint := sdk.NewDec(int64(proportion))
	maximumPowerOnPeriodBigint := sdk.NewDec(int64(maximumPowerOnPeriod))
	result := proportionBigint.Quo(maximumPowerOnPeriodBigint)
	return result
}
