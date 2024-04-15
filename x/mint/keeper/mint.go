package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/mint/types"
)

// ______________________________________________________________________

// GetMinter returns the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("Stored minter should not have been nil")
	}
	k.cdc.MustUnmarshal(b, &minter)
	return
}

// SetMinter set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, coins)
}

// AllocateExponentialInflation allocates coins from the inflation to external
// modules according to allocation proportions:
//   - staking rewards -> sdk `auth` module fee collector
//   - claims rewards -> `x/claims` module
func (k Keeper) AllocateExponentialInflation(
	ctx sdk.Context,
	mintedCoin sdk.Coin,
	params types.Params,
) error {
	distribution := params.InflationDistribution

	// Allocate staking rewards into fee collector account
	stakingRewards := sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.StakingRewards)}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		k.feeCollectorName,
		stakingRewards,
	); err != nil {
		return err
	}

	// Allocate claims rewards to claims module account
	claimsRewards := sdk.Coins{k.GetProportions(ctx, mintedCoin, distribution.ClaimsRewards)}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(
		ctx,
		types.ModuleName,
		k.claimsCollectorName,
		claimsRewards,
	); err != nil {
		return err
	}
	// Allocate community pool amount (remaining module balance) to community
	// pool address
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	inflationBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddr)
	if err := k.distrKeeper.FundCommunityPool(
		ctx,
		inflationBalance,
		moduleAddr,
	); err != nil {
		return err
	}

	return nil
}

// GetProportions calculates the proportion of coins that is to be
// allocated during inflation for a given distribution.
func (k Keeper) GetProportions(
	_ sdk.Context,
	coin sdk.Coin,
	distribution sdk.Dec,
) sdk.Coin {
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: sdk.NewDecFromInt(coin.Amount).Mul(distribution).TruncateInt(),
	}
}

// GetTokenSupply returns the total supply of the token
func (k Keeper) GetTokenSupply(ctx sdk.Context) sdkmath.Int {
	tabiCoin := k.bankKeeper.GetSupply(ctx, "")
	sdTabiCoin := k.bankKeeper.GetSupply(ctx, "")
	tabiCoin = tabiCoin.Add(sdTabiCoin)
	return tabiCoin.Amount
}

// GetDailyIssuance returns the amount of tokens issued per day
func (k Keeper) GetDailyIssuance(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)

	return params.Inflation.MulInt(types.InitialIssue)
}
