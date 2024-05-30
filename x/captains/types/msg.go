package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgCreateCaptainNode{}
	_ sdk.Msg = &MsgCommitReport{}
	_ sdk.Msg = &MsgAddAuthorizedMembers{}
	_ sdk.Msg = &MsgRemoveAuthorizedMembers{}
	_ sdk.Msg = &MsgUpdateSaleLevel{}
	_ sdk.Msg = &MsgCommitComputingPower{}
	_ sdk.Msg = &MsgClaimComputingPower{}
)

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.Params.CaptainsTotalCount == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "captains total count cannot be zero")
	}

	if msg.Params.MinimumPowerOnPeriod > msg.Params.MaximumPowerOnPeriod {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum power on period cannot be greater than maximum power on period")
	}

	if !msg.Params.HalvingEraCoefficient.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "halving era coefficient must be greater than zero")
	}

	if !msg.Params.CaptainsConstant.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "captains constant must be greater than zero")
	}

	if !msg.Params.TechProgressCoefficientCardinality.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "tech progress coefficient cardinality must be greater than zero")
	}

	if msg.Params.CurrentSaleLevel == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "current sale level cannot be zero")
	}

	if len(msg.Params.AuthorizedMembers) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "authorized members cannot be empty")
	}

	for _, member := range msg.Params.AuthorizedMembers {
		if _, err := sdk.AccAddressFromBech32(member); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

// NewMsgCreateCaptainNode creates a new MsgCreateCaptainNode instance
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

// NewMsgCommitReport creates a new MsgCommitReport instance
// report: interface of *ReportDigest, *ReportBatch, *ReportEnd
func NewMsgCommitReport(authority string, reportType ReportType, report any) (*MsgCommitReport, error) {
	res := MsgCommitReport{
		Authority:  authority,
		ReportType: reportType,
	}

	switch v := report.(type) {
	case *ReportDigest:
		anyV, err := types.NewAnyWithValue(v)
		if err != nil {
			return nil, err
		}
		res.Report = anyV
	case *ReportBatch:
		anyV, err := types.NewAnyWithValue(v)
		if err != nil {
			return nil, err
		}
		res.Report = anyV
	case *ReportEnd:
		anyV, err := types.NewAnyWithValue(v)
		if err != nil {
			return nil, err
		}
		res.Report = anyV
	default:
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid report type")
	}

	return &res, nil
}

// ValidateBasic Implements Msg.
func (msg *MsgCommitReport) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.ReportType > ReportType_REPORT_TYPE_END || msg.ReportType == ReportType_REPORT_TYPE_UNSPECIFIED {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid report type")
	}

	if msg.Report == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "report cannot be nil")
	}

	bz := msg.Report.GetValue()

	switch msg.ReportType {
	case ReportType_REPORT_TYPE_DIGEST:
		var digest ReportDigest
		err := digest.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		err = digest.ValidateBasic()
		if err != nil {
			return err
		}
	case ReportType_REPORT_TYPE_BATCH:
		var batch ReportBatch
		err := batch.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		err = batch.ValidateBasic()
		if err != nil {
			return err
		}
	case ReportType_REPORT_TYPE_END:
		var end ReportEnd
		err := end.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		if end.EpochId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch id is zero")
		}
	}
	return nil
}

// GetSigners Implements Msg.
func (msg *MsgCommitReport) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{fromAddress}
}

// NewAddAuthorizedMembers creates a new MsgAddAuthorizedMembers instance
func NewAddAuthorizedMembers(authority string, members []string) *MsgAddAuthorizedMembers {
	return &MsgAddAuthorizedMembers{
		Authority: authority,
		Members:   members,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgAddAuthorizedMembers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if len(msg.Members) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "member cannot be empty")
	}

	seenMap := make(map[string]bool)
	for _, member := range msg.Members {
		if seenMap[member] {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate member")
		}
		if _, err := sdk.AccAddressFromBech32(member); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
		seenMap[member] = true
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgAddAuthorizedMembers) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// NewMsgRemoveAuthorizedMembers creates a new MsgRemoveAuthorizedMembers instance
func NewMsgRemoveAuthorizedMembers(authority string, members []string) *MsgRemoveAuthorizedMembers {
	return &MsgRemoveAuthorizedMembers{
		Authority: authority,
		Members:   members,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgRemoveAuthorizedMembers) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if len(msg.Members) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "member cannot be empty")
	}

	seenMap := make(map[string]bool)
	for _, member := range msg.Members {
		if seenMap[member] {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate member")
		}
		if _, err := sdk.AccAddressFromBech32(member); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
		seenMap[member] = true
	}
	return nil
}

// GetSigners Implements Msg.
func (msg *MsgRemoveAuthorizedMembers) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateSaleLevel creates a new MsgUpdateSaleLevel instance
func NewMsgUpdateSaleLevel(authority string, level uint64) *MsgUpdateSaleLevel {
	return &MsgUpdateSaleLevel{
		Authority: authority,
		SaleLevel: level,
	}
}

// ValidateBasic Implements Msg.
func (msg *MsgUpdateSaleLevel) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.SaleLevel < 1 || msg.SaleLevel > 5 {
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

	seenMap := make(map[string]bool)
	for _, reward := range msg.ComputingPowerRewards {
		if seenMap[reward.Owner] {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate receiver")
		}
		if reward.Amount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "computing power cannot be zero")
		}
		if _, err := sdk.AccAddressFromBech32(reward.Owner); err != nil {
			return errorsmod.Wrap(err, "invalid receiver address")
		}
		seenMap[reward.Owner] = true
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
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "computing power cannot be zero")
	}

	return nil
}

// GetSigners Implements Msg.
func (msg *MsgClaimComputingPower) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}
