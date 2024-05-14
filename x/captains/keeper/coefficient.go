package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcTechProgressCoefficient calculates tech progress coefficient
// TechProgressCoefficient = 1.6 ^ (captainNodeSaleLevel - 1)
func (k Keeper) CalcTechProgressCoefficient(ctx sdk.Context) sdk.Dec {
	captainNodeSaleLevel := k.GetSaleLevel(ctx)
	cardinality := k.GetParams(ctx).TechProgressCoefficientCardinality
	return cardinality.Power(captainNodeSaleLevel - 1)
}

// GetHalvingEraCoefficient returns the tech progress coefficient
func (k Keeper) GetHalvingEraCoefficient(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).HalvingEraCoefficient
}

// GetCaptainsConstant returns the captains constant
func (k Keeper) GetCaptainsConstant(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).CaptainsConstant
}
