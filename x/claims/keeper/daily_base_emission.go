package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/tabilabs/tabi/x/mint/types"
)

var (
	ConstantA = sdk.NewDec(int64(3_000_000)) // constant 3,000,000
)

// CalculateCaptainNodeBlockActualEmission calculates the block actual emission of captain node
// BlockActualEmission = DailyActualEmission / BlocksPerDay
func (k Keeper) CalculateCaptainNodeBlockActualEmission(ctx sdk.Context) sdk.Dec {
	// Calculate daily actual emission
	dailyActualEmission := k.CalculateCaptainNodeDailyActualEmission(ctx)
	blocksPerDay := sdk.NewDec(int64(minttypes.BlocksPerDay))
	// Calculate block actual emission
	blockActualEmission := dailyActualEmission.Quo(blocksPerDay)
	return blockActualEmission
}

// CalculateCaptainNodeDailyActualEmission calculates the daily actual emission of captain node
// DailyActualEmission = DailyBaseEmission * (OperationalRate^0.5) * (PledgeRate ^ 0.5)
func (k Keeper) CalculateCaptainNodeDailyActualEmission(ctx sdk.Context) sdk.Dec {
	// Calculate daily base emission
	dailyBaseEmission := k.CalculateCaptainNodeDailyBaseEmission(ctx)
	//OperationalRate
	operationalRate := k.captainNodeKeeper.GetOperationalRate(ctx)
	// CalculatePledgeRate calculates the pledge rate of global
	pledgeRate := k.CalculatePledgeRate(ctx)
	// Calculate daily actual emission
	dailyActualEmission := dailyBaseEmission.Mul(operationalRate.Power(0.5)).Mul(pledgeRate.Power(0.5))
	return dailyActualEmission
}

// CalculateCaptainNodeDailyBaseEmission calculates the daily base emission of captain node
// DailyBaseEmission = ConstantA * TechProgressCoefficient * HalvingEra
func (k Keeper) CalculateCaptainNodeDailyBaseEmission(ctx sdk.Context) sdk.Dec {
	techProgressCoefficient := k.CalculateTechProgressCoefficient(ctx)
	halvingEra := k.CalculateHalvingEra(ctx)
	// Calculate daily base emission
	dailyBaseEmission := ConstantA.Mul(techProgressCoefficient).Mul(halvingEra)
	return dailyBaseEmission
}
