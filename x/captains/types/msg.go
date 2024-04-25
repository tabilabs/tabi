package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// captions message types
const (
	TypeMsgCreateCaptainNode      = "create_captain_node"
	TypeMsgCommitReport           = "commit_report"
	TypeMsgAddAuthorizedMember    = "add_authorized_member"
	TypeMsgRemoveAuthorizedMember = "remove_authorized_member"
	TypeMsgUpdateSaleLevel        = "update_sale_level"
	TypeMsgCommitComputingPower   = "commit_computing_power"
	TypeMsgClaimComputingPower    = "claim_computing_power"
)

var (
	_ sdk.Msg = &MsgCreateCaptainNode{}
	_ sdk.Msg = &MsgCommitReport{}
	_ sdk.Msg = &MsgAddAuthorizedMembers{}
	_ sdk.Msg = &MsgRemoveAuthorizedMembers{}
	_ sdk.Msg = &MsgUpdateSaleLevel{}
	_ sdk.Msg = &MsgCommitComputingPower{}
	_ sdk.Msg = &MsgClaimComputingPower{}
)

func NewMsgCreateCaptainNode(authority, owner, divisionId string) *MsgCreateCaptainNode {
	return &MsgCreateCaptainNode{
		Authority:  authority,
		Owner:      owner,
		DivisionId: divisionId,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgCreateCaptainNode) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrap(err, "invalid owner address")
	}
	if len(msg.DivisionId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "division id cannot be empty")
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgCreateCaptainNode) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgCommitReport(authority string, report []byte) *MsgCommitReport {
	// TODO: fixme
	return &MsgCommitReport{
		Authority: authority,
		Report:    nil,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgCommitReport) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")

	}

	// FIXME: fix here after designing the report structure
	panic("implement full validation!")

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgCommitReport) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

func NewAddAuthorizedMembers() *MsgAddAuthorizedMembers {
	return &MsgAddAuthorizedMembers{}
}

// ValidateBasic Implements Msg.
func (msg *MsgAddAuthorizedMembers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if len(msg.Members) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "member cannot be empty")
	}

	for _, member := range msg.Members {
		if _, err := sdk.AccAddressFromBech32(member); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
	}

	return nil
}

func (msg *MsgAddAuthorizedMembers) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgRemoveAuthorizedMembers() *MsgRemoveAuthorizedMembers {
	return &MsgRemoveAuthorizedMembers{}
}

// ValidateBasic Implements Msg.
func (msg *MsgRemoveAuthorizedMembers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if len(msg.Members) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "member cannot be empty")
	}

	for _, caller := range msg.Members {
		if _, err := sdk.AccAddressFromBech32(caller); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgRemoveAuthorizedMembers) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgUpdateSaleLevel() *MsgUpdateSaleLevel {
	return &MsgUpdateSaleLevel{}
}

// ValidateBasic Implements Msg.
func (msg *MsgUpdateSaleLevel) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.SaleLevel <= 1 || msg.SaleLevel > 5 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "sale level must be between 1 and 5")
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgUpdateSaleLevel) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgCommitComputingPower(
	rewards []ClaimableComputingPower,
	authority string,
) *MsgCommitComputingPower {
	return &MsgCommitComputingPower{
		Authority:             authority,
		ComputingPowerRewards: rewards,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgCommitComputingPower) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")

	}
	if len(msg.ComputingPowerRewards) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "computing powers to commit cannot be empty")
	}

	for _, reward := range msg.ComputingPowerRewards {
		if reward.Amount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "computing power cannot be zero")
		}
		if _, err := sdk.AccAddressFromBech32(reward.Owner); err != nil {
			return errorsmod.Wrap(err, "invalid receiver address")
		}
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgCommitComputingPower) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgWithdrawComputingPower(nodeId string, computingPowerAmount uint64, sender string) *MsgClaimComputingPower {
	return &MsgClaimComputingPower{
		NodeId:               nodeId,
		ComputingPowerAmount: computingPowerAmount,
		Sender:               sender,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgClaimComputingPower) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if len(msg.NodeId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node id cannot be empty")
	}
	if msg.ComputingPowerAmount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "experience cannot be zero")
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgClaimComputingPower) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}
