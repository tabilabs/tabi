package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
	"github.com/tabilabs/tabi/x/evm/types"
)

const (
	denom = "avetabi"
)

func (k Keeper) WithdrawRewards(ctx sdk.Context, sender, receiver sdk.Address) (sdk.Coins, error) {
	// 1. Get the Node associated with the sender and traverse the epochs associated with the Node
	nodes := k.captainsKeeper.GetNodesByOwner(ctx, sender.Bytes())
	// calculate the rewards
	reward, err := k.CalculateRewards(ctx, nodes)
	if err != nil {
		return sdk.Coins{}, err
	}
	// send the rewards to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Bytes(), reward); err != nil {
		return sdk.Coins{}, err
	}
	for _, node := range nodes {
		if err := k.captainsKeeper.UpdateNodeHistoricalEmissionOnLastClaim(
			ctx,
			node.Id,
		); err != nil {
			return sdk.Coins{}, err
		}
	}

	return sdk.Coins{}, nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, nodes []captainnodetypes.Node) (sdk.Coins, error) {
	// Calculate the rewards for each node
	totalRewards := sdk.DecCoins{}
	for _, node := range nodes {
		reward, err := k.CalculateRewardsByNodeId(ctx, node.Id)
		if err != nil {
			return sdk.Coins{}, err
		}
		// Sum the rewards
		totalRewards = totalRewards.Add(reward...)
	}
	truncatedCoins, _ := totalRewards.TruncateDecimal()
	return truncatedCoins, nil
}

func (k Keeper) CalculateRewardsByNodeId(ctx sdk.Context, nodeId string) (sdk.DecCoins, error) {
	// Get Current epoch
	epochSequence := k.captainsKeeper.GetCurrentEpoch(ctx) - 1
	historicalEmission := k.captainsKeeper.GetNodeHistoricalEmissionOnEpoch(ctx, epochSequence, nodeId)
	historicalEmissionOnLastClaim := k.captainsKeeper.GetNodeHistoricalEmissionOnLastClaim(ctx, nodeId)
	reward := historicalEmission.Sub(historicalEmissionOnLastClaim)

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom, reward)), nil
}
