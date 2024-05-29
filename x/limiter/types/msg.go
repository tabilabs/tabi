package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgLimiterSwitch{}
	_ sdk.Msg = &MsgAddAllowListMember{}
	_ sdk.Msg = &MsgRemoveAllowListMember{}
)

// NewMsgUpdateParams defines a message to update the params of the limiter module
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	return ValidateParams(&msg.Params)
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgLimiterSwitch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgLimiterSwitch) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgAddAllowListMember) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrap(err, "invalid member address")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgAddAllowListMember) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddr}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRemoveAllowListMember) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errorsmod.Wrap(err, "invalid member address")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgRemoveAllowListMember) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddr}
}
