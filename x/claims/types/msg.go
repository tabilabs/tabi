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
	_ sdk.Msg = &MsgWithdrawNodeReward{}
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

func NewMsgWithdrawNodeReward(nodeOwnerAddr sdk.AccAddress, nodeId string) *MsgWithdrawNodeReward {
	return &MsgWithdrawNodeReward{
		NodeOwnerAddress: nodeOwnerAddr.String(),
		NodeId:           nodeId,
	}
}

func (msg MsgWithdrawNodeReward) Route() string { return ModuleName }
func (msg MsgWithdrawNodeReward) Type() string  { return TypeMsgWithdrawNodeReward }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawNodeReward) GetSigners() []sdk.AccAddress {
	nodeOwner, _ := sdk.AccAddressFromBech32(msg.NodeOwnerAddress)
	return []sdk.AccAddress{nodeOwner}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawNodeReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawNodeReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.NodeOwnerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid node owner address: %s", err)
	}
	return nil
}
