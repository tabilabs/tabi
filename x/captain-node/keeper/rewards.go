package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetNodeBlockReward(ctx sdk.Context) {

}

func (k Keeper) CalculateDailyIssuance(ctx sdk.Context) {

}

func (k Keeper) CalculateHalvingEra(ctx sdk.Context) {

}

func (k Keeper) CalculateTechProgressCoefficient(ctx sdk.Context) {

}

func (k Keeper) CalculateOperationalRate(ctx sdk.Context) {

}

func (k Keeper) CalculatePledgeTotalCount(ctx sdk.Context) {

}

func (k Keeper) CalculatePledgeRate(ctx sdk.Context) {

}

func (k Keeper) CalculateNodeBlockActualEmission(ctx sdk.Context) {

}

func (k Keeper) CalculatePowerOnPeriod(ctx sdk.Context, proportion uint64) *big.Int {
	params := k.GetParams(ctx)
	proportionBigInt := new(big.Int).SetUint64(params.MaximumPowerOnPeriod)
	maximumPowerOnPeriodBigInt := new(big.Int).SetUint64(proportion)
	result := sdk.NewDecFromBigInt(proportionBigInt).Quo(sdk.NewDecFromBigInt(maximumPowerOnPeriodBigInt))
	return result.BigInt()
}
