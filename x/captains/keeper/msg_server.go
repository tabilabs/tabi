package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the mint MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// UpdateParams implement the interface of types.MsgServer
func (m msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority.String(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

// CreateCaptainNode implement the interface of types.MsgServer
func (m msgServer) CreateCaptainNode(
	goCtx context.Context,
	msg *types.MsgCreateCaptainNode,
) (*types.MsgCreateCaptainNodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	nodeID, err := m.k.CreateNode(ctx, msg.DivisionId, owner)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateNode,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeID),
			sdk.NewAttribute(types.AttributeKeyDivisionID, msg.DivisionId),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
		),
	})

	return &types.MsgCreateCaptainNodeResponse{
		NodeId: nodeID,
	}, nil
}

// CommitReport implement the interface of types.MsgServer
func (m msgServer) CommitReport(
	goCtx context.Context,
	msg *types.MsgCommitReport,
) (*types.MsgCommitReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	report, err := m.k.ValidateReport(ctx, msg.ReportType, msg.Report)
	if err != nil {
		return &types.MsgCommitReportResponse{}, err
	}

	if err := m.k.CommitReport(ctx, report); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
		),
		sdk.NewEvent(
			types.EventTypeCommitReport,
			sdk.NewAttribute(types.AttributeKeyReportType, msg.ReportType.String()),
		),
	})

	return &types.MsgCommitReportResponse{}, nil
}

// AddAuthorizedMembers implement the interface of types.MsgServer
func (m msgServer) AddAuthorizedMembers(
	goCtx context.Context,
	msg *types.MsgAddAuthorizedMembers,
) (*types.MsgAddAuthorizedMembersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	if err := m.k.SetAuthorizedMembers(ctx, msg.Members); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
	))

	return &types.MsgAddAuthorizedMembersResponse{}, nil
}

// RemoveAuthorizedMembers implement the interface of types.MsgServer
func (m msgServer) RemoveAuthorizedMembers(
	goCtx context.Context,
	msg *types.MsgRemoveAuthorizedMembers,
) (*types.MsgRemoveAuthorizedMembersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	if err := m.k.DeleteAuthorizedMembers(ctx, msg.Members); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
	))

	return &types.MsgRemoveAuthorizedMembersResponse{}, nil
}

// UpdateSaleLevel implement the interface of types.MsgServer
func (m msgServer) UpdateSaleLevel(
	goCtx context.Context,
	msg *types.MsgUpdateSaleLevel,
) (*types.MsgUpdateSaleLevelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	before, err := m.k.UpdateSaleLevel(ctx, msg.SaleLevel)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateSaleLevel,
			sdk.NewAttribute(types.AttributeKeySaleLevelBefore, fmt.Sprintf("%d", before)),
			sdk.NewAttribute(types.AttributeKeySaleLevelAfter, fmt.Sprintf("%d", msg.SaleLevel)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
		),
	})

	return &types.MsgUpdateSaleLevelResponse{}, nil
}

// CommitComputingPower implement the interface of types.MsgServer
func (m msgServer) CommitComputingPower(
	goCtx context.Context,
	msg *types.MsgCommitComputingPower,
) (*types.MsgCommitComputingPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	if len(msg.ComputingPowerRewards) == 0 {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid user experience; empty",
		)
	}

	events := make([]sdk.Event, 0)
	for _, cpr := range msg.ComputingPowerRewards {
		owner, err := sdk.AccAddressFromBech32(cpr.Owner)
		if err != nil {
			return nil, err
		}
		before, after, err := m.k.CommitComputingPower(ctx, cpr.Amount, owner)
		if err != nil {
			return nil, err
		}

		events = append(events, sdk.NewEvent(
			types.EventTypeCommitComputingPower,
			sdk.NewAttribute(types.AttributeKeyOwner, cpr.Owner),
			sdk.NewAttribute(types.AttributeKeyComputingPowerBefore, fmt.Sprintf("%d", before)),
			sdk.NewAttribute(types.AttributeKeyComputingPowerAfter, fmt.Sprintf("%d", after)),
		))
	}

	events = append(events, sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
	))
	ctx.EventManager().EmitEvents(events)

	return &types.MsgCommitComputingPowerResponse{}, nil
}

// ClaimComputingPower implement the interface of types.MsgServer
func (m msgServer) ClaimComputingPower(
	goCtx context.Context,
	msg *types.MsgClaimComputingPower,
) (*types.MsgClaimComputingPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if !m.k.HasAuthorizedMember(ctx, authority) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in authorized members list",
		)
	}

	owner, found := m.k.GetNodeOwner(ctx, msg.NodeId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNodeNotExists, msg.NodeId)
	}

	if err := m.k.UpdateNode(ctx, msg.NodeId, msg.ComputingPowerAmount, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimComputingPower,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyComputingPower, fmt.Sprintf("%d", msg.ComputingPowerAmount)),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgClaimComputingPowerResponse{}, nil
}
