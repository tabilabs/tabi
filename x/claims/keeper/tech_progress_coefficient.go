package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	TechProgressCoefficientCardinality = sdk.NewDecWithPrec(16, 1) // constant 1.6
)

// CalculateTechProgressCoefficient calculates tech progress coefficient
// TechProgressCoefficient = 1.6 ^ (captainNodeSaleLevel - 1)
func (k Keeper) CalculateTechProgressCoefficient(ctx sdk.Context) sdk.Dec {
	captainNodeSaleLevel := k.captainNodeKeeper.GetSaleLevel(ctx)
	// Calculate cardinality raised to the power of captainNodeSaleLevel-1
	techProgressCoefficient := TechProgressCoefficientCardinality.Power(captainNodeSaleLevel - 1)
	return techProgressCoefficient
}
