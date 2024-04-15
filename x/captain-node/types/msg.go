package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// caption-node message types
const (
	TypeMsgMint                 = "mint"
	TypeMsgReceiveExperience    = "receive_experience"
	TypeMsgUpdatePowerOnPeriod  = "update_power_on_period"
	TypeMsgUpdateUserExperience = "update_user_experience"
	TypeMsgAddCaller            = "add_caller"
	TypeMsgRemoveCaller         = "remove_caller"
	TypeMsgUpdateSaleLevel      = "update_sale_level"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgMint{}
	_ sdk.Msg = &MsgWithdrawExperience{}
	_ sdk.Msg = &MsgUpdatePowerOnPeriod{}
	_ sdk.Msg = &MsgUpdateUserExperience{}
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

func NewMsgMint(divisionId string, receiver string, sender string) *MsgMint {
	return &MsgMint{
		DivisionId: divisionId,
		Receiver:   receiver,
		Sender:     sender,
	}
}

// Route Implements Msg.
func (m MsgMint) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgMint) Type() string { return TypeMsgMint }

// ValidateBasic Implements Msg.
func (msg MsgMint) ValidateBasic() error {
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
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgWithdrawExperience(nodeId string, experience uint64, sender string) *MsgWithdrawExperience {
	return &MsgWithdrawExperience{
		NodeId:     nodeId,
		Experience: experience,
		Sender:     sender,
	}
}

// Route Implements Msg.
func (m MsgWithdrawExperience) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgWithdrawExperience) Type() string { return TypeMsgReceiveExperience }

// ValidateBasic Implements Msg.
func (msg MsgWithdrawExperience) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if len(msg.NodeId) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node id cannot be empty")
	}
	if msg.Experience == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "experience cannot be zero")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawExperience) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawExperience) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgUpdatePowerOnPeriod(
	captainNodePowerOnPeriods []*CaptainNodePowerOnPeriod,
	sender string,
) *MsgUpdatePowerOnPeriod {
	return &MsgUpdatePowerOnPeriod{
		CaptainNodePowerOnPeriods: captainNodePowerOnPeriods,
		Sender:                    sender,
	}
}

// Route Implements Msg.
func (m MsgUpdatePowerOnPeriod) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgUpdatePowerOnPeriod) Type() string { return TypeMsgUpdatePowerOnPeriod }

// ValidateBasic Implements Msg.
func (msg MsgUpdatePowerOnPeriod) ValidateBasic() error {
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
func (msg MsgUpdatePowerOnPeriod) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdatePowerOnPeriod) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func NewMsgUpdateUserExperience(
	userExperiences []*UserExperience,
	sender string,
) *MsgUpdateUserExperience {
	return &MsgUpdateUserExperience{
		Sender:          sender,
		UserExperiences: userExperiences,
	}
}

// Route Implements Msg.
func (m MsgUpdateUserExperience) Route() string { return RouterKey }

// Type Implements Msg.
func (m MsgUpdateUserExperience) Type() string { return TypeMsgUpdateUserExperience }

// ValidateBasic Implements Msg.
func (msg MsgUpdateUserExperience) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")

	}
	if len(msg.UserExperiences) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user experiences cannot be empty")
	}

	for _, userExperience := range msg.UserExperiences {
		if userExperience.Experience == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "experience cannot be zero")
		}
		if _, err := sdk.AccAddressFromBech32(userExperience.Receiver); err != nil {
			return errorsmod.Wrap(err, "invalid receiver address")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateUserExperience) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgUpdateUserExperience) GetSigners() []sdk.AccAddress {
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
