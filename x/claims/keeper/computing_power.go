package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DexE = sdk.NewDec(math.E)
)

// CalculateCaptainNodeComputingPowerRateForXn calculates the computing power rate for captain node
// ComputingPowerRateForXn = ComputingPowerForXn / ComputingPower
func (k Keeper) CalculateCaptainNodeComputingPowerRateForXn(ctx sdk.Context, nodeId string) sdk.Dec {
	computingPowerForXn := k.CalculateCaptainNodeComputingPowerForXn(ctx, nodeId)
	captainNodeComputingPower := k.CalculateCaptainNodeComputingPower(ctx)
	return computingPowerForXn.Quo(captainNodeComputingPower)
}

// CalculateCaptainNodeComputingPower calculates the computing power of captain node
// ComputingPower = sum(ComputingPowerForXn)
func (k Keeper) CalculateCaptainNodeComputingPower(ctx sdk.Context) sdk.Dec {
	nodes := k.captainNodeKeeper.GetNodes(ctx)
	computingPower := sdk.ZeroDec()
	for _, node := range nodes {
		computingPowerForXn := k.CalculateCaptainNodeComputingPowerForXn(ctx, node.Id)
		computingPower = computingPower.Add(computingPowerForXn)
	}
	return computingPower
}

// CalculateCaptainNodeComputingPowerForXn calculates the computing power for captain node
// ComputingPowerForXn = BaseComputingPowerFromLevelForXn * (e^(PledgeRateForXn / 0.5)) * PowerOnForXn
func (k Keeper) CalculateCaptainNodeComputingPowerForXn(ctx sdk.Context, nodeId string) sdk.Dec {
	PowerOnForXn := k.captainNodeKeeper.GetNodePowerOnPeriod(ctx, nodeId)
	node, found := k.captainNodeKeeper.GetNode(ctx, nodeId)
	if !found {
		// todo
		return sdk.ZeroDec()
	}
	// get the node's division
	division, found := k.captainNodeKeeper.GetDivision(ctx, node.DivisionId)
	if !found {
		// todo
		return sdk.ZeroDec()
	}
	baseComputingPowerFromLevelForXn := sdk.NewDec(int64(division.ComputingPower))

	pledgeRateForXn := k.CalculatePledgeRateForXN(ctx, sdk.AccAddress(node.Owner))

	// ComputingPowerForXn = BaseComputingPowerFromLevelForXn * (e^(PledgeRateForXn / 0.5)) * PowerOnForXn
	dexEExponent := pledgeRateForXn.Quo(sdk.NewDec(0.5)).TruncateInt64()
	computingPowerForXn := baseComputingPowerFromLevelForXn.Mul(DexE.Power(uint64(dexEExponent))).Mul(PowerOnForXn)
	return computingPowerForXn

}
