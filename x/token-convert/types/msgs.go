package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// token-convert module message types
const (
	TypeMsgConvertTabi   = "convert_tabi"
	TypeMsgConvertVetabi = "convert_vetabi"
	TypeMsgWithdrawTabi  = "withdraw_tabi"
	TypeMsgCancelConvert = "cancel_convert"
)

// NOTE: we don't impl legacy msg anymore
var (
	_ sdk.Msg = &MsgConvertTabi{}
	_ sdk.Msg = &MsgConvertVetabi{}
	_ sdk.Msg = &MsgWithdrawTabi{}
	_ sdk.Msg = &MsgCancelConvert{}
)

// NewMsgConvertTabi is a constructor function for MsgConvertTabi
func NewMsgConvertTabi(coin sdk.Coin, sender sdk.AccAddress) *MsgConvertTabi {
	return &MsgConvertTabi{
		Coin:   coin,
		Sender: sender.String(),
	}
}

func (m *MsgConvertTabi) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", m.Sender)
	}

	if !m.Coin.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin amount: %s", m.Coin.Amount)
	}

	if m.Coin.Denom != MinDenomTabi {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin denom: %s", m.Coin.Denom)
	}

	return nil
}

func (m *MsgConvertTabi) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgConvertVetabi is a constructor function for MsgConvertVetabi
func NewMsgConvertVetabi(coin sdk.Coin, sender sdk.AccAddress, strategy string) *MsgConvertVetabi {
	return &MsgConvertVetabi{
		Coin:     coin,
		Sender:   sender.String(),
		Strategy: strategy,
	}
}

func (m *MsgConvertVetabi) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", m.Sender)
	}

	if !m.Coin.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin amount: %s", m.Coin.Amount)
	}

	if m.Coin.Denom != MinDenomVetabi {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin denom: %s", m.Coin.Denom)
	}

	if len(m.Strategy) == 0 {
		return errorsmod.Wrap(ErrInvalidStrategy, "strategy is empty")
	}

	return nil
}

func (m *MsgConvertVetabi) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgWithdrawTabi is a constructor function for MsgWithdrawTabi
func NewMsgWithdrawTabi(voucherId string, sender sdk.AccAddress) *MsgWithdrawTabi {
	return &MsgWithdrawTabi{
		VoucherId: voucherId,
		Sender:    sender.String(),
	}
}

func (m *MsgWithdrawTabi) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", m.Sender)
	}

	if len(m.VoucherId) == 0 {
		return errorsmod.Wrapf(ErrInvalidVoucher, "voucher id is empty")
	}

	return nil
}

func (m *MsgWithdrawTabi) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgCancelConvert is a constructor function for MsgCancelConvert
func NewMsgCancelConvert(voucherId string, sender sdk.AccAddress) *MsgCancelConvert {
	return &MsgCancelConvert{
		VoucherId: voucherId,
		Sender:    sender.String(),
	}
}

func (m *MsgCancelConvert) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", m.Sender)
	}

	if len(m.VoucherId) == 0 {
		return errorsmod.Wrapf(ErrInvalidVoucher, "voucher id is empty")
	}

	return nil
}

func (m *MsgCancelConvert) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
