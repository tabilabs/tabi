package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// caption-node message types
const (
	TypeMsgMint                = "mint"
	TypeMsgReceiveExperience   = "receive_experience"
	TypeMsgUpdatePowerOnPeriod = "update_power_on_period"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgReceiveExperience{}
	_ sdk.Msg = &MsgUpdatePowerOnPeriod{}
)

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic executes sanity validation on the provided data
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}
	return m.Params.Validate()
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgMint() *MsgMint {
	return &MsgMint{}
}

// Route Implements Msg.
func (m MsgMint) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgMint) Type() string { return TypeMsgMint }

// ValidateBasic Implements Msg.
func (msg MsgMint) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgReceiveExperience() *MsgReceiveExperience {
	return &MsgReceiveExperience{}
}

// Route Implements Msg.
func (m MsgReceiveExperience) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgReceiveExperience) Type() string { return TypeMsgReceiveExperience }

// ValidateBasic Implements Msg.
func (msg MsgReceiveExperience) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgReceiveExperience) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgReceiveExperience) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgUpdatePowerOnPeriod() *MsgUpdatePowerOnPeriod {
	return &MsgUpdatePowerOnPeriod{}
}

// Route Implements Msg.
func (m MsgUpdatePowerOnPeriod) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgUpdatePowerOnPeriod) Type() string { return TypeMsgUpdatePowerOnPeriod }

// ValidateBasic Implements Msg.
func (msg MsgUpdatePowerOnPeriod) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdatePowerOnPeriod) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdatePowerOnPeriod) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}
