package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/evm/types"
)

func (k Keeper) WithdrawRewards(ctx sdk.Context, sender, receiver sdk.Address) (sdk.Coins, error) {
	// calculate the rewards
	reward, err := k.CalculateRewardsByOwner(ctx, sender)
	if err != nil {
		return sdk.Coins{}, err
	}
	// send the rewards to the receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Bytes(), reward); err != nil {
		return sdk.Coins{}, err
	}
	// prune the epochs
	// todo: implement the pruneEpochs function
	return sdk.Coins{}, nil
}

func (k Keeper) CalculateRewardsByOwner(ctx sdk.Context, owner sdk.Address) (sdk.Coins, error) {
	// 1. Get the Node associated with the sender and traverse the epochs associated with the Node
	nodes := k.captainsKeeper.GetNodesByOwner(ctx, owner.Bytes())
	for _, node := range nodes {
		epochs := k.captainsKeeper.GetEpochs(ctx, node.Id, owner)
		// 2. For each epoch, calculate the rewards
		for _, epoch := range epochs {
			// 3. Calculate the rewards
			// 4. Add the rewards to the total rewards
		}
	}

	return sdk.Coins{}, nil
}

func (k Keeper) CalculateRewardsByNodeId(ctx sdk.Context, nodeId string) (sdk.Coins, error) {
	//todo: implement the CalculateRewardsByNodeId function
	return sdk.Coins{}, nil
}
