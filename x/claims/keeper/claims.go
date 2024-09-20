package keeper

import (
	tabitypes "github.com/tabilabs/tabi/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
	"github.com/tabilabs/tabi/x/claims/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
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

	if totalRewards.IsZero() {
		return sdk.Coins{}, types.ErrZeroRewards
	}

	// Truncate the rewards
	truncatedCoins, _ := totalRewards.TruncateDecimal()

	// mint vetabi to the module
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, truncatedCoins); err != nil {
		return sdk.Coins{}, nil
	}

	// send the rewards to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Bytes(), truncatedCoins); err != nil {
		return sdk.Coins{}, errorsmod.Wrapf(
			types.ErrSendCoins,
			"error while sending coins from module(%s) to account(%s)",
			types.ModuleName, receiver.String())
	}
	for _, node := range nodes {
		if err := k.captainsKeeper.UpdateGlobalAndNodeClaimedEmission(
			ctx,
			node.Id,
		); err != nil {
			return sdk.Coins{}, errorsmod.Wrapf(
				types.ErrUpdateNodeHistoricalEmissionOnLastClaim,
				"error while updating node(%s) historical emission on last claim", node.Id)
		}
	}

	return truncatedCoins, nil
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

	if epochSequence == 0 {
		return sdk.DecCoins{}, types.ErrFirstEpoch
	}

	// all historical emission
	historicalEmission := k.captainsKeeper.CalcNodeCumulativeEmissionByEpoch(ctx, epochSequence, nodeId)
	// historical emission on last claim
	historicalEmissionOnLastClaim := k.captainsKeeper.GetNodeClaimedEmission(ctx, nodeId)
	reward := historicalEmission.Sub(historicalEmissionOnLastClaim)

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(tabitypes.AttoVeTabi, reward)), nil
}

func (k Keeper) ClaimedRewards(ctx sdk.Context, owner sdk.Address) (sdk.DecCoins, error) {
	// Calculate the rewards for each node
	totalClaimedRewards := sdk.ZeroDec()
	nodes := k.captainsKeeper.GetNodesByOwner(ctx, owner.Bytes())
	// check if the sender has not held node
	if len(nodes) == 0 {
		return sdk.DecCoins{}, types.ErrHolderNotFound
	}

	for _, node := range nodes {
		claimedReward := k.captainsKeeper.GetNodeClaimedEmission(ctx, node.Id)
		// Sum the rewards
		totalClaimedRewards = totalClaimedRewards.Add(claimedReward)
	}

	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(tabitypes.AttoVeTabi, totalClaimedRewards)), nil
}
