package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgWithdrawNodeReward = "withdraw_node_reward"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgFundCommunityPool{}
	_ sdk.Msg = &MsgWithdrawReward{}
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
	return m.Params.ValidateBasic()
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgFundCommunityPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic executes sanity validation on the provided data
func (m *MsgFundCommunityPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}
	return nil
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (m *MsgFundCommunityPool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Depositor)
	return []sdk.AccAddress{addr}
}

func NewMsgWithdrawNodeReward(nodeOwnerAddr sdk.AccAddress, nodeId string) *MsgWithdrawReward {
	return &MsgWithdrawReward{
		NodeOwnerAddress: nodeOwnerAddr.String(),
		NodeId:           nodeId,
	}
}

func (msg MsgWithdrawReward) Route() string { return ModuleName }
func (msg MsgWithdrawReward) Type() string  { return TypeMsgWithdrawNodeReward }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawReward) GetSigners() []sdk.AccAddress {
	nodeOwner, _ := sdk.AccAddressFromBech32(msg.NodeOwnerAddress)
	return []sdk.AccAddress{nodeOwner}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.NodeOwnerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid node owner address: %s", err)
	}
	return nil
}
