package keeper

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

const (
	MinDenomTabi   = "utabi"
	MinDenomVetabi = "uvetabi"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	authKeeper      types.AccountKeeper
	bankKeeper      types.BankKeeper
	instantStrategy string // TODO: make this a module param
}

func NewKeeper(cdc codec.BinaryCodec, key *storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper) Keeper {
	return Keeper{
		storeKey:   *key,
		cdc:        cdc,
		authKeeper: ak,
		bankKeeper: bk,
	}
}

// ConvertTabi converts tabi to vetabi.
func (k Keeper) ConvertTabi(ctx sdk.Context, sender sdk.AccAddress, coin sdk.Coin) error {
	// send tabi to the module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	// burn tabi from the module
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	// mint vetabi to the module
	mintCoin := sdk.NewCoin(MinDenomVetabi, coin.Amount)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
		return err
	}

	// send vetabi from module to the sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(mintCoin)); err != nil {
		return err
	}

	return nil
}

// LockVetabiAndCreateVoucher locks vetabi and creates a voucher for future withdraw.
func (k Keeper) LockVetabiAndCreateVoucher(ctx sdk.Context, sender sdk.AccAddress, strategy types.Strategy, coin sdk.Coin) (string, uint64, error) {
	if strategy.Name == k.instantStrategy {
		err := k.InstantWithdrawVetabi(ctx, sender, coin)
		return "", 0, err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return "", 0, err
	}

	voucherID := k.createVoucher(ctx, sender, strategy.Name, coin)
	k.setVoucherByOwner(ctx, sender, voucherID)

	expiryTime := ctx.BlockTime().Add(time.Duration(strategy.Period) * time.Second).String()

	return expiryTime, voucherID, nil
}

// InstantWithdrawVetabi convert vetabi to tabi with instant strategy.
func (k Keeper) InstantWithdrawVetabi(ctx sdk.Context, sender sdk.AccAddress, coin sdk.Coin) error {
	strategy, found := k.GetStrategy(ctx, k.instantStrategy)
	if !found {
		return sdkerrors.Wrapf(types.ErrStrategyNotFound, "strategy instant not found")
	}

	// send vetabi to the module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	// burn vetabi from the module
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	mintAmt := sdk.NewDecFromInt(coin.Amount).Mul(strategy.ConversionRate).RoundInt()
	mintCoin := sdk.NewCoin(MinDenomTabi, mintAmt)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
		return err
	}

	// send tabi from module to the sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(mintCoin)); err != nil {
		return err
	}

	return nil
}

// WithdrawTabi withdraws tabi according to the given voucher.
func (k Keeper) WithdrawTabi(ctx sdk.Context, sender sdk.AccAddress, voucher types.Voucher) (sdk.Coin, sdk.Coin, error) {
	strategy, found := k.GetStrategy(ctx, voucher.Strategy)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrapf(types.ErrStrategyNotFound, "strategy %s not found", voucher.Strategy)
	}

	// release_ratio = (current_time - created_time) / period
	createdTime := voucher.CreatedTime
	currentTime := ctx.BlockTime().Unix()
	releaseRatio := sdk.NewDec(currentTime - createdTime).Quo(sdk.NewDec(strategy.Period))

	// burnable_vetabi_amt = locked_vatabi_amt * release_ratio
	// returnable_vetabi_amt = locked_vetabi_amt * (1 - release_ratio)
	// withdrawable_tabi_amt = burnable_vetabi_amt * conversion_rate
	lockedVetabiAmt := voucher.Amount.Amount
	burnableVetabiAmt := sdk.NewDecFromInt(lockedVetabiAmt).Mul(releaseRatio).RoundInt()
	returnableVetabiAmt := lockedVetabiAmt.Sub(burnableVetabiAmt)
	withdrawableTabiAmt := sdk.NewDecFromInt(burnableVetabiAmt).Mul(strategy.ConversionRate).RoundInt()

	withdrawableTabi := sdk.NewCoin(MinDenomTabi, withdrawableTabiAmt)
	burnableVetabi := sdk.NewCoin(MinDenomVetabi, burnableVetabiAmt)
	returnableVetabi := sdk.NewCoin(MinDenomVetabi, returnableVetabiAmt)

	// make sure the owner has the minimum amount of tabi to withdraw
	if !withdrawableTabi.IsPositive() {
		return sdk.Coin{}, sdk.Coin{}, sdkerrors.Wrapf(types.ErrInsufficientFunds, "insufficient tabi to withdraw")
	}

	// 1-a. mint withdrawable tabi to the module account
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(withdrawableTabi))
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// 1-b. send withdrawable tabi to the owner
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(withdrawableTabi))
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// 2. burn burnable vetabi from the module account
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burnableVetabi))
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// check if there is any returnable vetabi
	if !returnableVetabiAmt.IsPositive() {
		return withdrawableTabi, sdk.Coin{}, nil
	}

	// 3. send returnable vetabi to the owner
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(returnableVetabi))
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	return withdrawableTabi, returnableVetabi, nil
}
