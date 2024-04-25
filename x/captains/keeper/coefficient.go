package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	TechProgressCoefficientCardinality = sdk.NewDecWithPrec(16, 1) // constant 1.6
)

// GetTechProgressCoefficient calculates tech progress coefficient
// TechProgressCoefficient = 1.6 ^ (captainNodeSaleLevel - 1)
func (k Keeper) GetTechProgressCoefficient(ctx sdk.Context) sdk.Dec {
	captainNodeSaleLevel := k.GetSaleLevel(ctx)
	// Calculate cardinality raised to the power of captainNodeSaleLevel-1
	techProgressCoefficient := TechProgressCoefficientCardinality.Power(captainNodeSaleLevel - 1)
	return techProgressCoefficient
}

// GetHalvingEraCoefficient returns the tech progress coefficient
func (k Keeper) GetHalvingEraCoefficient(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).HalvingEraCoefficient
}

// GetCaptainsConstant returns the captains constant
func (k Keeper) GetCaptainsConstant(ctx sdk.Context) sdk.Dec {
	cc := k.GetParams(ctx).CaptainsConstant
	return sdk.NewDec(int64(cc))
}
