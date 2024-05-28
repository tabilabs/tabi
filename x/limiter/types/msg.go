package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
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

	seenMap := make(map[string]bool)
	for _, addr := range msg.Params.AllowList {
		if _, ok := seenMap[addr]; ok {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate whitelist address")
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrap(err, "invalid whitelist address")
		}
		seenMap[addr] = true
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	fromAddr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddr}
}
