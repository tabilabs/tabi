package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// caption-node message types
const (
	TypeMsgCreateCaptainNode    = "mint"
	TypeMsgReceiveExperience    = "receive_experience"
	TypeMsgCommitReport         = "update_power_on_period"
	TypeMsgRewardComputingPower = "update_user_experience"
	TypeMsgAddCaller            = "add_caller"
	TypeMsgRemoveCaller         = "remove_caller"
	TypeMsgUpdateSaleLevel      = "update_sale_level"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgCreateCaptainNode{}
	_ sdk.Msg = &MsgWithdrawComputingPower{}
	_ sdk.Msg = &MsgCommitReport{}
	_ sdk.Msg = &MsgRewardComputingPower{}
	_ sdk.Msg = &MsgAddCaller{}
	_ sdk.Msg = &MsgRemoveCaller{}
	_ sdk.Msg = &MsgUpdateSaleLevel{}
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
		return errorsmod.Wrap(err, "invalid authority address")
	}
	return m.Params.Validate()
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgCreateCaptainNode(divisionId string, receiver string, sender string) *MsgCreateCaptainNode {
	return &MsgCreateCaptainNode{
		DivisionId: divisionId,
		Receiver:   receiver,
		Sender:     sender,
	}
}

// Route Implements Msg.
func (m MsgCreateCaptainNode) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgCreateCaptainNode) Type() string { return TypeMsgCreateCaptainNode }

// ValidateBasic Implements Msg.
func (msg MsgCreateCaptainNode) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Receiver); err != nil {
		return errorsmod.Wrap(err, "invalid receiver address")
	}
	if len(msg.DivisionId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "division id cannot be empty")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateCaptainNode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateCaptainNode) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgWithdrawComputingPower(nodeId string, computingPowerAmount uint64, sender string) *MsgWithdrawComputingPower {
	return &MsgWithdrawComputingPower{
		NodeId:               nodeId,
		ComputingPowerAmount: computingPowerAmount,
		Sender:               sender,
	}
}

// Route Implements Msg.
func (m MsgWithdrawComputingPower) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgWithdrawComputingPower) Type() string { return TypeMsgReceiveExperience }

// ValidateBasic Implements Msg.
func (msg MsgWithdrawComputingPower) ValidateBasic() error {
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

// GetSignBytes Implements Msg.
func (msg MsgWithdrawComputingPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawComputingPower) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgCommitReport(
	captainNodePowerOnPeriods []*CaptainNodePowerOnPeriod,
	sender string,
) *MsgCommitReport {
	return &MsgCommitReport{
		CaptainNodePowerOnPeriods: captainNodePowerOnPeriods,
		Sender:                    sender,
	}
}

// Route Implements Msg.
func (m MsgCommitReport) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgCommitReport) Type() string { return TypeMsgCommitReport }

// ValidateBasic Implements Msg.
func (msg MsgCommitReport) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")

	}
	if len(msg.CaptainNodePowerOnPeriods) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "captain node power on periods cannot be empty")
	}

	for _, captainNodePowerOnPeriod := range msg.CaptainNodePowerOnPeriods {
		if captainNodePowerOnPeriod.PowerOnPeriodRate.IsZero() || captainNodePowerOnPeriod.PowerOnPeriodRate.GT(sdk.NewDec(1)) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "power on period must be between 0 and 1")
		}
		captainNodePowerOnPeriod.NodeId = strings.TrimSpace(captainNodePowerOnPeriod.NodeId)
		if captainNodePowerOnPeriod.NodeId == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node id cannot be empty")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCommitReport) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgCommitReport) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgRewardComputingPower(
	extractableComputingPowers []*ExtractableComputingPower,
	sender string,
) *MsgRewardComputingPower {
	return &MsgRewardComputingPower{
		Sender:                     sender,
		ExtractableComputingPowers: extractableComputingPowers,
	}
}

// Route Implements Msg.
func (m MsgRewardComputingPower) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgRewardComputingPower) Type() string { return TypeMsgRewardComputingPower }

// ValidateBasic Implements Msg.
func (msg MsgRewardComputingPower) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")

	}
	if len(msg.ExtractableComputingPowers) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user experiences cannot be empty")
	}

	for _, userExperience := range msg.ExtractableComputingPowers {
		if userExperience.Amount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "experience cannot be zero")
		}
		if _, err := sdk.AccAddressFromBech32(userExperience.Owner); err != nil {
			return errorsmod.Wrap(err, "invalid receiver address")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRewardComputingPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRewardComputingPower) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgAddCaller() *MsgAddCaller {
	return &MsgAddCaller{}
}

// Route Implements Msg.
func (m MsgAddCaller) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgAddCaller) Type() string { return TypeMsgAddCaller }

// ValidateBasic Implements Msg.
func (msg MsgAddCaller) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid sender address")
	}

	if len(msg.Callers) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "caller cannot be empty")
	}

	for _, caller := range msg.Callers {
		if _, err := sdk.AccAddressFromBech32(caller); err != nil {
			return errorsmod.Wrap(err, "invalid caller address")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddCaller) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgAddCaller) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgRemoveCaller() *MsgRemoveCaller {
	return &MsgRemoveCaller{}
}

// Route Implements Msg.
func (m MsgRemoveCaller) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgRemoveCaller) Type() string { return TypeMsgRemoveCaller }

// ValidateBasic Implements Msg.
func (msg MsgRemoveCaller) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid sender address")
	}

	if len(msg.Callers) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "caller cannot be empty")
	}

	for _, caller := range msg.Callers {
		if _, err := sdk.AccAddressFromBech32(caller); err != nil {
			return errorsmod.Wrap(err, "invalid caller address")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveCaller) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveCaller) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

func NewMsgUpdateSaleLevel() *MsgUpdateSaleLevel {
	return &MsgUpdateSaleLevel{}
}

// Route Implements Msg.
func (m MsgUpdateSaleLevel) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgUpdateSaleLevel) Type() string { return TypeMsgUpdateSaleLevel }

// ValidateBasic Implements Msg.
func (msg MsgUpdateSaleLevel) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid sender address")
	}

	if msg.SaleLevel <= 1 || msg.SaleLevel > 5 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "sale level must be between 1 and 5")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateSaleLevel) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdateSaleLevel) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
