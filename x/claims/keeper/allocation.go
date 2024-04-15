package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/claims/types"
)

// AllocateTokens performs reward and fee distribution to all captain node
func (k Keeper) AllocateTokens(ctx sdk.Context) {

	// Calculate block actual emission of captain node
	totalRewards := k.CalculateCaptainNodeBlockActualEmission(ctx)

	totalRewardsCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), totalRewards.RoundInt())
	totalRewardsInt := sdk.NewCoins(totalRewardsCoin)
	totalRewardsDecCoins := sdk.NewDecCoinsFromCoins(totalRewardsInt...)
	// transfer rewards to claims module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ClaimsCollectorName,
		types.ModuleName,
		totalRewardsInt,
	)
	if err != nil {
		panic(err)
	}

	// Get All captain nodes
	captainNodes := k.captainNodeKeeper.GetNodes(ctx)

	remaining := totalRewardsDecCoins
	// Allocate tokens to captain nodes
	for _, captainNode := range captainNodes {
		// Calculate captain node reward
		captainNodeBlockReward := k.CalculateCaptainNodeBlockReward(
			ctx,
			captainNode.Id,
			totalRewardsInt,
		)
		captainNodeOwner := sdk.AccAddress(captainNode.Owner)

		captainNodeBlockRewardDecCoins := sdk.NewDecCoinsFromCoins(captainNodeBlockReward...)

		// Allocate Tokens To Captain Node Of Owner
		k.AllocateTokensToCaptainNodeOfOwner(
			ctx,
			captainNodeOwner,
			sdk.NewDecCoinsFromCoins(captainNodeBlockReward...),
		)
		remaining = remaining.Sub(captainNodeBlockRewardDecCoins)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	feePool := k.GetFeePool(ctx)

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
}

func (k Keeper) AllocateTokensToCaptainNodeOfOwner(ctx sdk.Context, owner sdk.AccAddress, reward sdk.DecCoins) {
	// Allocate tokens to captain node of owner
	// todo
}
