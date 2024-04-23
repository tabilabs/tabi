package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func PledgeRatioGlobal(ctx sdk.Context, epochID string) sdk.Dec {
	panic("implement me")
}

func PledgeRatioByNode(ctx sdk.Context, nodeID string, epochID string) {
	panic("implement me")
}

func PledgeRatioByOwner(ctx sdk.Context, owner sdk.AccAddress, epochID string) {
	panic("implement me")
}

// PledgeSum returns the total pledge amount of captains' owners on the epoch end.
func PledgeSum(ctx sdk.Context, epochID string) sdk.Dec {
	panic("implement me")
}

func PledgeByOwner(ctx sdk.Context, owner sdk.AccAddress, epochID string) sdk.Dec {
	panic("implement me")
}

// TODO: legacy functions

// CalculatePledgeRateForXN calculates the pledge rate of the owner
// pledgeRate = pledgeCoin / mintCoin
func (k Keeper) CalculatePledgeRateForXN(ctx sdk.Context, owner sdk.AccAddress) sdk.Dec {
	pledgeCoin := k.calculatePledgeTotalCountFromForXn(ctx, owner)
	mintCoin := k.calculateMintTotalCountFromXN(ctx, owner)
	return sdk.NewDecFromInt(pledgeCoin.Amount).Quo(sdk.NewDecFromInt(mintCoin.Amount))
}

// calculatePledgeTotalCountFromForXn calculates the total pledge amount of the owner
func (k Keeper) calculatePledgeTotalCountFromForXn(ctx sdk.Context, owner sdk.AccAddress) sdk.Coin {
	stakingParams := k.stakingKeeper.GetParams(ctx)

	// WARN: Can it be delegated to a candidate validator?
	maxRetrieve := stakingParams.GetMaxValidators()
	delegations := k.stakingKeeper.GetDelegatorDelegations(ctx, owner, uint16(maxRetrieve))
	totolBalance := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdk.ZeroInt())
	for _, delegation := range delegations {
		val, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
		if !found {
			//todo
			continue
		}
		delegationBalance := sdk.NewCoin(
			k.stakingKeeper.BondDenom(ctx),
			val.TokensFromShares(delegation.Shares).TruncateInt(),
		)
		totolBalance = totolBalance.Add(delegationBalance)
	}
	return totolBalance
}

// calculateMintTotalCountFromXN calculates the total mint amount of the owner
// WARN: This function is not implemented yet
// WARN: mintAmount can be zero
func (k Keeper) calculateMintTotalCountFromXN(ctx sdk.Context, owner sdk.AccAddress) sdk.Coin {
	// todo
	return sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdk.ZeroInt())
}
