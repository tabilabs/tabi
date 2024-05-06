package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec/types"
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
	_ sdk.Msg = &MsgUpdateParams{}
)

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
func NewMsgCommitReport(authority string, reportType ReportType, report any) *MsgCommitReport {
	res := MsgCommitReport{
		Authority:  authority,
		ReportType: reportType,
		Report:     nil,
	}

	switch v := report.(type) {
	case *ReportDigest:
		anyv, _ := types.NewAnyWithValue(v)
		res.Report = anyv
	case *ReportBatch:
		anyv, _ := types.NewAnyWithValue(v)
		res.Report = anyv
	case *ReportEnd:
		anyv, _ := types.NewAnyWithValue(v)
		res.Report = anyv
	}

	return &res
}

// ValidateBasic Implements Msg.
func (msg *MsgCommitReport) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if msg.Report == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "report cannot be nil")
	}

	if msg.ReportType > ReportType_REPORT_TYPE_END || msg.ReportType == ReportType_REPORT_TYPE_DIGEST {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid report type")
	}

	bz := msg.Report.GetValue()

	switch msg.ReportType {
	case ReportType_REPORT_TYPE_DIGEST:
		var digest ReportDigest
		err := digest.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		if digest.EpochId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch id is zero")
		}
		if digest.TotalBatchCount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total batch is zero")
		}
		if digest.TotalNodeCount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total node is zero")
		}
		if digest.MaximumNodeCountPerBatch == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max node per batch id is zero")
		}
		if digest.GlobalOnOperationRatio.IsZero() {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "operation ratio is zero")
		}

	case ReportType_REPORT_TYPE_BATCH:
		var batch ReportBatch
		err := batch.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		if batch.EpochId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "epoch id is zero")
		}
		if batch.BatchId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "batch id is zero")
		}
		if batch.NodeCount == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node count is zero")
		}
		if uint64(len(batch.NodeIds)) != batch.NodeCount {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "node count node ids length unmatched")
		}

	case ReportType_REPORT_TYPE_END:
		var end ReportEnd
		err := end.Unmarshal(bz)
		if err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		if end.Epoch == 0 {
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

	for _, member := range msg.Members {
		if _, err := sdk.AccAddressFromBech32(member); err != nil {
			return errorsmod.Wrap(err, "invalid member address")
		}
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
