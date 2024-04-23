package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	TechProgressCoefficientCardinality = sdk.NewDecWithPrec(16, 1) // constant 1.6
)

// CalculateTechProgressCoefficient calculates tech progress coefficient
// TechProgressCoefficient = 1.6 ^ (captainNodeSaleLevel - 1)
func (k Keeper) CalculateTechProgressCoefficient(ctx sdk.Context) sdk.Dec {
	captainNodeSaleLevel := k.GetSaleLevel(ctx)
	// Calculate cardinality raised to the power of captainNodeSaleLevel-1
	techProgressCoefficient := TechProgressCoefficientCardinality.Power(captainNodeSaleLevel - 1)
	return techProgressCoefficient
}

func (k Keeper) GetHalvingEraCoefficient(ctx sdk.Context) sdk.Dec {
	panic("implement me")
}

func (k Keeper) SetHalvingEraCoefficient(ctx sdk.Context, era uint64) {
	panic("implement me")
}
