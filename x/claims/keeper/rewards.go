package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// CalculateCaptainNodeBlockReward calculates the block reward for a captain node
// The block reward is calculated as follows:
func (k Keeper) CalculateCaptainNodeBlockReward(
	ctx sdk.Context,
	nodeId string,
	totalRewards sdk.Coins,
) sdk.Coins {
	computingPowerRateForXnode := k.CalculateCaptainNodeComputingPowerRateForXn(ctx, nodeId)
	captainNodeBlockReward := totalRewards.MulInt(computingPowerRateForXnode.RoundInt())
	return captainNodeBlockReward
}
