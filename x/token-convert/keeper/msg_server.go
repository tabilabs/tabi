package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tabilabs/tabi/x/token-convert/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the token-convert MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}

// ConvertTabi converts tabi to vetabi.
func (m msgServer) ConvertTabi(goCtx context.Context, msg *types.MsgConvertTabi) (*types.MsgConvertTabiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if the sender has enough coins
	balance := m.bankKeeper.GetBalance(ctx, sender, MinDenomTabi)
	_, hasNeg := sdk.Coins{balance}.SafeSub(msg.Coin)
	if hasNeg {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance: %s%s", balance, MinDenomTabi)
	}

	// execute conversion
	err = m.Keeper.ConvertTabi(ctx, sender, msg.Coin)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConvertTabi,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Coin.Amount.String()),
		),
	)

	return &types.MsgConvertTabiResponse{}, nil
}

// ConvertVetabi handles user's covert vetabi to tabi request and returns a voucher for future redeem.
func (m msgServer) ConvertVetabi(goCtx context.Context, msg *types.MsgConvertVetabi) (*types.MsgConvertVetabiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check balances
	balance := m.bankKeeper.GetBalance(ctx, sender, MinDenomTabi)
	_, hasNeg := sdk.Coins{balance}.SafeSub(msg.Coin)
	if hasNeg {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance: %s%s", balance, MinDenomVetabi)
	}

	// get strategy
	strategy, found := m.GetStrategy(ctx, msg.Strategy)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrStrategyNotFound, "strategy %s not found", msg.Strategy)
	}

	expiryTime, voucherID, err := m.LockVetabiAndCreateVoucher(ctx, sender, strategy, msg.Coin)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConvertVetabi,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Coin.String()),
			sdk.NewAttribute(types.AttributeKeyVoucherID, strconv.FormatUint(voucherID, 10)),
			sdk.NewAttribute(types.AttributeKeyExpiryTime, expiryTime),
		),
	)

	return &types.MsgConvertVetabiResponse{
		VoucherId:  voucherID,
		ExpiryTime: expiryTime,
	}, nil
}

func (m msgServer) WithdrawTabi(goCtx context.Context, msg *types.MsgWithdrawTabi) (*types.MsgWithdrawTabiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	voucher, found := m.GetVoucher(ctx, msg.VoucherId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrVoucherNotFound, "voucher %d not found", msg.VoucherId)
	}

	if voucher.Owner != msg.Sender {
		return nil, sdkerrors.Wrapf(types.ErrInvalidVoucherOwner, "voucher %d is not owned by %s", msg.VoucherId, msg.Sender)
	}

	// execute withdrawal
	tabiWithdrawn, vetabiReturned, err := m.Keeper.WithdrawTabi(ctx, sender, voucher)
	if err != nil {
		return nil, err
	}

	// delete voucher
	m.deleteVoucher(ctx, msg.VoucherId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawTabi,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyAmount, tabiWithdrawn.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, vetabiReturned.String()),
		),
	)

	return &types.MsgWithdrawTabiResponse{
		TabiWithdrawn:  tabiWithdrawn,
		VetabiReturned: vetabiReturned,
	}, nil
}

// CancelConvert cancels the conversion and returns the locked token to the sender.
func (m msgServer) CancelConvert(goCtx context.Context, msg *types.MsgCancelConvert) (*types.MsgCancelConvertResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	voucher, found := m.GetVoucher(ctx, msg.VoucherId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrVoucherNotFound, "voucher %d not found", msg.VoucherId)
	}

	if voucher.Owner != msg.Sender {
		return nil, sdkerrors.Wrapf(types.ErrInvalidVoucherOwner, "voucher %d is not owned by %s", msg.VoucherId, msg.Sender)
	}

	moduleAcc := m.authKeeper.GetModuleAddress(types.ModuleName)
	balance := m.bankKeeper.GetBalance(ctx, moduleAcc, MinDenomVetabi)
	_, hasNeg := sdk.Coins{balance}.SafeSub(voucher.Amount)
	if hasNeg {
		// WARN: this error shall never happen
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient balance: %s%s", balance, MinDenomVetabi)
	}

	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(voucher.Amount))
	if err != nil {
		return nil, err
	}

	// delete voucher
	m.deleteVoucher(ctx, msg.VoucherId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelConvert,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyAmount, voucher.Amount.String()),
		),
	)

	return &types.MsgCancelConvertResponse{
		VetabiUnlocked: voucher.Amount,
	}, nil
}
