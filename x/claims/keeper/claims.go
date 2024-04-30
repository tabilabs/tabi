package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tabitypes "github.com/tabilabs/tabi/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
	"github.com/tabilabs/tabi/x/claims/types"
)

func (k Keeper) WithdrawRewards(ctx sdk.Context, sender, receiver sdk.Address) (sdk.Coins, error) {
	// Get the Node associated with the sender and traverse the epochs associated with the Node
	nodes := k.captainsKeeper.GetNodesByOwner(ctx, sender.Bytes())
	// check if the sender has not held node
	if len(nodes) == 0 {
		return sdk.Coins{}, types.ErrHolderNotFound
	}
	// calculate the rewards
	totalRewards, err := k.CalculateRewards(ctx, nodes)
	if err != nil {
		return sdk.Coins{}, types.ErrCalculateRewards
	}

	// Truncate the rewards
	truncatedCoins, _ := totalRewards.TruncateDecimal()
	// send the rewards to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Bytes(), truncatedCoins); err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(
			types.ErrSendCoins,
			"error while sending coins from module(%s) to account(%s)",
			types.ModuleName, receiver.String())
	}
	for _, node := range nodes {
		if err := k.captainsKeeper.UpdateNodeHistoricalEmissionOnLastClaim(
			ctx,
			node.Id,
		); err != nil {
			return sdk.Coins{}, errorsmod.Wrapf(
				types.ErrUpdateNodeHistoricalEmissionOnLastClaim,
				"error while updating node(%s) historical emission on last claim", node.Id)
		}
	}

	return sdk.Coins{}, nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, nodes []captainnodetypes.Node) (sdk.DecCoins, error) {
	// Calculate the rewards for each node
	totalRewards := sdk.DecCoins{}
	for _, node := range nodes {
		reward, err := k.CalculateRewardsByNodeId(ctx, node.Id)
		if err != nil {
			return sdk.DecCoins{}, err
		}
		// Sum the rewards
		totalRewards = totalRewards.Add(reward...)
	}

	return totalRewards, nil
}

func (k Keeper) CalculateRewardsByNodeId(ctx sdk.Context, nodeId string) (sdk.DecCoins, error) {
	// Get Current epoch
	epochSequence := k.captainsKeeper.GetCurrentEpoch(ctx) - 1
	historicalEmission := k.captainsKeeper.CalAndGetNodeHistoricalEmissionOnEpoch(ctx, epochSequence, nodeId)
	historicalEmissionOnLastClaim := k.captainsKeeper.GetNodeHistoricalEmissionOnLastClaim(ctx, nodeId)
	reward := historicalEmission.Sub(historicalEmissionOnLastClaim)

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(tabitypes.AttoVeTabi, reward)), nil
}
