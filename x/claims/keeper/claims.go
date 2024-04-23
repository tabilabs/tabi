package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/evm/types"
)

func (k Keeper) WithdrawRewards(ctx sdk.Context, sender, receiver sdk.Address) (sdk.Coins, error) {
	// todo
	// send the rewards to the receiver
	reward, err := k.CalculateRewards(ctx, sender)
	if err != nil {
		return sdk.Coins{}, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver.Bytes(), reward); err != nil {
		return sdk.Coins{}, err
	}
	// prune the epochs
	return sdk.Coins{}, nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, sender sdk.Address) (sdk.Coins, error) {
	// todo
	// 1. Get the Node associated with the sender and traverse the epochs associated with the Node
	epochs := k.captainsKeeper.GetEpochs(ctx, sender)
	// 2. For each epoch, calculate the rewards

	for _, epoch := range epochs {
		// 3. Calculate the rewards
		// 4. Add the rewards to the total rewards
	}

	return sdk.Coins{}, nil
}
